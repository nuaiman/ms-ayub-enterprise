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
// EXPENSE MANAGEMENT
// =============================================================================

// CreateExpenseHandler - POST /api/expenses
func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[EXPENSES] CreateExpenseHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Title       string     `json:"title"`
		Amount      float64    `json:"amount"`
		ExpenseDate *time.Time `json:"expense_date,omitempty"`
		Notes       *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[EXPENSES] CreateExpenseHandler ERROR: invalid request body - %v", err)
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
		log.Printf("[EXPENSES] CreateExpenseHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default expense date if not provided
	expenseDate := req.ExpenseDate
	if expenseDate == nil {
		now := time.Now()
		expenseDate = &now
	}

	expense := &models.Expense{
		UserID:      userID,
		Title:       req.Title,
		Amount:      req.Amount,
		ExpenseDate: *expenseDate,
		Notes:       req.Notes,
	}

	id, err := h.app.Models.Expense.Insert(r.Context(), expense)
	if err != nil {
		log.Printf("[EXPENSES] CreateExpenseHandler ERROR: failed to insert expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create expense")
		return
	}

	expense.ID = id

	log.Printf("[EXPENSES] CreateExpenseHandler SUCCESS: created expense ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created expense: " + req.Title,
		EntityType:  "expenses",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "expense created successfully", expense)
}

// GetExpenseHandler - GET /api/expenses/{id}
func (h *Handler) GetExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[EXPENSES] GetExpenseHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[EXPENSES] GetExpenseHandler called - ID: %d", id)

	expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[EXPENSES] GetExpenseHandler ERROR: failed to fetch expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expense")
		return
	}

	if expense == nil {
		log.Printf("[EXPENSES] GetExpenseHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	log.Printf("[EXPENSES] GetExpenseHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "expense details fetched", expense)
}

// GetAllExpensesHandler - GET /api/expenses
func (h *Handler) GetAllExpensesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[EXPENSES] GetAllExpensesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var expenses []models.Expense
	var err error

	switch {
	case search != "":
		expenses, err = h.app.Models.Expense.Search(r.Context(), search)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		expenses, err = h.app.Models.Expense.GetByUserID(r.Context(), id)
	case startDate != "" && endDate != "":
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		expenses, err = h.app.Models.Expense.GetByDateRange(r.Context(), start, end)
	default:
		expenses, err = h.app.Models.Expense.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[EXPENSES] GetAllExpensesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expenses")
		return
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	log.Printf("[EXPENSES] GetAllExpensesHandler SUCCESS: fetched %d expenses", len(expenses))
	utils.SuccessJson(w, http.StatusOK, "expenses fetched successfully", expenses)
}

// UpdateExpenseHandler - PATCH /api/expenses/{id}
func (h *Handler) UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[EXPENSES] UpdateExpenseHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[EXPENSES] UpdateExpenseHandler called - ID: %d", id)

	type request struct {
		Title       *string    `json:"title,omitempty"`
		Amount      *float64   `json:"amount,omitempty"`
		ExpenseDate *time.Time `json:"expense_date,omitempty"`
		Notes       *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[EXPENSES] UpdateExpenseHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing expense
	expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[EXPENSES] UpdateExpenseHandler ERROR: failed to fetch expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expense")
		return
	}
	if expense == nil {
		log.Printf("[EXPENSES] UpdateExpenseHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	// Apply updates
	if req.Title != nil {
		if *req.Title == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		expense.Title = *req.Title
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
			return
		}
		expense.Amount = *req.Amount
	}
	if req.ExpenseDate != nil {
		expense.ExpenseDate = *req.ExpenseDate
	}
	if req.Notes != nil {
		expense.Notes = req.Notes
	}

	if err := h.app.Models.Expense.Update(r.Context(), expense); err != nil {
		log.Printf("[EXPENSES] UpdateExpenseHandler ERROR: failed to update expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update expense")
		return
	}

	log.Printf("[EXPENSES] UpdateExpenseHandler SUCCESS: updated expense ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated expense #" + strconv.FormatInt(id, 10),
		EntityType:  "expenses",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated expense
	updatedExpense, _ := h.app.Models.Expense.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "expense updated successfully", updatedExpense)
}

// DeleteExpenseHandler - DELETE /api/expenses/{id}
func (h *Handler) DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[EXPENSES] DeleteExpenseHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[EXPENSES] DeleteExpenseHandler called - ID: %d", id)

	// Check if expense exists
	expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[EXPENSES] DeleteExpenseHandler ERROR: failed to fetch expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expense")
		return
	}
	if expense == nil {
		log.Printf("[EXPENSES] DeleteExpenseHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	if err := h.app.Models.Expense.Delete(r.Context(), id); err != nil {
		log.Printf("[EXPENSES] DeleteExpenseHandler ERROR: failed to delete expense - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete expense")
		return
	}

	log.Printf("[EXPENSES] DeleteExpenseHandler SUCCESS: deleted expense ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted expense: " + expense.Title,
		EntityType:  "expenses",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "expense deleted successfully", nil)
}
