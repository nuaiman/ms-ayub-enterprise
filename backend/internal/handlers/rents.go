package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// =============================================================================
// RENT MANAGEMENT
// =============================================================================

// CreateRentHandler - POST /api/rents
func (h *Handler) CreateRentHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RENTS] CreateRentHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		GodownID  int64   `json:"godown_id"`
		MonthYear string  `json:"month_year"`
		Amount    float64 `json:"amount"`
		Notes     *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[RENTS] CreateRentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.GodownID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "godown_id is required")
		return
	}
	if req.MonthYear == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "month_year is required")
		return
	}
	if req.Amount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
		return
	}

	// Validate month_year format (YYYY-MM)
	if len(req.MonthYear) != 7 || req.MonthYear[4] != '-' {
		utils.ErrorJson(w, http.StatusBadRequest, "month_year must be in YYYY-MM format")
		return
	}

	// Validate: Cannot create rent for future months
	currentMonth := time.Now().Format("2006-01")
	if req.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot create rent for future months")
		return
	}

	// Check if godown exists and is active
	godown, err := h.app.Models.Godown.GetByID(r.Context(), req.GodownID)
	if err != nil {
		log.Printf("[RENTS] CreateRentHandler ERROR: failed to fetch godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
		return
	}
	if godown == nil {
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}
	if !godown.IsActive {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot create rent for inactive godown")
		return
	}

	// Check if rent already exists for this godown and month
	exists, err := h.app.Models.Rent.ExistsForGodownMonth(r.Context(), req.GodownID, req.MonthYear)
	if err != nil {
		log.Printf("[RENTS] CreateRentHandler ERROR: failed to check existing rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify existing rent")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "rent already exists for this godown and month")
		return
	}

	// Get current user ID
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[RENTS] CreateRentHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	rent := &models.Rent{
		UserID:    userID,
		GodownID:  req.GodownID,
		MonthYear: req.MonthYear,
		Amount:    req.Amount,
		Notes:     req.Notes,
	}

	id, err := h.app.Models.Rent.Insert(r.Context(), rent)
	if err != nil {
		log.Printf("[RENTS] CreateRentHandler ERROR: failed to insert rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create rent")
		return
	}

	rent.ID = id

	log.Printf("[RENTS] CreateRentHandler SUCCESS: created rent ID=%d for godown ID=%d", id, req.GodownID)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created rent for godown #" + strconv.FormatInt(req.GodownID, 10) + " for " + req.MonthYear,
		EntityType:  "rents",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "rent created successfully", rent)
}

// GetRentHandler - GET /api/rents/{id}
func (h *Handler) GetRentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[RENTS] GetRentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[RENTS] GetRentHandler called - ID: %d", id)

	rent, err := h.app.Models.Rent.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[RENTS] GetRentHandler ERROR: failed to fetch rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rent")
		return
	}

	if rent == nil {
		log.Printf("[RENTS] GetRentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "rent not found")
		return
	}

	log.Printf("[RENTS] GetRentHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "rent details fetched", rent)
}

// GetAllRentsHandler - GET /api/rents
func (h *Handler) GetAllRentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RENTS] GetAllRentsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	godownID := r.URL.Query().Get("godown_id")
	monthYear := r.URL.Query().Get("month_year")
	status := r.URL.Query().Get("status")
	outstandingOnly := r.URL.Query().Get("outstanding") == "true"

	var rents []models.Rent
	var err error

	switch {
	case outstandingOnly:
		rents, err = h.app.Models.Rent.GetOutstanding(r.Context())
	case godownID != "":
		id, _ := strconv.ParseInt(godownID, 10, 64)
		rents, err = h.app.Models.Rent.GetByGodownID(r.Context(), id)
	case monthYear != "":
		rents, err = h.app.Models.Rent.GetByMonth(r.Context(), monthYear)
	case status != "":
		rents, err = h.app.Models.Rent.GetByStatus(r.Context(), status)
	default:
		rents, err = h.app.Models.Rent.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[RENTS] GetAllRentsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rents")
		return
	}

	if rents == nil {
		rents = []models.Rent{}
	}

	log.Printf("[RENTS] GetAllRentsHandler SUCCESS: fetched %d rents", len(rents))
	utils.SuccessJson(w, http.StatusOK, "rents fetched successfully", rents)
}

