package repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zodius/api-war/model"
)

type repo struct {
	client *redis.Client
	lock   *sync.Mutex
}

func NewRepo(client *redis.Client) model.Repo {
	return &repo{
		client: client,
		lock:   new(sync.Mutex),
	}
}

func (r *repo) GetUser(username string) (model.User, error) {
	values, err := r.client.HMGet(context.Background(), fmt.Sprintf("user:%s", username),
		"password", "id",
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, model.ErrNotFound
		}
		return model.User{}, err
	}

	password, ok := values[0].(string)
	if !ok {
		return model.User{}, model.ErrNotFound
	}

	idStr, ok := values[1].(string)
	if !ok {
		return model.User{}, model.ErrNotFound
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Username: username,
		Password: password,
		ID:       id,
	}, nil
}

func (r *repo) CreateUser(username, password string) error {
	// create user requires lock to prevent race condition in user count
	r.lock.Lock()
	defer r.lock.Unlock()
	// get user count as id
	userCount, err := r.client.Get(context.Background(), "usercount").Int()
	if err != nil {
		// if not exist, set 0
		if errors.Is(err, redis.Nil) {
			userCount = 0
			err = r.client.Set(context.Background(), "usercount", 0, 0).Err()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	userID := userCount + 1

	// create user
	if err := r.client.HSet(context.Background(), fmt.Sprintf("user:%s", username),
		"password", password,
		"id", userID,
	).Err(); err != nil {
		return err
	}

	// increment user count
	if err := r.client.Incr(context.Background(), "usercount").Err(); err != nil {
		return err
	}

	// create score
	if err := r.client.ZAdd(context.Background(), "score:conquerCount", redis.Z{
		Score:  0,
		Member: username,
	}).Err(); err != nil {
		return err
	}

	conquerTypes := []string{"webservice", "restful", "graphql", "grpc"}
	for _, conquerType := range conquerTypes {
		// create conquer history score
		if err := r.client.ZAdd(context.Background(), fmt.Sprintf("score:conquerHistory:%s", conquerType), redis.Z{
			Score:  0,
			Member: username,
		}).Err(); err != nil {
			return err
		}

		bitmapKey := fmt.Sprintf("user:%s:conquerField:%s", username, conquerType)
		if err := r.client.SetBit(context.Background(), bitmapKey, 0, 0).Err(); err != nil {
			return err
		}
	}

	// add user to users zset
	if err := r.client.ZAdd(context.Background(), "users", redis.Z{
		Score:  float64(userID),
		Member: username,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repo) CreateToken(username string) (token string, err error) {
	token, err = randomToken()
	if err != nil {
		return "", err
	}

	err = r.client.Set(context.Background(), fmt.Sprintf("token:%s", token), username, 15*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *repo) GetTokenUsername(token string) (username string, err error) {
	username, err = r.client.Get(context.Background(), fmt.Sprintf("token:%s", token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.ErrNotFound
		}
		return "", err
	}
	return username, nil
}

func (r *repo) GetMap(startInput, endInput int) (model.Map, error) {
	mapMap := make(map[int]model.Field, endInput-startInput+1)

	conquerTypes := []string{"webservice", "restful", "graphql", "grpc"}
	for _, conquerType := range conquerTypes {
		start, end := startInput, endInput
		// get conquerer within range in batch
		// for each batch, get conquerer using hmget
		for start <= end {
			// check remaining is less than batchsize
			batchSize := model.BatchSize
			if start+batchSize > end {
				batchSize = end - start + 1
			}
			fields := make([]string, 0, batchSize)
			for i := start; i < start+batchSize-1; i++ {
				fields = append(fields, strconv.Itoa(i))
			}
			conquerers, err := r.client.HMGet(context.Background(), fmt.Sprintf("fields:%s:conquerer", conquerType),
				fields...,
			).Result()
			if err != nil {
				return model.Map{}, err
			}
			for i, conquerer := range conquerers {
				fieldID := start + i
				field, ok := mapMap[fieldID]
				if !ok {
					field = model.Field{
						FieldID:   fieldID,
						Conquerer: make([]model.Owner, 0),
					}
				}
				conquererName := ""
				if conquerer != nil {
					conquererName = conquerer.(string)
				}

				field.Conquerer = append(field.Conquerer, model.Owner{
					ConquerType: conquerType,
					Owner:       conquererName,
				})
				mapMap[fieldID] = field
			}
			start += batchSize
		}
	}

	fields := make([]model.Field, 0, len(mapMap))
	for _, field := range mapMap {
		fields = append(fields, field)
	}
	return model.Map{
		Fields: fields,
	}, nil
}

func (r *repo) GetUserList() ([]model.User, error) {
	users := make([]model.User, 0)
	// get all users
	zrange, err := r.client.ZRangeWithScores(context.Background(), "users", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for _, z := range zrange {
		users = append(users, model.User{
			Username: z.Member.(string),
			ID:       int(z.Score),
		})
	}
	return users, nil
}

func (r *repo) GetUserConquerField(username string, conquerType string) ([]int, error) {
	fieldCount := model.FieldCount
	batchSize := model.BatchSize

	result := make([]int, 0)

	// get bitmap in batch
	for i := 0; i < fieldCount; i += batchSize {
		ids := make([]int, 0)
		// calculate start and end
		start := i
		end := i + batchSize - 1
		// count bit within range
		count, err := r.client.BitCount(context.Background(), fmt.Sprintf("user:%s:conquerField:%s", username, conquerType), &redis.BitCount{
			Start: int64(start),
			End:   int64(end),
		}).Result()
		if err != nil {
			return nil, err
		}
		if count == 0 {
			continue
		}

		// get bitpos in batch
		for start <= end {
			pos, err := r.client.BitPos(context.Background(), fmt.Sprintf("user:%s:conquerField:%s", username, conquerType),
				1,
				int64(start), int64(end),
			).Result()
			if err != nil {
				return nil, err
			}
			if pos == -1 {
				break
			}
			ids = append(ids, int(pos))
			start = int(pos) + 1
		}

		result = append(result, ids...)
	}
	return result, nil
}

func (r *repo) GetScoreboard() ([]model.Score, error) {
	// make hashmap for calculate
	scoreMap := make(map[string]model.Score)

	zrangeKey := []string{
		"score:conquerCount",
		"score:conquerHistory:webservice",
		"score:conquerHistory:restful",
		"score:conquerHistory:graphql",
		"score:conquerHistory:grpc",
	}

	for _, key := range zrangeKey {
		zrange, err := r.client.ZRangeWithScores(context.Background(), key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		for _, z := range zrange {
			score, ok := scoreMap[z.Member.(string)]
			if !ok {
				// if not exist, create new score
				score = model.Score{
					Username:            z.Member.(string),
					ConquerHistoryCount: make(map[string]int),
				}
			}

			switch key {
			case "score:conquerCount":
				score.ConquerFieldCount = int(z.Score)
			case "score:conquerHistory:webservice":
				score.ConquerHistoryCount["webservice"] = int(z.Score)
			case "score:conquerHistory:restful":
				score.ConquerHistoryCount["restful"] = int(z.Score)
			case "score:conquerHistory:graphql":
				score.ConquerHistoryCount["graphql"] = int(z.Score)
			case "score:conquerHistory:grpc":
				score.ConquerHistoryCount["grpc"] = int(z.Score)
			}

			scoreMap[z.Member.(string)] = score
		}
	}

	// convert map to slice
	var scoreList []model.Score
	for _, score := range scoreMap {
		scoreList = append(scoreList, score)
	}
	return scoreList, nil
}

func (r *repo) SetFieldConquerer(fieldID int, conquerType, username string) error {
	// set bit in user bitmap
	if err := r.client.SetBit(context.Background(), fmt.Sprintf("user:%s:conquerField:%s", username, conquerType), int64(fieldID), 1).Err(); err != nil {
		return err
	}
	// set conquerer in field map
	if err := r.client.HSet(context.Background(),
		fmt.Sprintf("fields:%s:conquerer", conquerType),
		fieldID, username,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *repo) AddScore(username string, fieldID int, conquerType string) error {
	// add score:conquerCount
	if err := r.client.ZIncrBy(context.Background(), "score:conquerCount", 1, username).Err(); err != nil {
		return err
	}
	// add score:conquerHistory:<conquerType>
	switch conquerType {
	case "webservice":
		if err := r.client.ZIncrBy(context.Background(), "score:conquerHistory:webservice", 1, username).Err(); err != nil {
			return err
		}
	case "restful":
		if err := r.client.ZIncrBy(context.Background(), "score:conquerHistory:restful", 1, username).Err(); err != nil {
			return err
		}
	case "graphql":
		if err := r.client.ZIncrBy(context.Background(), "score:conquerHistory:graphql", 1, username).Err(); err != nil {
			return err
		}
	case "grpc":
		if err := r.client.ZIncrBy(context.Background(), "score:conquerHistory:grpc", 1, username).Err(); err != nil {
			return err
		}
	}

	return nil
}

func randomToken() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
