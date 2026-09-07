package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"strconv"
)

// =============================================================================
// AUTHENTICATION
// =============================================================================

// LoginHandler - POST /api/users/login
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] LoginHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] LoginHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		log.Printf("[USERS] LoginHandler ERROR: missing required fields")
		utils.ErrorJson(w, http.StatusBadRequest, "username and password are required")
		return
	}

	log.Printf("[USERS] Login attempt for username: %s", req.Username)

	user, err := h.app.Models.User.GetByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("[USERS] LoginHandler ERROR: database error - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to authenticate")
		return
	}

	if user == nil {
		log.Printf("[USERS] LoginHandler ERROR: user not found - %s", req.Username)
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !user.IsActive {
		log.Printf("[USERS] LoginHandler ERROR: account deactivated - %s", req.Username)
		utils.ErrorJson(w, http.StatusForbidden, "account is deactivated")
		return
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		log.Printf("[USERS] LoginHandler ERROR: invalid password for user - %s", req.Username)
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Printf("[USERS] LoginHandler ERROR: failed to generate refresh token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate session")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), user.ID, &refreshToken); err != nil {
		log.Printf("[USERS] LoginHandler ERROR: failed to update refresh token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not save session")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, h.app.Config.JWTKey)
	if err != nil {
		log.Printf("[USERS] LoginHandler ERROR: failed to generate access token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
	})

	user.Password = ""

	log.Printf("[USERS] LoginHandler SUCCESS: user logged in - ID=%d, Username=%s", user.ID, user.Username)

	utils.SuccessJson(w, http.StatusOK, "login successful", map[string]interface{}{
		"access_token": accessToken,
		"user":         user,
	})
}

// RefreshHandler - POST /api/users/refresh-session
func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] RefreshHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		log.Printf("[USERS] RefreshHandler ERROR: missing refresh token")
		utils.ErrorJson(w, http.StatusUnauthorized, "missing session")
		return
	}

	user, err := h.app.Models.User.GetByRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		log.Printf("[USERS] RefreshHandler ERROR: invalid refresh token - %v", err)
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid session")
		return
	}

	if user == nil {
		log.Printf("[USERS] RefreshHandler ERROR: user not found for refresh token")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid session")
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Printf("[USERS] RefreshHandler ERROR: failed to generate refresh token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate session")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), user.ID, &newRefreshToken); err != nil {
		log.Printf("[USERS] RefreshHandler ERROR: failed to update refresh token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not update session")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, h.app.Config.JWTKey)
	if err != nil {
		log.Printf("[USERS] RefreshHandler ERROR: failed to generate access token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
	})

	log.Printf("[USERS] RefreshHandler SUCCESS: refreshed session for user ID=%d", user.ID)

	utils.SuccessJson(w, http.StatusOK, "session refreshed", map[string]interface{}{
		"access_token": accessToken,
	})
}

// LogoutHandler - DELETE /api/users/logout
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] LogoutHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[USERS] LogoutHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), userID, nil); err != nil {
		log.Printf("[USERS] LogoutHandler ERROR: failed to clear refresh token - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	log.Printf("[USERS] LogoutHandler SUCCESS: user logged out - ID=%d", userID)

	utils.SuccessJson(w, http.StatusOK, "logged out successfully", nil)
}

// =============================================================================
// CURRENT USER
// =============================================================================

// GetCurrentUserHandler - GET /api/users/current-user
func (h *Handler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] GetCurrentUserHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[USERS] GetCurrentUserHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] GetCurrentUserHandler ERROR: database error - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		log.Printf("[USERS] GetCurrentUserHandler ERROR: user not found - ID=%d", userID)
		utils.ErrorJson(w, http.StatusUnauthorized, "user not found")
		return
	}

	user.Password = ""

	log.Printf("[USERS] GetCurrentUserHandler SUCCESS: user ID=%d", userID)
	utils.SuccessJson(w, http.StatusOK, "user details fetched", user)
}

// =============================================================================
// PASSWORD MANAGEMENT
// =============================================================================