// UpdateRentHandler - PATCH /api/rents/{id}
func (h *Handler) UpdateRentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[RENTS] UpdateRentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[RENTS] UpdateRentHandler called - ID: %d", id)

	type request struct {
		Amount *float64 `json:"amount,omitempty"`
		Notes  *string  `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[RENTS] UpdateRentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing rent
	rent, err := h.app.Models.Rent.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[RENTS] UpdateRentHandler ERROR: failed to fetch rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rent")
		return
	}
	if rent == nil {
		log.Printf("[RENTS] UpdateRentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "rent not found")
		return
	}

	// Only draft rents can be updated
	if rent.Status != "draft" {
		log.Printf("[RENTS] UpdateRentHandler ERROR: cannot update non-draft rent - Status: %s", rent.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft rents can be updated")
		return
	}

	// Validate: Cannot update future months
	currentMonth := time.Now().Format("2006-01")
	if rent.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot update rent for future months")
		return
	}

	// Apply updates
	if req.Amount != nil {
		if *req.Amount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
			return
		}
		rent.Amount = *req.Amount
	}
	if req.Notes != nil {
		rent.Notes = req.Notes
	}

	if err := h.app.Models.Rent.Update(r.Context(), rent); err != nil {
		log.Printf("[RENTS] UpdateRentHandler ERROR: failed to update rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update rent")
		return
	}

	log.Printf("[RENTS] UpdateRentHandler SUCCESS: updated rent ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated rent #" + strconv.FormatInt(id, 10),
		EntityType:  "rents",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated rent
	updatedRent, _ := h.app.Models.Rent.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "rent updated successfully", updatedRent)
}

// MarkRentAsPaidHandler - PATCH /api/rents/{id}/pay
func (h *Handler) MarkRentAsPaidHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[RENTS] MarkRentAsPaidHandler called - ID: %d", id)

	type request struct {
		PaymentDate     *string `json:"payment_date,omitempty"`
		PaymentMethod   *string `json:"payment_method,omitempty"`
		ReferenceNumber *string `json:"reference_number,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate payment method if provided
	if req.PaymentMethod != nil && *req.PaymentMethod != "" {
		validMethods := map[string]bool{
			"cash": true, "bank_transfer": true, "check": true, "mobile_banking": true,
		}
		if !validMethods[*req.PaymentMethod] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid payment_method. Valid: cash, bank_transfer, check, mobile_banking")
			return
		}
	}

	// Get existing rent
	rent, err := h.app.Models.Rent.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: failed to fetch rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rent")
		return
	}
	if rent == nil {
		log.Printf("[RENTS] MarkRentAsPaidHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "rent not found")
		return
	}

	// Only draft rents can be marked as paid
	if rent.Status != "draft" {
		log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: rent not draft - Status: %s", rent.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft rents can be marked as paid")
		return
	}

	// Validate: Cannot pay future months
	currentMonth := time.Now().Format("2006-01")
	if rent.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot pay rent for future months")
		return
	}

	// Parse payment date
	var paymentDate *time.Time
	if req.PaymentDate != nil && *req.PaymentDate != "" {
		layouts := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
		}
		var parsed time.Time
		var parseErr error
		for _, layout := range layouts {
			parsed, parseErr = time.Parse(layout, *req.PaymentDate)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: failed to parse payment date - %v", parseErr)
			utils.ErrorJson(w, http.StatusBadRequest, "invalid payment_date format. Use YYYY-MM-DD or RFC3339")
			return
		}
		paymentDate = &parsed
	}

	if paymentDate == nil {
		now := time.Now()
		paymentDate = &now
	}

	if err := h.app.Models.Rent.MarkAsPaid(r.Context(), id, paymentDate, req.PaymentMethod, req.ReferenceNumber); err != nil {
		log.Printf("[RENTS] MarkRentAsPaidHandler ERROR: failed to mark as paid - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to mark rent as paid")
		return
	}

	log.Printf("[RENTS] MarkRentAsPaidHandler SUCCESS: rent ID=%d marked as paid", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Marked rent #" + strconv.FormatInt(id, 10) + " as paid",
		EntityType:  "rents",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated rent
	updatedRent, _ := h.app.Models.Rent.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "rent marked as paid successfully", updatedRent)
}

