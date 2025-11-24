package jwt

import "github.com/golang-jwt/jwt/v5"

const (
	RefreshToken = "refresh_token"
	AccessToken  = "access_token"
)

type Claims struct {
	Uid   string `json:"uid"`
	Login string `json:"login"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}
