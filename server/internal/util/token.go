package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id string, username string) (string, error) {
	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := TokenClaims{
		UserId:   id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "jsusmachaca",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenstring string) (jwt.MapClaims, error) {
	SECRET_KEY := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Error mapping token")
	}

	return claims, nil
}
