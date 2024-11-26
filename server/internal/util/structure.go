package util

import "github.com/golang-jwt/jwt"

type TokenClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

type ResponseWithToken struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