// ChangePasswordHandler - PATCH /api/users/change-password
func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] ChangePasswordHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		log.Printf("[USERS] ChangePasswordHandler ERROR: missing required fields")
		utils.ErrorJson(w, http.StatusBadRequest, "current_password and new_password are required")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: password validation failed - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.CurrentPassword == req.NewPassword {
		log.Printf("[USERS] ChangePasswordHandler ERROR: new password same as current")
		utils.ErrorJson(w, http.StatusBadRequest, "new password must be different from current password")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[USERS] ChangePasswordHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: user not found - ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if !utils.ComparePassword(user.Password, req.CurrentPassword) {
		log.Printf("[USERS] ChangePasswordHandler ERROR: incorrect current password for user ID=%d", userID)
		utils.ErrorJson(w, http.StatusUnauthorized, "current password is incorrect")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: failed to hash password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := h.app.Models.User.UpdatePassword(r.Context(), userID, hashedPassword); err != nil {
		log.Printf("[USERS] ChangePasswordHandler ERROR: failed to update password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	log.Printf("[USERS] ChangePasswordHandler SUCCESS: user ID=%d", userID)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Changed own password",
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "password changed successfully", nil)
}

// ChangeUserPasswordHandler - PATCH /api/users/{id}/change-password (Admin only)
func (h *Handler) ChangeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] ChangeUserPasswordHandler called - UserID: %d", userID)

	type request struct {
		NewPassword string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.NewPassword == "" {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: missing new_password")
		utils.ErrorJson(w, http.StatusBadRequest, "missing new_password")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: password validation failed - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: user not found - ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: cannot modify admin user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: failed to hash password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := h.app.Models.User.AdminUpdatePasswordAndInvalidateSessions(r.Context(), userID, hashed); err != nil {
		log.Printf("[USERS] ChangeUserPasswordHandler ERROR: failed to update password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	log.Printf("[USERS] ChangeUserPasswordHandler SUCCESS: changed password for user ID=%d", userID)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Admin changed password for user #" + strconv.FormatInt(userID, 10),
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user password updated successfully", nil)
}

// ResetAllPasswordsHandler - PATCH /api/users/reset-all-passwords (Admin only)
func (h *Handler) ResetAllPasswordsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] ResetAllPasswordsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		NewPassword string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] ResetAllPasswordsHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.NewPassword == "" {
		log.Printf("[USERS] ResetAllPasswordsHandler ERROR: missing new_password")
		utils.ErrorJson(w, http.StatusBadRequest, "missing new_password")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		log.Printf("[USERS] ResetAllPasswordsHandler ERROR: password validation failed - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[USERS] ResetAllPasswordsHandler ERROR: failed to hash password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := h.app.Models.User.ResetAllPasswordsExceptRoot(r.Context(), hashed); err != nil {
		log.Printf("[USERS] ResetAllPasswordsHandler ERROR: failed to reset passwords - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to reset passwords")
		return
	}

	log.Printf("[USERS] ResetAllPasswordsHandler SUCCESS: reset all passwords")

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Reset all user passwords",
		EntityType:  "users",
		EntityID:    0,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "all user passwords reset successfully", nil)
}

// =============================================================================
// USER MANAGEMENT (CRUD)
// =============================================================================

// CreateUserHandler - POST /api/users (Admin only)
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] CreateUserHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Name          string  `json:"name"`
		Username      string  `json:"username"`
		Password      string  `json:"password"`
		Role          string  `json:"role"`
		Email         *string `json:"email"`
		Phone         *string `json:"phone"`
		Address       *string `json:"address"`
		IDType        *string `json:"id_type"`
		IDNumber      *string `json:"id_number"`
		ImageURL      *string `json:"image_url"`
		MonthlySalary float64 `json:"monthly_salary"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] CreateUserHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Username == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "username is required")
		return
	}
	if req.Password == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "password is required")
		return
	}
	if req.Role == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "role is required")
		return
	}

	// Validate role - prevent creating admin users via API
	switch req.Role {
	case "manager", "staff", "accounts":
		// Valid roles
	case "admin":
		utils.ErrorJson(w, http.StatusForbidden, "cannot create admin users via API")
		return
	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid role")
		return
	}

	// Validate email if provided
	if req.Email != nil && *req.Email != "" {
		if !utils.IsValidEmail(*req.Email) {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}

		exists, err := h.app.Models.User.ExistsByEmail(r.Context(), *req.Email)
		if err != nil {
			log.Printf("[USERS] CreateUserHandler ERROR: failed to check email - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to check email")
			return
		}
		if exists {
			utils.ErrorJson(w, http.StatusConflict, "email already in use")
			return
		}
	}

	// Validate ID type if provided
	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if username already exists
	exists, err := h.app.Models.User.ExistsByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("[USERS] CreateUserHandler ERROR: failed to check username - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to check username")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "username already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("[USERS] CreateUserHandler ERROR: failed to hash password - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		Name:          req.Name,
		Username:      req.Username,
		Password:      hashedPassword,
		Role:          req.Role,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		IDType:        req.IDType,
		IDNumber:      req.IDNumber,
		ImageURL:      req.ImageURL,
		MonthlySalary: req.MonthlySalary,
	}

	id, err := h.app.Models.User.Insert(r.Context(), &user)
	if err != nil {
		log.Printf("[USERS] CreateUserHandler ERROR: failed to insert user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	user.ID = id
	user.Password = ""

	log.Printf("[USERS] CreateUserHandler SUCCESS: created user ID=%d, Username=%s", id, req.Username)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "create",
		Description: "Created user: " + req.Username,
		EntityType:  "users",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "user created successfully", user)
}

// GetUserHandler - GET /api/users/{id}
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] GetUserHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] GetUserHandler called - ID: %d", userID)

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] GetUserHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		log.Printf("[USERS] GetUserHandler NOT FOUND: ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	user.Password = ""

	log.Printf("[USERS] GetUserHandler SUCCESS: ID=%d", userID)
	utils.SuccessJson(w, http.StatusOK, "user details fetched", user)
}

// GetAllUsersHandler - GET /api/users
func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[USERS] GetAllUsersHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	role := r.URL.Query().Get("role")
	activeOnly := r.URL.Query().Get("active") == "true"

	var users []models.User
	var err error

	switch {
	case search != "":
		users, err = h.app.Models.User.Search(r.Context(), search)
	case role != "":
		users, err = h.app.Models.User.GetByRole(r.Context(), role)
	case activeOnly:
		users, err = h.app.Models.User.GetActiveUsers(r.Context())
	default:
		users, err = h.app.Models.User.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[USERS] GetAllUsersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	for i := range users {
		users[i].Password = ""
		users[i].RefreshToken = nil
	}

	log.Printf("[USERS] GetAllUsersHandler SUCCESS: fetched %d users", len(users))
	utils.SuccessJson(w, http.StatusOK, "users fetched successfully", users)
}

// UpdateUserProfileHandler - PATCH /api/users/{id}/profile
func (h *Handler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	targetUserID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] UpdateUserProfileHandler called - TargetID: %d", targetUserID)

	type request struct {
		Name     string  `json:"name"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
		Address  *string `json:"address"`
		IDType   *string `json:"id_type"`
		IDNumber *string `json:"id_number"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	currentUserID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	targetUser, err := h.app.Models.User.GetByID(r.Context(), targetUserID)
	if err != nil {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: failed to fetch target user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if targetUser == nil {
		log.Printf("[USERS] UpdateUserProfileHandler NOT FOUND: target user ID=%d", targetUserID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	// Check if updating own profile or has admin/manager role
	if targetUserID != currentUserID {
		currentUser, err := h.app.Models.User.GetByID(r.Context(), currentUserID)
		if err != nil || currentUser == nil {
			log.Printf("[USERS] UpdateUserProfileHandler ERROR: failed to verify permissions - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify permissions")
			return
		}
		if currentUser.Role != "admin" && currentUser.Role != "manager" {
			log.Printf("[USERS] UpdateUserProfileHandler ERROR: permission denied - UserID=%d, Role=%s",
				currentUserID, currentUser.Role)
			utils.ErrorJson(w, http.StatusForbidden, "you can only update your own profile")
			return
		}
	}

	// Prevent modifying admin users (except root)
	if targetUser.Role == "admin" && targetUserID != currentUserID {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: cannot modify admin user - TargetID=%d, CurrentID=%d",
			targetUserID, currentUserID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	// Validate email if provided
	if req.Email != nil && *req.Email != "" {
		if !utils.IsValidEmail(*req.Email) {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}

		existingUser, err := h.app.Models.User.GetByEmail(r.Context(), *req.Email)
		if err == nil && existingUser != nil && existingUser.ID != targetUserID {
			utils.ErrorJson(w, http.StatusConflict, "email already in use")
			return
		}
	}

	// Validate ID type if provided
	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
	}

	if err := h.app.Models.User.UpdateProfile(r.Context(), targetUserID, user); err != nil {
		log.Printf("[USERS] UpdateUserProfileHandler ERROR: failed to update - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	log.Printf("[USERS] UpdateUserProfileHandler SUCCESS: user ID=%d", targetUserID)

	updatedUser, err := h.app.Models.User.GetByID(r.Context(), targetUserID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated user")
		return
	}

	updatedUser.Password = ""

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Updated profile for user #" + strconv.FormatInt(targetUserID, 10),
		EntityType:  "users",
		EntityID:    targetUserID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user updated successfully", updatedUser)
}

// ChangeUserRoleHandler - PATCH /api/users/{id}/change-role (Admin only)
func (h *Handler) ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] ChangeUserRoleHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] ChangeUserRoleHandler called - UserID: %d", userID)

	type request struct {
		Role string `json:"role"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] ChangeUserRoleHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "role is required")
		return
	}

	validRoles := map[string]bool{
		"admin": true, "manager": true, "accounts": true, "staff": true,
	}
	if !validRoles[req.Role] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid role")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] ChangeUserRoleHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		log.Printf("[USERS] ChangeUserRoleHandler NOT FOUND: user ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		log.Printf("[USERS] ChangeUserRoleHandler ERROR: cannot modify admin user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	if err := h.app.Models.User.UpdateRole(r.Context(), userID, req.Role); err != nil {
		log.Printf("[USERS] ChangeUserRoleHandler ERROR: failed to update role - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update role")
		return
	}

	log.Printf("[USERS] ChangeUserRoleHandler SUCCESS: changed role for user ID=%d to '%s'", userID, req.Role)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Changed role for user #" + strconv.FormatInt(userID, 10) + " to " + req.Role,
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user role updated successfully", nil)
}

