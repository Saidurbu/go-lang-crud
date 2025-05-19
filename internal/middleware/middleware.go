package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Saidurbu/go-lang-crud/internal/handlers/student"
	"github.com/golang-jwt/jwt/v5"
)

var secret_key = []byte("secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret_key, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), student.EmailContextKey(), claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
