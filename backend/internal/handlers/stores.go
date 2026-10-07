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
		LotID     int64      `json:"lot_id"`
		GodownID  int64      `json:"godown_id"`
		Weight    float64    `json:"weight"`
		Quantity  float64    `json:"quantity"`
		StartDate *time.Time `json:"start_date,omitempty"`
		IsActive  *bool      `json:"is_active,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.LotID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "lot_id is required")
		return
	}
	if req.GodownID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "godown_id is required")
		return
	}
	if req.Weight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
		return
	}
	if req.Quantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}

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

	godownExists, err := h.app.Models.Godown.Exists(r.Context(), req.GodownID)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to verify godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
		return
	}
	if !godownExists {
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[STORES] CreateStoreHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	startDate := time.Now()
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	store := &models.Store{
		UserID:    userID,
		LotID:     req.LotID,
		GodownID:  req.GodownID,
		Weight:    req.Weight,
		Quantity:  req.Quantity,
		StartDate: startDate,
		IsActive:  isActive,
	}

	id, err := h.app.Models.Store.Insert(r.Context(), store)
	if err != nil {
		log.Printf("[STORES] CreateStoreHandler ERROR: failed to insert store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create store")
		return
	}

	store.ID = id

	log.Printf("[STORES] CreateStoreHandler SUCCESS: created store ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created store #" + strconv.FormatInt(id, 10) + " for lot #" + strconv.FormatInt(req.LotID, 10),
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
	activeOnly := r.URL.Query().Get("active") == "true"
	withInventory := r.URL.Query().Get("with_inventory") == "true"

	var stores []models.Store
	var err error

	switch {
	case withInventory:
		stores, err = h.app.Models.Store.GetWithInventory(r.Context())
	case lotID != "":
		id, parseErr := strconv.ParseInt(lotID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid lot_id")
			return
		}
		stores, err = h.app.Models.Store.GetByLotID(r.Context(), id)
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
		Weight    *float64   `json:"weight,omitempty"`
		Quantity  *float64   `json:"quantity,omitempty"`
		StartDate *time.Time `json:"start_date,omitempty"`
		IsActive  *bool      `json:"is_active,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORES] UpdateStoreHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

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

	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		store.Weight = *req.Weight
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
			return
		}
		store.Quantity = *req.Quantity
	}
	if req.StartDate != nil {
		store.StartDate = *req.StartDate
	}
	if req.IsActive != nil {
		store.IsActive = *req.IsActive
	}

	if err := h.app.Models.Store.Update(r.Context(), store); err != nil {
		log.Printf("[STORES] UpdateStoreHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update store")
		return
	}

	log.Printf("[STORES] UpdateStoreHandler SUCCESS: updated store ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated store #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

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

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	action := "deactivated"
	if newStatus {
		action = "activated"
	}
	h.logAudit(r.Context(), &models.Log{
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

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
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

// =============================================================================
// REFRESH CURRENT-MONTH BILLS
// =============================================================================

// RefreshStoreBillsHandler - POST /api/stores/{id}/refresh-bills
//
// Refreshes the current month's frozen values (weight_at_billing,
// quantity_at_billing, weight_unit_at_billing, quantity_unit_at_billing,
// total_amount) for whatever godown / majhi / customer_store bills already
// exist for this store in the current month. Only the current month is
// touched; historical bills remain frozen.
//
// Rate is NOT touched. Only the frozen quantity snapshot is refreshed, so
// the total recomputes from the same rate × fresh weight/quantity.
func (h *Handler) RefreshStoreBillsHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] RefreshStoreBillsHandler called - StoreID: %d", id)

	store, err := h.app.Models.Store.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		log.Printf("[STORES] RefreshStoreBillsHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	currentMonth := time.Now().Format("2006-01")

	refreshed := map[string]int{}

	// ---- Godown bills (current month) ----
	godownBills, err := h.app.Models.GodownBill.GetByStoreID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to fetch godown bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown bills")
		return
	}
	for i := range godownBills {
		b := &godownBills[i]
		if b.MonthYear != currentMonth {
			continue
		}
		b.WeightAtBilling = store.Weight
		b.QuantityAtBilling = store.Quantity
		b.WeightUnitAtBilling = &lot.WeightUnit
		b.QuantityUnitAtBilling = &lot.QuantityUnit
		b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
		if err := h.app.Models.GodownBill.Update(r.Context(), b); err != nil {
			log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to update godown bill #%d - %v", b.ID, err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to refresh godown bill")
			return
		}
		refreshed["godown"]++
	}

	// ---- Customer store bills (current month) ----
	customerBills, err := h.app.Models.CustomerStoreBill.GetByStoreID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to fetch customer store bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer store bills")
		return
	}
	for i := range customerBills {
		b := &customerBills[i]
		if b.MonthYear != currentMonth {
			continue
		}
		b.WeightAtBilling = store.Weight
		b.QuantityAtBilling = store.Quantity
		b.WeightUnitAtBilling = &lot.WeightUnit
		b.QuantityUnitAtBilling = &lot.QuantityUnit
		b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
		if err := h.app.Models.CustomerStoreBill.Update(r.Context(), b); err != nil {
			log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to update customer store bill #%d - %v", b.ID, err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to refresh customer store bill")
			return
		}
		refreshed["customer_store"]++
	}

	// ---- Majhi bills (one-time, no month — refresh all for this store) ----
	majhiBills, err := h.app.Models.MajhiBill.GetByStoreID(r.Context(), id)
	if err != nil {
		log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to fetch majhi bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi bills")
		return
	}
	for i := range majhiBills {
		b := &majhiBills[i]
		b.WeightAtBilling = store.Weight
		b.QuantityAtBilling = store.Quantity
		b.WeightUnitAtBilling = &lot.WeightUnit
		b.QuantityUnitAtBilling = &lot.QuantityUnit
		b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
		if err := h.app.Models.MajhiBill.Update(r.Context(), b); err != nil {
			log.Printf("[STORES] RefreshStoreBillsHandler ERROR: failed to update majhi bill #%d - %v", b.ID, err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to refresh majhi bill")
			return
		}
		refreshed["majhi"]++
	}

	log.Printf("[STORES] RefreshStoreBillsHandler SUCCESS: store=%d refreshed=%v", id, refreshed)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Refreshed current-month bills for store #" + strconv.FormatInt(id, 10),
		EntityType:  "stores",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "bills refreshed successfully", map[string]any{
		"refreshed": refreshed,
	})
}
