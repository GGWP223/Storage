package auth

import "github.com/golang-jwt/jwt/v5"

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Uid       string `json:"uid"`
	CreatedAt string `json:"create_data"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type Claims struct {
	Uid   string `json:"uid"`
	Login string `json:"login"`
	jwt.RegisteredClaims
}
