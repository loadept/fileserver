package middleware

import (
	"log"
	"net/http"

	"github.com/jsusmachaca/fileserver/internal/util"
	"github.com/jsusmachaca/go-router/pkg/response"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) < 1 {
			response.JsonErrorFromString(w, "Token not provided", http.StatusUnauthorized)
			return
		}
		tokenString := authorization[len("Bearer "):]

		_, err := util.ValidateToken(tokenString)
		if err != nil {
			log.Println(err.Error())
			response.JsonErrorFromString(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
