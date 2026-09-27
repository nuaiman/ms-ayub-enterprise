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
// CUSTOMER ADDITIONAL CHARGE MANAGEMENT
// =============================================================================

var validAdditionalChargeEntityTypes = map[string]bool{
	"lot":       true,
	"store":     true,
	"delivery":  true,
	"transport": true,
	"damage":    true,
	"godown":    true,
}

// CreateAdditionalChargeHandler - POST /api/customer-additional-charges
func (h *Handler) CreateAdditionalChargeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID  int64   `json:"customer_id"`
		EntityType  string  `json:"entity_type"`
		EntityID    int64   `json:"entity_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_id is required")
		return
	}
	if req.EntityType == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_type is required")
		return
	}
	if !validAdditionalChargeEntityTypes[req.EntityType] {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_type must be one of: lot, store, delivery, transport, damage, godown")
		return
	}
	if req.EntityID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_id is required")
		return
	}
	if req.Amount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
		return
	}
	if req.Description == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "description is required")
		return
	}

	exists, err := h.app.Models.Customer.Exists(r.Context(), req.CustomerID)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	charge := &models.AdditionalCharge{
		UserID:            userID,
		CustomerID:        req.CustomerID,
		EntityType:        req.EntityType,
		EntityID:          req.EntityID,
		Amount:            req.Amount,
		Description:       req.Description,
		CustomerTotalPaid: 0,
	}

	id, err := h.app.Models.AdditionalCharge.Insert(r.Context(), charge)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler ERROR: failed to insert charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create additional charge")
		return
	}

	charge.ID = id

	log.Printf("[ADDITIONAL_CHARGES] CreateAdditionalChargeHandler SUCCESS: created charge ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created additional charge #" + strconv.FormatInt(id, 10) + " for customer #" + strconv.FormatInt(req.CustomerID, 10),
		EntityType:  "customer_additional_charges",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "additional charge created successfully", charge)
}

// GetAdditionalChargeHandler - GET /api/customer-additional-charges/{id}
func (h *Handler) GetAdditionalChargeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ADDITIONAL_CHARGES] GetAdditionalChargeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] GetAdditionalChargeHandler called - ID: %d", id)

	charge, err := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] GetAdditionalChargeHandler ERROR: failed to fetch charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch additional charge")
		return
	}
	if charge == nil {
		log.Printf("[ADDITIONAL_CHARGES] GetAdditionalChargeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "additional charge not found")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] GetAdditionalChargeHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "additional charge details fetched", charge)
}

// GetAllAdditionalChargesHandler - GET /api/customer-additional-charges
func (h *Handler) GetAllAdditionalChargesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ADDITIONAL_CHARGES] GetAllAdditionalChargesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	customerID := r.URL.Query().Get("customer_id")
	entityType := r.URL.Query().Get("entity_type")
	entityID := r.URL.Query().Get("entity_id")
	search := r.URL.Query().Get("search")

	var charges []models.AdditionalCharge
	var err error

	switch {
	case entityType != "" && entityID != "":
		id, parseErr := strconv.ParseInt(entityID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid entity_id")
			return
		}
		if !validAdditionalChargeEntityTypes[entityType] {
			utils.ErrorJson(w, http.StatusBadRequest, "entity_type must be one of: lot, store, delivery, transport, damage, godown")
			return
		}
		charges, err = h.app.Models.AdditionalCharge.GetByEntity(r.Context(), entityType, id)
	case customerID != "":
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		charges, err = h.app.Models.AdditionalCharge.GetByCustomerID(r.Context(), id)
	default:
		charges, err = h.app.Models.AdditionalCharge.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] GetAllAdditionalChargesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch additional charges")
		return
	}

	if charges == nil {
		charges = []models.AdditionalCharge{}
	}

	if search != "" {
		searchLower := strings.ToLower(search)
		filtered := make([]models.AdditionalCharge, 0, len(charges))
		for _, c := range charges {
			if strings.Contains(strings.ToLower(c.Description), searchLower) {
				filtered = append(filtered, c)
			}
		}
		charges = filtered
	}

	log.Printf("[ADDITIONAL_CHARGES] GetAllAdditionalChargesHandler SUCCESS: fetched %d charges", len(charges))
	utils.SuccessJson(w, http.StatusOK, "additional charges fetched successfully", charges)
}

// UpdateAdditionalChargeHandler - PATCH /api/customer-additional-charges/{id}
func (h *Handler) UpdateAdditionalChargeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler called - ID: %d", id)

	type request struct {
		Amount      *float64 `json:"amount,omitempty"`
		Description *string  `json:"description,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	charge, err := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler ERROR: failed to fetch charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch additional charge")
		return
	}
	if charge == nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "additional charge not found")
		return
	}

	if req.Amount != nil {
		if *req.Amount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
			return
		}
		charge.Amount = *req.Amount
	}
	if req.Description != nil {
		if *req.Description == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "description cannot be empty")
			return
		}
		charge.Description = *req.Description
	}

	if err := h.app.Models.AdditionalCharge.Update(r.Context(), charge); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler ERROR: failed to update charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update additional charge")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeHandler SUCCESS: updated charge ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated additional charge #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_additional_charges",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "additional charge updated successfully", updated)
}

// UpdateAdditionalChargeCustomerPaymentHandler - PATCH /api/customer-additional-charges/{id}/customer-payment
func (h *Handler) UpdateAdditionalChargeCustomerPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler called - ID: %d", id)

	type request struct {
		CustomerTotalPaid        float64    `json:"customer_total_paid"`
		CustomerTotalPaidThrough *time.Time `json:"customer_total_paid_through,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerTotalPaid < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_total_paid cannot be negative")
		return
	}

	charge, err := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler ERROR: failed to fetch charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch additional charge")
		return
	}
	if charge == nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "additional charge not found")
		return
	}

	if req.CustomerTotalPaid > charge.Amount {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_total_paid cannot exceed the charge amount")
		return
	}

	if err := h.app.Models.AdditionalCharge.UpdateCustomerPayment(r.Context(), id, req.CustomerTotalPaid, req.CustomerTotalPaidThrough); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler ERROR: failed to update payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer payment")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] UpdateAdditionalChargeCustomerPaymentHandler SUCCESS: updated payment for charge ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer payment for additional charge #" + strconv.FormatInt(id, 10) + " to " + strconv.FormatFloat(req.CustomerTotalPaid, 'f', 2, 64),
		EntityType:  "customer_additional_charges",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer payment updated successfully", updated)
}

// DeleteAdditionalChargeHandler - DELETE /api/customer-additional-charges/{id}
func (h *Handler) DeleteAdditionalChargeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler called - ID: %d", id)

	charge, err := h.app.Models.AdditionalCharge.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler ERROR: failed to fetch charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch additional charge")
		return
	}
	if charge == nil {
		log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "additional charge not found")
		return
	}

	if err := h.app.Models.AdditionalCharge.Delete(r.Context(), id); err != nil {
		log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler ERROR: failed to delete charge - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete additional charge")
		return
	}

	log.Printf("[ADDITIONAL_CHARGES] DeleteAdditionalChargeHandler SUCCESS: deleted charge ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted additional charge #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_additional_charges",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "additional charge deleted successfully", nil)
}

// containsFold is a small case-insensitive substring check used by the list
// handler's search filter.
func containsFold(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	hl := len(haystack)
	nl := len(needle)
	if nl > hl {
		return false
	}
	for i := 0; i+nl <= hl; i++ {
		match := true
		for j := 0; j < nl; j++ {
			if toLower(haystack[i+j]) != toLower(needle[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
