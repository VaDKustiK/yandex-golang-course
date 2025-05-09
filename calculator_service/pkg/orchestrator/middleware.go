package orchestrator

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type ctxKey string

const userCtxKey = ctxKey("user_id")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := h[len("Bearer "):]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) uint {
	if uid, ok := ctx.Value(userCtxKey).(uint); ok {
		return uid
	}
	return 0
}
