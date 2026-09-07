package middlewares

import (
	"backend/internal/app"
	"backend/internal/utils"
	"fmt"
	"net/http"
)

func Authorize(app *app.Application, roles ...string) func(http.HandlerFunc) http.HandlerFunc {

	fmt.Printf("[AUTHZ] AUTHORIZE SETUP - Allowed roles: %v\n", roles)

	allowed := make(map[string]struct{})

	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("========================================")
			fmt.Println("[AUTHZ] AUTHORIZE MIDDLEWARE CALLED")
			fmt.Printf("[AUTHZ] Method: %s, Path: %s\n", r.Method, r.URL.Path)
			fmt.Printf("[AUTHZ] Required roles: %v\n", roles)

			userID, ok := r.Context().Value(UserIDKey).(int64)

			fmt.Printf("[AUTHZ] UserID from context: %d, ok: %v\n", userID, ok)

			if !ok {
				fmt.Println("[AUTHZ] ERROR: UserID not found in context - authentication likely failed or wasn't called")
				utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
				return
			}

			fmt.Printf("[AUTHZ] Looking up user in database, ID: %d\n", userID)

			user, err := app.Models.User.GetByID(r.Context(), userID)

			if err != nil {
				fmt.Printf("[AUTHZ] ERROR: Database error fetching user %d: %v\n", userID, err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
				return
			}

			if user == nil {
				fmt.Printf("[AUTHZ] ERROR: User not found in database, ID: %d\n", userID)
				utils.ErrorJson(w, http.StatusUnauthorized, "user not found")
				return
			}

			fmt.Printf("[AUTHZ] User found: ID=%d, Username=%s, Role=%s, Active=%v\n",
				user.ID, user.Username, user.Role, user.IsActive)

			fmt.Printf("[AUTHZ] Checking role: User role '%s' against allowed %v\n", user.Role, roles)

			if _, exists := allowed[user.Role]; !exists {
				fmt.Printf("[AUTHZ] ERROR: Permission denied. User role '%s' not in allowed roles %v\n", user.Role, roles)
				utils.ErrorJson(w, http.StatusForbidden, "permission denied")
				return
			}

			fmt.Printf("[AUTHZ] Authorization successful for user %s (role: %s)\n", user.Username, user.Role)
			fmt.Println("========================================")

			next(w, r)
		}
	}
}
