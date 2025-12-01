package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sssseraphim/fitterBy/internal/auth"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserTypeKey contextKey = "user_type"
	EmailKey    contextKey = "email"
)

func AuthMiddleware(jwtConfig *auth.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "Authorization header invalid"}`, http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims, valid, err := jwtConfig.ValidateAccessToken(tokenString)
			if err != nil || !valid {
				http.Error(w, `{"error": "Token invalid or expired"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserTypeKey, claims.UserType)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