// ToggleUserActiveHandler - PATCH /api/users/{id}/toggle-active (Admin only)
func (h *Handler) ToggleUserActiveHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] ToggleUserActiveHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] ToggleUserActiveHandler called - UserID: %d", userID)

	if userID == 1 {
		log.Printf("[USERS] ToggleUserActiveHandler ERROR: cannot modify root user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify root user")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] ToggleUserActiveHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		log.Printf("[USERS] ToggleUserActiveHandler NOT FOUND: user ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		log.Printf("[USERS] ToggleUserActiveHandler ERROR: cannot modify admin user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	newStatus := !user.IsActive

	if err := h.app.Models.User.SetActiveStatus(r.Context(), userID, newStatus); err != nil {
		log.Printf("[USERS] ToggleUserActiveHandler ERROR: failed to update status - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update user status")
		return
	}

	user.IsActive = newStatus
	log.Printf("[USERS] ToggleUserActiveHandler SUCCESS: user ID=%d now active=%v", userID, newStatus)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	action := "deactivated"
	if newStatus {
		action = "activated"
	}
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "User " + action + ": #" + strconv.FormatInt(userID, 10),
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	user.Password = ""
	utils.SuccessJson(w, http.StatusOK, "user status updated successfully", user)
}

// UpdateUserSalaryHandler - PATCH /api/users/{id}/salary (Admin only)
func (h *Handler) UpdateUserSalaryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] UpdateUserSalaryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] UpdateUserSalaryHandler called - UserID: %d", userID)

	type request struct {
		MonthlySalary float64 `json:"monthly_salary"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[USERS] UpdateUserSalaryHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.MonthlySalary < 0 {
		log.Printf("[USERS] UpdateUserSalaryHandler ERROR: salary cannot be negative - %.2f", req.MonthlySalary)
		utils.ErrorJson(w, http.StatusBadRequest, "monthly_salary cannot be negative")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] UpdateUserSalaryHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		log.Printf("[USERS] UpdateUserSalaryHandler NOT FOUND: user ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if err := h.app.Models.User.UpdateMonthlySalary(r.Context(), userID, req.MonthlySalary); err != nil {
		log.Printf("[USERS] UpdateUserSalaryHandler ERROR: failed to update salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update salary")
		return
	}

	log.Printf("[USERS] UpdateUserSalaryHandler SUCCESS: updated salary for user ID=%d to %.2f",
		userID, req.MonthlySalary)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Updated salary for user #" + strconv.FormatInt(userID, 10) + " to " + strconv.FormatFloat(req.MonthlySalary, 'f', 2, 64),
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	user.MonthlySalary = req.MonthlySalary
	user.Password = ""

	utils.SuccessJson(w, http.StatusOK, "salary updated successfully", user)
}

// DeleteUserHandler - DELETE /api/users/{id} (Admin only)
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[USERS] DeleteUserHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[USERS] DeleteUserHandler called - UserID: %d", userID)

	if userID == 1 {
		log.Printf("[USERS] DeleteUserHandler ERROR: cannot delete root user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot delete root user")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("[USERS] DeleteUserHandler ERROR: failed to fetch user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		log.Printf("[USERS] DeleteUserHandler NOT FOUND: user ID=%d", userID)
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		log.Printf("[USERS] DeleteUserHandler ERROR: cannot delete admin user - ID=%d", userID)
		utils.ErrorJson(w, http.StatusForbidden, "cannot delete admin users")
		return
	}

	if err := h.app.Models.User.Delete(r.Context(), userID); err != nil {
		log.Printf("[USERS] DeleteUserHandler ERROR: failed to delete user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	log.Printf("[USERS] DeleteUserHandler SUCCESS: deleted user ID=%d", userID)

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      currentUserID,
		Action:      "delete",
		Description: "Deleted user: " + user.Username,
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user deleted successfully", nil)
}
