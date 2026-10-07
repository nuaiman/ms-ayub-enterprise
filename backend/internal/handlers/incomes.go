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
// INCOME MANAGEMENT
// =============================================================================

// CreateIncomeHandler - POST /api/incomes
func (h *Handler) CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INCOMES] CreateIncomeHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Title      string     `json:"title"`
		Amount     float64    `json:"amount"`
		IncomeDate *time.Time `json:"income_date,omitempty"`
		Notes      *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[INCOMES] CreateIncomeHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Amount <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	// Get current user ID
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[INCOMES] CreateIncomeHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default income date if not provided
	incomeDate := req.IncomeDate
	if incomeDate == nil {
		now := time.Now()
		incomeDate = &now
	}

	income := &models.Income{
		UserID:     userID,
		Title:      req.Title,
		Amount:     req.Amount,
		IncomeDate: *incomeDate,
		Notes:      req.Notes,
	}

	id, err := h.app.Models.Income.Insert(r.Context(), income)
	if err != nil {
		log.Printf("[INCOMES] CreateIncomeHandler ERROR: failed to insert income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create income")
		return
	}

	income.ID = id

	log.Printf("[INCOMES] CreateIncomeHandler SUCCESS: created income ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created income: " + req.Title,
		EntityType:  "incomes",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "income created successfully", income)
}

// GetIncomeHandler - GET /api/incomes/{id}
func (h *Handler) GetIncomeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[INCOMES] GetIncomeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[INCOMES] GetIncomeHandler called - ID: %d", id)

	income, err := h.app.Models.Income.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[INCOMES] GetIncomeHandler ERROR: failed to fetch income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch income")
		return
	}

	if income == nil {
		log.Printf("[INCOMES] GetIncomeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "income not found")
		return
	}

	log.Printf("[INCOMES] GetIncomeHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "income details fetched", income)
}

// GetAllIncomesHandler - GET /api/incomes
func (h *Handler) GetAllIncomesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INCOMES] GetAllIncomesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var incomes []models.Income
	var err error

	switch {
	case search != "":
		incomes, err = h.app.Models.Income.Search(r.Context(), search)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		incomes, err = h.app.Models.Income.GetByUserID(r.Context(), id)
	case startDate != "" && endDate != "":
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		incomes, err = h.app.Models.Income.GetByDateRange(r.Context(), start, end)
	default:
		incomes, err = h.app.Models.Income.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[INCOMES] GetAllIncomesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch incomes")
		return
	}

	if incomes == nil {
		incomes = []models.Income{}
	}

	log.Printf("[INCOMES] GetAllIncomesHandler SUCCESS: fetched %d incomes", len(incomes))
	utils.SuccessJson(w, http.StatusOK, "incomes fetched successfully", incomes)
}

// UpdateIncomeHandler - PATCH /api/incomes/{id}
func (h *Handler) UpdateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[INCOMES] UpdateIncomeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[INCOMES] UpdateIncomeHandler called - ID: %d", id)

	type request struct {
		Title      *string    `json:"title,omitempty"`
		Amount     *float64   `json:"amount,omitempty"`
		IncomeDate *time.Time `json:"income_date,omitempty"`
		Notes      *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[INCOMES] UpdateIncomeHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing income
	income, err := h.app.Models.Income.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[INCOMES] UpdateIncomeHandler ERROR: failed to fetch income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch income")
		return
	}
	if income == nil {
		log.Printf("[INCOMES] UpdateIncomeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "income not found")
		return
	}

	// Apply updates
	if req.Title != nil {
		if *req.Title == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		income.Title = *req.Title
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
			return
		}
		income.Amount = *req.Amount
	}
	if req.IncomeDate != nil {
		income.IncomeDate = *req.IncomeDate
	}
	if req.Notes != nil {
		income.Notes = req.Notes
	}

	if err := h.app.Models.Income.Update(r.Context(), income); err != nil {
		log.Printf("[INCOMES] UpdateIncomeHandler ERROR: failed to update income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update income")
		return
	}

	log.Printf("[INCOMES] UpdateIncomeHandler SUCCESS: updated income ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated income #" + strconv.FormatInt(id, 10),
		EntityType:  "incomes",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated income
	updatedIncome, _ := h.app.Models.Income.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "income updated successfully", updatedIncome)
}

// DeleteIncomeHandler - DELETE /api/incomes/{id}
func (h *Handler) DeleteIncomeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[INCOMES] DeleteIncomeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[INCOMES] DeleteIncomeHandler called - ID: %d", id)

	// Check if income exists
	income, err := h.app.Models.Income.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[INCOMES] DeleteIncomeHandler ERROR: failed to fetch income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch income")
		return
	}
	if income == nil {
		log.Printf("[INCOMES] DeleteIncomeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "income not found")
		return
	}

	if err := h.app.Models.Income.Delete(r.Context(), id); err != nil {
		log.Printf("[INCOMES] DeleteIncomeHandler ERROR: failed to delete income - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete income")
		return
	}

	log.Printf("[INCOMES] DeleteIncomeHandler SUCCESS: deleted income ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted income: " + income.Title,
		EntityType:  "incomes",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "income deleted successfully", nil)
}
