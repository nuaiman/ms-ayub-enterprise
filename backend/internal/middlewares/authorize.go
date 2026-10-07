package middlewares

import (
	"log"
	"net/http"

	"backend/internal/app"
	"backend/internal/utils"
)

func Authorize(app *app.Application, roles ...string) func(http.HandlerFunc) http.HandlerFunc {

	allowed := make(map[string]struct{})
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int64)
			if !ok {
				log.Printf("[AUTHZ] UserID not in context for %s %s", r.Method, r.URL.Path)
				utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
				return
			}

			user, err := app.Models.User.GetByID(r.Context(), userID)
			if err != nil {
				log.Printf("[AUTHZ] DB error fetching user %d: %v", userID, err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
				return
			}
			if user == nil {
				log.Printf("[AUTHZ] User not found: %d", userID)
				utils.ErrorJson(w, http.StatusUnauthorized, "user not found")
				return
			}

			if _, exists := allowed[user.Role]; !exists {
				log.Printf("[AUTHZ] Permission denied for user %s (role=%s) on %s %s",
					user.Username, user.Role, r.Method, r.URL.Path)
				utils.ErrorJson(w, http.StatusForbidden, "permission denied")
				return
			}

			next(w, r)
		}
	}
}
