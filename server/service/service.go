package service

import (
	"errors"

	"github.com/zodius/api-war/model"
)

type service struct {
	repo model.Repo
}

func NewService(
	repo model.Repo,
) model.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(username, password string) error {
	// check if username already exists
	_, err := s.repo.GetUser(username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// create new user
			return s.repo.CreateUser(username, password)
		} else {
			return err
		}
	}

	return model.ErrUserExist
}

func (s *service) Login(username, password string) (token string, err error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	if user.Password != password {
		return "", model.ErrInvalidCredentials
	}

	return s.repo.CreateToken(username)
}

func (s *service) GetCurrentMap() (Map model.Map, err error) {
	return s.repo.GetMap()
}

func (s *service) GetUserList(token string) (userList []model.User, err error) {
	// verify token
	_, err = s.repo.GetTokenUsername(token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserList()
}

func (s *service) GetUserConquerField(token string, conquerType string) (fields []int, err error) {
	username, err := s.repo.GetTokenUsername(token)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserConquerField(username, conquerType)
}

func (s *service) ConquerField(token string, fieldID int, conquerType string) error {
	username, err := s.repo.GetTokenUsername(token)
	if err != nil {
		return err
	}

	err = s.repo.SetFieldConquerer(fieldID, conquerType, username)
	if err != nil {
		return err
	}

	// add score
	return s.repo.AddScore(username, fieldID, conquerType)
}

func (s *service) GetScoreboard() (scoreList []model.Score, err error) {
	return s.repo.GetScoreboard()
}