// CancelRentHandler - PATCH /api/rents/{id}/cancel
func (h *Handler) CancelRentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[RENTS] CancelRentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[RENTS] CancelRentHandler called - ID: %d", id)

	// Get existing rent
	rent, err := h.app.Models.Rent.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[RENTS] CancelRentHandler ERROR: failed to fetch rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rent")
		return
	}
	if rent == nil {
		log.Printf("[RENTS] CancelRentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "rent not found")
		return
	}

	// Cannot cancel already cancelled rents
	if rent.Status == "cancelled" {
		utils.ErrorJson(w, http.StatusBadRequest, "rent is already cancelled")
		return
	}

	// Validate: Cannot cancel future months
	currentMonth := time.Now().Format("2006-01")
	if rent.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot cancel rent for future months")
		return
	}

	if err := h.app.Models.Rent.Cancel(r.Context(), id); err != nil {
		log.Printf("[RENTS] CancelRentHandler ERROR: failed to cancel rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to cancel rent")
		return
	}

	log.Printf("[RENTS] CancelRentHandler SUCCESS: rent ID=%d cancelled", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Cancelled rent #" + strconv.FormatInt(id, 10),
		EntityType:  "rents",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated rent
	updatedRent, _ := h.app.Models.Rent.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "rent cancelled successfully", updatedRent)
}

// DeleteRentHandler - DELETE /api/rents/{id}
func (h *Handler) DeleteRentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[RENTS] DeleteRentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[RENTS] DeleteRentHandler called - ID: %d", id)

	// Get existing rent
	rent, err := h.app.Models.Rent.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[RENTS] DeleteRentHandler ERROR: failed to fetch rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rent")
		return
	}
	if rent == nil {
		log.Printf("[RENTS] DeleteRentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "rent not found")
		return
	}

	// Only draft rents can be deleted
	if rent.Status != "draft" {
		log.Printf("[RENTS] DeleteRentHandler ERROR: cannot delete non-draft rent - Status: %s", rent.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft rents can be deleted")
		return
	}

	// Validate: Cannot delete future months
	currentMonth := time.Now().Format("2006-01")
	if rent.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete rent for future months")
		return
	}

	if err := h.app.Models.Rent.Delete(r.Context(), id); err != nil {
		log.Printf("[RENTS] DeleteRentHandler ERROR: failed to delete rent - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete rent")
		return
	}

	log.Printf("[RENTS] DeleteRentHandler SUCCESS: deleted rent ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted rent #" + strconv.FormatInt(id, 10),
		EntityType:  "rents",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "rent deleted successfully", nil)
}

// GetCurrentMonthRentsHandler - GET /api/rents/current-month
func (h *Handler) GetCurrentMonthRentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RENTS] GetCurrentMonthRentsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	currentMonth := time.Now().Format("2006-01")
	ctx := r.Context()

	// Get current user ID for creating draft rents
	userID, ok := ctx.Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[RENTS] GetCurrentMonthRentsHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Ensure all active godowns have rent records for current month
	if err := h.app.Models.Rent.EnsureCurrentMonthRents(ctx, userID); err != nil {
		log.Printf("[RENTS] GetCurrentMonthRentsHandler ERROR: failed to ensure rents - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to prepare rents")
		return
	}

	// Fetch all rents for current month
	rents, err := h.app.Models.Rent.GetByMonth(ctx, currentMonth)
	if err != nil {
		log.Printf("[RENTS] GetCurrentMonthRentsHandler ERROR: failed to fetch rents - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch rents")
		return
	}

	if rents == nil {
		rents = []models.Rent{}
	}

	log.Printf("[RENTS] GetCurrentMonthRentsHandler SUCCESS: fetched %d rents", len(rents))
	utils.SuccessJson(w, http.StatusOK, "current month rents fetched", rents)
}
