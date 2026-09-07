package middlewares

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("========================================")
		fmt.Println("[AUTH] AUTHENTICATE MIDDLEWARE CALLED")
		fmt.Printf("[AUTH] Method: %s, Path: %s\n", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		fmt.Printf("[AUTH] Authorization Header: '%s'\n", authHeader)

		if authHeader == "" {
			fmt.Println("[AUTH] ERROR: Missing authorization header")
			utils.ErrorJson(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			fmt.Printf("[AUTH] ERROR: Header doesn't start with 'Bearer ', got: '%s'\n", authHeader)
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)

		tokenPreview := token
		if len(tokenPreview) > 30 {
			tokenPreview = tokenPreview[:30]
		}
		fmt.Printf("[AUTH] Token extracted (length: %d): '%s...'\n", len(token), tokenPreview)

		userID, err := utils.VerifyAccessToken(token, secret)
		if err != nil {
			fmt.Printf("[AUTH] ERROR: Token verification failed: %v\n", err)
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		fmt.Printf("[AUTH] Token verified! UserID: %d\n", *userID)
		fmt.Printf("[AUTH] Setting UserID in context: %d\n", *userID)

		ctx := context.WithValue(r.Context(), UserIDKey, *userID)

		fmt.Println("[AUTH] Authentication successful, calling next handler")
		fmt.Println("========================================")

		next(w, r.WithContext(ctx))
	}
}
