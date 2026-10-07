package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// =============================================================================
// RESPONSE SHAPES
// =============================================================================

type CustomerAdditionalBillResponse struct {
	models.CustomerAdditionalBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

// CreateCustomerAdditionalBillHandler - POST /api/customer-additional-bills
func (h *Handler) CreateCustomerAdditionalBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID  int64   `json:"customer_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_id is required")
		return
	}

	if req.Amount <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "description is required")
		return
	}

	customerExists, err := h.app.Models.Customer.Exists(r.Context(), req.CustomerID)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !customerExists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	bill := &models.CustomerAdditionalBill{
		UserID:      userID,
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		Description: description,
	}

	id, err := h.app.Models.CustomerAdditionalBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer additional bill")
		return
	}

	bill.ID = id

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] CreateCustomerAdditionalBillHandler SUCCESS: created bill ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer additional bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_additional_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer additional bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

// GetCustomerAdditionalBillHandler - GET /api/customer-additional-bills/{id}
func (h *Handler) GetCustomerAdditionalBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerAdditionalBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer additional bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer additional bill not found")
		return
	}

	responses, err := h.buildCustomerAdditionalBillResponses(r, []models.CustomerAdditionalBill{*bill})
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetCustomerAdditionalBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "customer additional bill details fetched", responses[0])
}

// GetAllCustomerAdditionalBillsHandler - GET /api/customer-additional-bills
func (h *Handler) GetAllCustomerAdditionalBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetAllCustomerAdditionalBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	customerID := r.URL.Query().Get("customer_id")

	var bills []models.CustomerAdditionalBill
	var err error

	if customerID != "" {
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		bills, err = h.app.Models.CustomerAdditionalBill.GetByCustomerID(r.Context(), id)
	} else {
		bills, err = h.app.Models.CustomerAdditionalBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetAllCustomerAdditionalBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer additional bills")
		return
	}

	responses, err := h.buildCustomerAdditionalBillResponses(r, bills)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetAllCustomerAdditionalBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] GetAllCustomerAdditionalBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "customer additional bills fetched successfully", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateCustomerAdditionalBillHandler - PATCH /api/customer-additional-bills/{id}
func (h *Handler) UpdateCustomerAdditionalBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler called - ID: %d", id)

	type request struct {
		Amount      *float64 `json:"amount,omitempty"`
		Description *string  `json:"description,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.CustomerAdditionalBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer additional bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer additional bill not found")
		return
	}

	if req.Amount != nil {
		if *req.Amount <= 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
			return
		}
		bill.Amount = *req.Amount
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "description cannot be empty")
			return
		}
		bill.Description = trimmed
	}

	if err := h.app.Models.CustomerAdditionalBill.Update(r.Context(), bill); err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer additional bill")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] UpdateCustomerAdditionalBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer additional bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_additional_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.CustomerAdditionalBill.GetByID(r.Context(), id)
	responses, _ := h.buildCustomerAdditionalBillResponses(r, []models.CustomerAdditionalBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "customer additional bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteCustomerAdditionalBillHandler - DELETE /api/customer-additional-bills/{id}
func (h *Handler) DeleteCustomerAdditionalBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerAdditionalBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer additional bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer additional bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "customer_additional", id)
	if err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.CustomerAdditionalBill.Delete(r.Context(), id); err != nil {
		log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer additional bill")
		return
	}

	log.Printf("[CUSTOMER_ADDITIONAL_BILLS] DeleteCustomerAdditionalBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer additional bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_additional_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer additional bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildCustomerAdditionalBillResponses(r *http.Request, bills []models.CustomerAdditionalBill) ([]CustomerAdditionalBillResponse, error) {
	if len(bills) == 0 {
		return []CustomerAdditionalBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "customer_additional", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]CustomerAdditionalBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.Amount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, CustomerAdditionalBillResponse{
			CustomerAdditionalBill: b,
			TotalPaid:              t.TotalPaid,
			TotalPaidThrough:       t.TotalPaidThrough,
			Remaining:              remaining,
		})
	}

	return out, nil
}
