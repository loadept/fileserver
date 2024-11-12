package util

import "github.com/golang-jwt/jwt"

type TokenClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
