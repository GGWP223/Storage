package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	RefreshToken(token string) (string, error)
	ValidateToken(token string) (*Claims, error)
	GenerateToken(login, uid, tokenType string) (string, error)
}

type service struct {
	jwtSecret string
}

func NewService(jwtSecret string) Service {
	return &service{jwtSecret: jwtSecret}
}

func (s *service) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)

	if err != nil {
		return "", err
	}

	if claims.Type != RefreshToken {
		return "", errors.New("wrong token type")
	}

	key, err := s.GenerateToken(claims.Login, claims.Uid, claims.Type)

	if err != nil {
		return "", err
	}

	return key, nil
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

func (s *service) GenerateToken(login, uid, tokenType string) (string, error) {
	claims := &Claims{
		Uid:   uid,
		Login: login,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "File Storage",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}
