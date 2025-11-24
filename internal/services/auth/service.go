package auth

import (
	proto "File_Storage/internal/api/jwt"
	"File_Storage/internal/services/jwt"
	"context"
	"errors"
)

type Service interface {
	Register(request Request) error
	Login(request Request) (access, refresh string, err error)
}

type service struct {
	repo      Repository
	jwtClient proto.JwtServiceClient
}

func NewService(repo Repository, client proto.JwtServiceClient) Service {
	return &service{repo: repo, jwtClient: client}
}

func (s *service) Register(request Request) error {
	_, err := s.repo.GetUser(request)

	if err == nil {
		return errors.New("user already exists")
	}

	return s.repo.CreateUser(request)
}

func (s *service) Login(request Request) (access, refresh string, err error) {
	user, err := s.repo.GetUser(request)

	if err != nil {
		return "", "", err
	}

	accessResponse, err := s.jwtClient.GenerateToken(context.Background(), &proto.GenerateRequest{
		Uid:       user.Uid,
		Login:     user.Login,
		TokenType: jwt.AccessToken,
	})

	if err != nil {
		return "", "", err
	}

	refreshResponse, err := s.jwtClient.GenerateToken(context.Background(), &proto.GenerateRequest{
		Uid:       user.Uid,
		Login:     user.Login,
		TokenType: jwt.RefreshToken,
	})

	if err != nil {
		return "", "", err
	}

	return accessResponse.Token, refreshResponse.Token, nil
}
