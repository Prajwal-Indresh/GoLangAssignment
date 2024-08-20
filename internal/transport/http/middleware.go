package http

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTAuthMiddleware(next http.HandlerFunc, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[len("Bearer "):]
		}

		token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*AuthClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		log.Printf("Extracted user ID: %s", claims.UserID)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
