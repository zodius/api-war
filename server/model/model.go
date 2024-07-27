package model

import (
	"errors"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrUserExist          = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	FieldCount     = 1000000
	BatchSize      = 1000
	TypeWebservice = "webservice"
	TypeRestful    = "restful"
	TypeGraphql    = "graphql"
	TypeGrpc       = "grpc"
)

/*
	Redis schema:
	- Hashmap:
		{"user:<username>" : {"password":<password>}, "id":<id>} }
		{"fields:<type>:conquerer": {<fieldID>:<owner>}}
	- Key:
		{"token:<token>" : <username>}
		{"usercount": int}
	- ZSet:
		{"users": [<username> <id>]}
		{"score:conquerCount": [<username> <count>]}
		{"score:conquerHistory:webservice": [<username> <count>]}
		{"score:conquerHistory:restful": [<username> <count>]}
		{"score:conquerHistory:graphql": [<username> <count>]}
		{"score:conquerHistory:grpc": [<username> <count>]}
	- Bitmap:
		{"user:<username>:conquerField:<type>": <fieldID>}
*/

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Owner struct {
	ConquerType string `json:"conquerType"` // webservice, restful, grqphql, grpc
	Owner       string `json:"owner"`
}

type Field struct {
	FieldID   int     `json:"field_id"`
	Conquerer []Owner `json:"conquerer"`
}

type Map struct {
	Fields []Field `json:"fields"`
}

func (m *Map) Representation() interface{} {
	// convert map object to json encoded string
	represent := make(map[int](map[string]string))
	for _, field := range m.Fields {
		for _, conquerer := range field.Conquerer {
			if _, ok := represent[field.FieldID]; !ok {
				represent[field.FieldID] = make(map[string]string)
			}
			represent[field.FieldID][conquerer.ConquerType] = conquerer.Owner
		}
	}
	return represent
}

type Score struct {
	Username            string         `json:"username"`
	ConquerFieldCount   int            `json:"conquerFieldCount"`
	ConquerHistoryCount map[string]int `json:"conquerHistoryCount"`
}

type Service interface {
	// auth
	Login(username, password string) (token string, err error)
	Register(username, password string) error
	// basic information
	GetCurrentMap() (Map Map, err error)
	GetUserList(token string) (userList []User, err error) // this is used to get username by id for each client
	// services for exploit
	GetUserConquerField(token string, conquerType string) ([]int, error)
	ConquerField(token string, fieldID int, conquerType string) error
	// scoreboard
	GetScoreboard() (scoreList []Score, err error)
}

type Repo interface {
	GetUser(username string) (User, error)
	CreateUser(username, password string) error
	CreateToken(username string) (token string, err error)
	GetTokenUsername(token string) (username string, err error)
	GetMap() (Map, error)
	GetUserList() (userList []User, err error)
	GetUserConquerField(username string, conquerType string) ([]int, error)
	GetScoreboard() (scoreList []Score, err error)
	SetFieldConquerer(fieldID int, conquerType, username string) error
	AddScore(username string, fieldID int, conquerType string) error
}
