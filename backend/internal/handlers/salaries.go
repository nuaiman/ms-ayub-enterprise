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
// SALARY MANAGEMENT
// =============================================================================

// CreateSalaryHandler - POST /api/salaries
func (h *Handler) CreateSalaryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SALARIES] CreateSalaryHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		EmployeeID int64   `json:"employee_id"`
		MonthYear  string  `json:"month_year"`
		Bonus      float64 `json:"bonus"`
		Deductions float64 `json:"deductions"`
		Notes      *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[SALARIES] CreateSalaryHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.EmployeeID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "employee_id is required")
		return
	}
	if req.MonthYear == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "month_year is required")
		return
	}

	// Validate month_year format (YYYY-MM)
	if len(req.MonthYear) != 7 || req.MonthYear[4] != '-' {
		utils.ErrorJson(w, http.StatusBadRequest, "month_year must be in YYYY-MM format")
		return
	}

	// Validate: Cannot create salary for future months
	currentMonth := time.Now().Format("2006-01")
	if req.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot create salary for future months")
		return
	}

	// Check if employee exists
	employee, err := h.app.Models.User.GetByID(r.Context(), req.EmployeeID)
	if err != nil {
		log.Printf("[SALARIES] CreateSalaryHandler ERROR: failed to fetch employee - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify employee")
		return
	}
	if employee == nil {
		utils.ErrorJson(w, http.StatusNotFound, "employee not found")
		return
	}

	// Check if salary already exists for this employee and month
	exists, err := h.app.Models.Salary.ExistsForEmployeeMonth(r.Context(), req.EmployeeID, req.MonthYear)
	if err != nil {
		log.Printf("[SALARIES] CreateSalaryHandler ERROR: failed to check existing salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify existing salary")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "salary already exists for this employee and month")
		return
	}

	// Get current user ID
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[SALARIES] CreateSalaryHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	salary := &models.Salary{
		UserID:     userID,
		EmployeeID: req.EmployeeID,
		MonthYear:  req.MonthYear,
		Bonus:      req.Bonus,
		Deductions: req.Deductions,
		Notes:      req.Notes,
	}

	id, err := h.app.Models.Salary.Insert(r.Context(), salary)
	if err != nil {
		log.Printf("[SALARIES] CreateSalaryHandler ERROR: failed to insert salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create salary")
		return
	}

	salary.ID = id

	log.Printf("[SALARIES] CreateSalaryHandler SUCCESS: created salary ID=%d for employee ID=%d", id, req.EmployeeID)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created salary for employee #" + strconv.FormatInt(req.EmployeeID, 10) + " for " + req.MonthYear,
		EntityType:  "salaries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "salary created successfully", salary)
}

// GetSalaryHandler - GET /api/salaries/{id}
func (h *Handler) GetSalaryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] GetSalaryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] GetSalaryHandler called - ID: %d", id)

	salary, err := h.app.Models.Salary.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[SALARIES] GetSalaryHandler ERROR: failed to fetch salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salary")
		return
	}

	if salary == nil {
		log.Printf("[SALARIES] GetSalaryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "salary not found")
		return
	}

	// Get employee's monthly salary for basic_salary
	employee, err := h.app.Models.User.GetByID(r.Context(), salary.EmployeeID)
	basicSalary := 0.0
	if err == nil && employee != nil {
		basicSalary = employee.MonthlySalary
	}

	type SalaryResponse struct {
		models.Salary
		BasicSalary float64 `json:"basic_salary"`
		TotalSalary float64 `json:"total_salary"`
	}

	response := SalaryResponse{
		Salary:      *salary,
		BasicSalary: basicSalary,
		TotalSalary: basicSalary + salary.Bonus - salary.Deductions,
	}

	log.Printf("[SALARIES] GetSalaryHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "salary details fetched", response)
}

// GetAllSalariesHandler - GET /api/salaries
func (h *Handler) GetAllSalariesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SALARIES] GetAllSalariesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	// Build filters from query params
	filters := make(map[string]string)
	monthYear := r.URL.Query().Get("month_year")

	if employeeID := r.URL.Query().Get("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if monthYear != "" {
		filters["month_year"] = monthYear
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	ctx := r.Context()
	currentMonth := time.Now().Format("2006-01")

	// If no month filter or filtering for current month, ensure all employees have salaries
	if monthYear == "" || monthYear == currentMonth {
		userID, ok := ctx.Value(middlewares.UserIDKey).(int64)
		if !ok {
			log.Printf("[SALARIES] GetAllSalariesHandler ERROR: invalid user context")
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
			return
		}

		err := h.app.Models.Salary.EnsureCurrentMonthSalaries(ctx, userID)
		if err != nil {
			log.Printf("[SALARIES] GetAllSalariesHandler ERROR: failed to ensure salaries - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to prepare salaries")
			return
		}
	}

	// Now fetch all salaries with the filters
	salaries, err := h.app.Models.Salary.GetAll(ctx, filters)
	if err != nil {
		log.Printf("[SALARIES] GetAllSalariesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salaries")
		return
	}

	// Build response with basic_salary and total_salary
	type SalaryResponse struct {
		models.Salary
		BasicSalary float64 `json:"basic_salary"`
		TotalSalary float64 `json:"total_salary"`
	}

	response := []SalaryResponse{}
	for _, salary := range salaries {
		employee, err := h.app.Models.User.GetByID(ctx, salary.EmployeeID)
		basicSalary := 0.0
		if err == nil && employee != nil {
			basicSalary = employee.MonthlySalary
		}

		totalSalary := basicSalary + salary.Bonus - salary.Deductions

		response = append(response, SalaryResponse{
			Salary:      salary,
			BasicSalary: basicSalary,
			TotalSalary: totalSalary,
		})
	}

	log.Printf("[SALARIES] GetAllSalariesHandler SUCCESS: fetched %d salaries", len(response))
	utils.SuccessJson(w, http.StatusOK, "salaries fetched successfully", response)
}

