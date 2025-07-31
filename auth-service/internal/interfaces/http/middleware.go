package http

import (
	"context"
	"net/http"
	"strings"

	jwtmgr "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/jwt"
)

type ctxKey string

const userIDKey ctxKey = "uid"

func AuthMiddleware(jwt *jwtmgr.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing bearer", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := jwt.Parse(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
