package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Register(request Request) error
	Login(request Request) (string, error)
	RefreshToken(token string) (string, error)
	ParseToken(token string) (string, string, error)
	ValidateToken(token string) (*Claims, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo, jwtSecret}
}

func (s *service) Register(request Request) error {
	_, err := s.repo.GetUser(request)

	if err == nil {
		return errors.New("user already exists")
	}

	return s.repo.CreateUser(request)
}

func (s *service) Login(request Request) (string, error) {
	user, err := s.repo.GetUser(request)

	if err != nil {
		return "", err
	}

	token, err := s.generateJWT(request.Login, user.Uid)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RefreshToken(token string) (string, error) {
	claims, err := s.ValidateToken(token)

	if err != nil {
		return "", err
	}

	key, err := s.generateJWT(claims.Login, claims.Uid)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *service) ParseToken(token string) (string, string, error) {
	claims, err := s.ValidateToken(token)

	if err != nil {
		return "", "", err
	}

	return claims.Login, claims.Uid, err
}

func (s *service) ValidateToken(token string) (*Claims, error) {
	var key, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := key.Claims.(*Claims)

	if !ok || !key.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *service) generateJWT(login, uid string) (string, error) {
	claims := &Claims{
		Uid:   uid,
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "File Storage",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}