// UpdateSalaryHandler - PATCH /api/salaries/{id}
func (h *Handler) UpdateSalaryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] UpdateSalaryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] UpdateSalaryHandler called - ID: %d", id)

	type request struct {
		Bonus      *float64 `json:"bonus,omitempty"`
		Deductions *float64 `json:"deductions,omitempty"`
		Notes      *string  `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[SALARIES] UpdateSalaryHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing salary
	salary, err := h.app.Models.Salary.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[SALARIES] UpdateSalaryHandler ERROR: failed to fetch salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salary")
		return
	}
	if salary == nil {
		log.Printf("[SALARIES] UpdateSalaryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "salary not found")
		return
	}

	// Only draft salaries can be updated
	if salary.Status != "draft" {
		log.Printf("[SALARIES] UpdateSalaryHandler ERROR: cannot update non-draft salary - Status: %s", salary.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft salaries can be updated")
		return
	}

	// Validate: Cannot update future months
	currentMonth := time.Now().Format("2006-01")
	if salary.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot update salary for future months")
		return
	}

	// Apply updates
	if req.Bonus != nil {
		if *req.Bonus < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "bonus cannot be negative")
			return
		}
		salary.Bonus = *req.Bonus
	}
	if req.Deductions != nil {
		if *req.Deductions < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "deductions cannot be negative")
			return
		}
		salary.Deductions = *req.Deductions
	}
	if req.Notes != nil {
		salary.Notes = req.Notes
	}

	if err := h.app.Models.Salary.Update(r.Context(), salary); err != nil {
		log.Printf("[SALARIES] UpdateSalaryHandler ERROR: failed to update salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update salary")
		return
	}

	log.Printf("[SALARIES] UpdateSalaryHandler SUCCESS: updated salary ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated salary #" + strconv.FormatInt(id, 10),
		EntityType:  "salaries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated salary
	updatedSalary, _ := h.app.Models.Salary.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "salary updated successfully", updatedSalary)
}

// MarkSalaryAsPaidHandler - PATCH /api/salaries/{id}/pay
// Marks a draft salary as paid with payment details
func (h *Handler) MarkSalaryAsPaidHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] MarkSalaryAsPaidHandler called - ID: %d", id)

	type request struct {
		PaymentDate     *time.Time `json:"payment_date,omitempty"`
		PaymentMethod   *string    `json:"payment_method,omitempty"`
		ReferenceNumber *string    `json:"reference_number,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate payment method if provided
	if req.PaymentMethod != nil && *req.PaymentMethod != "" {
		validMethods := map[string]bool{
			"cash":           true,
			"bank_transfer":  true,
			"check":          true,
			"mobile_banking": true,
		}
		if !validMethods[*req.PaymentMethod] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid payment_method. Valid: cash, bank_transfer, check, mobile_banking")
			return
		}
	}

	// Get existing salary
	salary, err := h.app.Models.Salary.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler ERROR: failed to fetch salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salary")
		return
	}
	if salary == nil {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "salary not found")
		return
	}

	// Only draft salaries can be marked as paid
	if salary.Status != "draft" {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler ERROR: salary not draft - Status: %s", salary.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft salaries can be marked as paid")
		return
	}

	// Validate: Cannot pay future months
	currentMonth := time.Now().Format("2006-01")
	if salary.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot pay salary for future months")
		return
	}

	// Set default payment date if not provided
	if req.PaymentDate == nil {
		now := time.Now()
		req.PaymentDate = &now
	}

	if err := h.app.Models.Salary.MarkAsPaid(r.Context(), id, req.PaymentDate, req.PaymentMethod, req.ReferenceNumber); err != nil {
		log.Printf("[SALARIES] MarkSalaryAsPaidHandler ERROR: failed to mark as paid - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to mark salary as paid")
		return
	}

	log.Printf("[SALARIES] MarkSalaryAsPaidHandler SUCCESS: salary ID=%d marked as paid", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Marked salary #" + strconv.FormatInt(id, 10) + " as paid",
		EntityType:  "salaries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated salary
	updatedSalary, _ := h.app.Models.Salary.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "salary marked as paid successfully", updatedSalary)
}

