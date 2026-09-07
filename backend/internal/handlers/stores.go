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
// STORE MANAGEMENT
// =============================================================================

// CreateStoreHandler - POST /api/stores
func (h *Handler) CreateStoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[STORES] CreateStoreHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		LotID           int64      `json:"lot_id"`
		GodownID        int64      `json:"godown_id"`
		StoreBillType   string     `json:"store_bill_type"`
		GodownCut       float64    `json:"godown_cut"`
		Quantity        float64    `json:"quantity"`
		QuantityUnit    string     `json:"quantity_unit"`
		Weight          float64    `json:"weight"`
		WeightUnit      string     `json:"weight_unit"`
		BillingStart    *time.Time `json:"billing_start,omitempty"`
		BillingEnd      *time.Time `json:"billing_end,omitempty"`
		LastPaidThrough *time.Time `json:"last_paid_through,omitempty"`
		LastPaidAmount  float64    `json:"last_paid_amount"`
		Notes           *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.LotID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "lot_id is required")
		return
	}
	if req.GodownID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "godown_id is required")
		return
	}
	if req.Quantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}
	if req.Weight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
		return
	}

	// Validate store_bill_type
	if req.StoreBillType != "" && req.StoreBillType != "weight" && req.StoreBillType != "quantity" {
		utils.ErrorJson(w, http.StatusBadRequest, "store_bill_type must be 'weight' or 'quantity'")
		return
	}

	// Check if lot exists
	exists, err := h.app.Models.Lot.Exists(r.Context(), req.LotID)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to verify lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// Check if godown exists
	exists, err = h.app.Models.Godown.Exists(r.Context(), req.GodownID)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to verify godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	// Check if store already exists for this lot and godown
	exists, err = h.app.Models.Store.ExistsByLotAndGodown(r.Context(), req.LotID, req.GodownID)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to check existing store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "store already exists for this lot and godown")
		return
	}

	// Set default values
	billingStart := req.BillingStart
	if billingStart == nil {
		now := time.Now()
		billingStart = &now
	}

	quantityUnit := req.QuantityUnit
	if quantityUnit == "" {
		quantityUnit = "units"
	}

	weightUnit := req.WeightUnit
	if weightUnit == "" {
		weightUnit = "kg"
	}

	storeBillType := req.StoreBillType
	if storeBillType == "" {
		storeBillType = "quantity"
	}

	store := &models.Store{
		LotID:           req.LotID,
		GodownID:        req.GodownID,
		StoreBillType:   storeBillType,
		GodownCut:       req.GodownCut,
		Quantity:        req.Quantity,
		QuantityUnit:    quantityUnit,
		Weight:          req.Weight,
		WeightUnit:      weightUnit,
		IsActive:        true,
		BillingStart:    *billingStart,
		BillingEnd:      req.BillingEnd,
		LastPaidThrough: req.LastPaidThrough,
		LastPaidAmount:  req.LastPaidAmount,
		Notes:           req.Notes,
	}

	id, err := h.app.Models.Store.Insert(r.Context(), store)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to insert store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create store")
		return
	}

	store.ID = id

	log.Printf("[STORES] CreateStoreHandler SUCCESS: created store ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created store #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "store created successfully", store)
}

// GetStoreHandler - GET /api/stores/{id}
func (h *Handler) GetStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] GetStoreHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] GetStoreHandler called - ID: %d", id)

	store, err := h.app.Models.Store.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] GetStoreHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}

	if store == nil {
		log.Printf("[STORES] GetStoreHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	log.Printf("[STORES] GetStoreHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "store details fetched", store)
}

// GetAllStoresHandler - GET /api/stores
func (h *Handler) GetAllStoresHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[STORES] GetAllStoresHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	lotID := r.URL.Query().Get("lot_id")
	godownID := r.URL.Query().Get("godown_id")
	activeOnly := r.URL.Query().Get("active") == "true"
	withInventory := r.URL.Query().Get("with_inventory") == "true"

	var stores []models.Store
	var err error

	switch {
	case withInventory:
		stores, err = h.app.Models.Store.GetWithInventory(r.Context())
	case lotID != "":
		id, _ := strconv.ParseInt(lotID, 10, 64)
		stores, err = h.app.Models.Store.GetByLotID(r.Context(), id)
	case godownID != "":
		id, _ := strconv.ParseInt(godownID, 10, 64)
		stores, err = h.app.Models.Store.GetByGodownID(r.Context(), id)
	case activeOnly:
		stores, err = h.app.Models.Store.GetActive(r.Context())
	default:
		stores, err = h.app.Models.Store.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[STORES] GetAllStoresHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch stores")
		return
	}

	if stores == nil {
		stores = []models.Store{}
	}

	log.Printf("[STORES] GetAllStoresHandler SUCCESS: fetched %d stores", len(stores))
	utils.SuccessJson(w, http.StatusOK, "stores fetched successfully", stores)
}

