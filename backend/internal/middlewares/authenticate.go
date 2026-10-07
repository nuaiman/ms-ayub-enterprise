package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"backend/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorJson(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)

		userID, err := utils.VerifyAccessToken(token, secret)
		if err != nil {
			log.Printf("[AUTH] token verification failed: %v", err)
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		if userID == nil {
			log.Printf("[AUTH] token verified but userID is nil")
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, *userID)
		next(w, r.WithContext(ctx))
	}
}