// CancelSalaryHandler - PATCH /api/salaries/{id}/cancel
// Cancels a salary (draft or paid)
func (h *Handler) CancelSalaryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] CancelSalaryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] CancelSalaryHandler called - ID: %d", id)

	// Get existing salary
	salary, err := h.app.Models.Salary.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[SALARIES] CancelSalaryHandler ERROR: failed to fetch salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salary")
		return
	}
	if salary == nil {
		log.Printf("[SALARIES] CancelSalaryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "salary not found")
		return
	}

	// Cannot cancel already cancelled salaries
	if salary.Status == "cancelled" {
		utils.ErrorJson(w, http.StatusBadRequest, "salary is already cancelled")
		return
	}

	// Validate: Cannot cancel future months
	currentMonth := time.Now().Format("2006-01")
	if salary.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot cancel salary for future months")
		return
	}

	if err := h.app.Models.Salary.Cancel(r.Context(), id); err != nil {
		log.Printf("[SALARIES] CancelSalaryHandler ERROR: failed to cancel salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to cancel salary")
		return
	}

	log.Printf("[SALARIES] CancelSalaryHandler SUCCESS: salary ID=%d cancelled", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Cancelled salary #" + strconv.FormatInt(id, 10),
		EntityType:  "salaries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated salary
	updatedSalary, _ := h.app.Models.Salary.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "salary cancelled successfully", updatedSalary)
}

// DeleteSalaryHandler - DELETE /api/salaries/{id}
func (h *Handler) DeleteSalaryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] DeleteSalaryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] DeleteSalaryHandler called - ID: %d", id)

	// Get existing salary
	salary, err := h.app.Models.Salary.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[SALARIES] DeleteSalaryHandler ERROR: failed to fetch salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salary")
		return
	}
	if salary == nil {
		log.Printf("[SALARIES] DeleteSalaryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "salary not found")
		return
	}

	// Only draft salaries can be deleted
	if salary.Status != "draft" {
		log.Printf("[SALARIES] DeleteSalaryHandler ERROR: cannot delete non-draft salary - Status: %s", salary.Status)
		utils.ErrorJson(w, http.StatusBadRequest, "only draft salaries can be deleted")
		return
	}

	// Validate: Cannot delete future months
	currentMonth := time.Now().Format("2006-01")
	if salary.MonthYear > currentMonth {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete salary for future months")
		return
	}

	if err := h.app.Models.Salary.Delete(r.Context(), id); err != nil {
		log.Printf("[SALARIES] DeleteSalaryHandler ERROR: failed to delete salary - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete salary")
		return
	}

	log.Printf("[SALARIES] DeleteSalaryHandler SUCCESS: deleted salary ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted salary #" + strconv.FormatInt(id, 10),
		EntityType:  "salaries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "salary deleted successfully", nil)
}

// GetSalariesByEmployeeHandler - GET /api/salaries/employee/{id}
func (h *Handler) GetSalariesByEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	employeeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[SALARIES] GetSalariesByEmployeeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[SALARIES] GetSalariesByEmployeeHandler called - EmployeeID: %d", employeeID)

	// Check if employee exists
	employee, err := h.app.Models.User.GetByID(r.Context(), employeeID)
	if err != nil {
		log.Printf("[SALARIES] GetSalariesByEmployeeHandler ERROR: failed to fetch employee - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify employee")
		return
	}
	if employee == nil {
		utils.ErrorJson(w, http.StatusNotFound, "employee not found")
		return
	}

	salaries, err := h.app.Models.Salary.GetByEmployeeID(r.Context(), employeeID)
	if err != nil {
		log.Printf("[SALARIES] GetSalariesByEmployeeHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salaries")
		return
	}

	if salaries == nil {
		salaries = []models.Salary{}
	}

	log.Printf("[SALARIES] GetSalariesByEmployeeHandler SUCCESS: fetched %d salaries for employee ID=%d", len(salaries), employeeID)
	utils.SuccessJson(w, http.StatusOK, "employee salaries fetched successfully", salaries)
}

// GetSalariesByMonthHandler - GET /api/salaries/month/{month}
func (h *Handler) GetSalariesByMonthHandler(w http.ResponseWriter, r *http.Request) {
	monthYear := r.PathValue("month")

	log.Printf("[SALARIES] GetSalariesByMonthHandler called - Month: %s", monthYear)

	// Validate month format
	if len(monthYear) != 7 || monthYear[4] != '-' {
		utils.ErrorJson(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	salaries, err := h.app.Models.Salary.GetByMonth(r.Context(), monthYear)
	if err != nil {
		log.Printf("[SALARIES] GetSalariesByMonthHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch salaries")
		return
	}

	if salaries == nil {
		salaries = []models.Salary{}
	}

	log.Printf("[SALARIES] GetSalariesByMonthHandler SUCCESS: fetched %d salaries for month %s", len(salaries), monthYear)
	utils.SuccessJson(w, http.StatusOK, "monthly salaries fetched successfully", salaries)
}