// UpdateStoreHandler - PATCH /api/stores/{id}
func (h *Handler) UpdateStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] UpdateStoreHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] UpdateStoreHandler called - ID: %d", id)

	type request struct {
		StoreBillType   *string    `json:"store_bill_type,omitempty"`
		GodownCut       *float64   `json:"godown_cut,omitempty"`
		Quantity        *float64   `json:"quantity,omitempty"`
		QuantityUnit    *string    `json:"quantity_unit,omitempty"`
		Weight          *float64   `json:"weight,omitempty"`
		WeightUnit      *string    `json:"weight_unit,omitempty"`
		BillingStart    *time.Time `json:"billing_start,omitempty"`
		BillingEnd      *time.Time `json:"billing_end,omitempty"`
		LastPaidThrough *time.Time `json:"last_paid_through,omitempty"`
		LastPaidAmount  *float64   `json:"last_paid_amount,omitempty"`
		Notes           *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORES] UpdateStoreHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing store
	store, err := h.app.Models.Store.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] UpdateStoreHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		log.Printf("[STORES] UpdateStoreHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Apply updates
	if req.StoreBillType != nil {
		if *req.StoreBillType != "weight" && *req.StoreBillType != "quantity" {
			utils.ErrorJson(w, http.StatusBadRequest, "store_bill_type must be 'weight' or 'quantity'")
			return
		}
		store.StoreBillType = *req.StoreBillType
	}
	if req.GodownCut != nil {
		if *req.GodownCut < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "godown_cut cannot be negative")
			return
		}
		store.GodownCut = *req.GodownCut
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
			return
		}
		store.Quantity = *req.Quantity
	}
	if req.QuantityUnit != nil {
		if *req.QuantityUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity_unit cannot be empty")
			return
		}
		store.QuantityUnit = *req.QuantityUnit
	}
	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		store.Weight = *req.Weight
	}
	if req.WeightUnit != nil {
		if *req.WeightUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "weight_unit cannot be empty")
			return
		}
		store.WeightUnit = *req.WeightUnit
	}
	if req.BillingStart != nil {
		store.BillingStart = *req.BillingStart
	}
	if req.BillingEnd != nil {
		store.BillingEnd = req.BillingEnd
	}
	if req.LastPaidThrough != nil {
		store.LastPaidThrough = req.LastPaidThrough
	}
	if req.LastPaidAmount != nil {
		if *req.LastPaidAmount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "last_paid_amount cannot be negative")
			return
		}
		store.LastPaidAmount = *req.LastPaidAmount
	}
	if req.Notes != nil {
		store.Notes = req.Notes
	}

	if err := h.app.Models.Store.Update(r.Context(), store); err != nil {
		log.Printf("[STORES] UpdateStoreHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update store")
		return
	}

	log.Printf("[STORES] UpdateStoreHandler SUCCESS: updated store ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated store #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated store
	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "store updated successfully", updatedStore)
}

// ToggleStoreActiveHandler - PATCH /api/stores/{id}/toggle-active
func (h *Handler) ToggleStoreActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] ToggleStoreActiveHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] ToggleStoreActiveHandler called - ID: %d", id)

	store, err := h.app.Models.Store.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] ToggleStoreActiveHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		log.Printf("[STORES] ToggleStoreActiveHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	newStatus := !store.IsActive

	if err := h.app.Models.Store.ToggleActive(r.Context(), id, newStatus); err != nil {
		log.Printf("[STORES] ToggleStoreActiveHandler ERROR: failed to toggle status - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update store status")
		return
	}

	log.Printf("[STORES] ToggleStoreActiveHandler SUCCESS: store ID=%d now active=%v", id, newStatus)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	action := "deactivated"
	if newStatus {
		action = "activated"
	}
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Store " + action + ": #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	store.IsActive = newStatus
	utils.SuccessJson(w, http.StatusOK, "store status updated successfully", store)
}

// DeleteStoreHandler - DELETE /api/stores/{id}
func (h *Handler) DeleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] DeleteStoreHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] DeleteStoreHandler called - ID: %d", id)

	store, err := h.app.Models.Store.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] DeleteStoreHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		log.Printf("[STORES] DeleteStoreHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	if err := h.app.Models.Store.Delete(r.Context(), id); err != nil {
		log.Printf("[STORES] DeleteStoreHandler ERROR: failed to delete store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete store")
		return
	}

	log.Printf("[STORES] DeleteStoreHandler SUCCESS: deleted store ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted store #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "store deleted successfully", nil)
}
