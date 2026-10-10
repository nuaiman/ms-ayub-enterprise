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
// PRE-DELIVERY ADJUSTMENT
// =============================================================================

// PreDeliveryAdjustHandler - POST /api/stores/{id}/pre-delivery-adjust
//
// Bumps the store's quantity/weight UP to a target value so a subsequent
// delivery or damage can proceed without tripping the stock guards.
//
// Only positive adjustments are supported (target must be >= current).
// To decrease a store, use the regular adjustments endpoint.
//
// After the store is updated, any CURRENT-MONTH bills for this store
// that are NOT yet attached to an invoice are refreshed with the new
// frozen snapshot. Invoiced bills are left untouched.
func (h *Handler) PreDeliveryAdjustHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORES] PreDeliveryAdjustHandler called - StoreID: %d", storeID)

	type request struct {
		TargetQuantity float64 `json:"target_quantity"`
		TargetWeight   float64 `json:"target_weight"`
		Reason         *string `json:"reason,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TargetQuantity < 0 || req.TargetWeight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "target values cannot be negative")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// Target must be >= current (this is a top-up only).
	if req.TargetQuantity < store.Quantity || req.TargetWeight < store.Weight {
		utils.ErrorJson(w, http.StatusBadRequest, "target values must be greater than or equal to the current store balance")
		return
	}

	quantityDelta := req.TargetQuantity - store.Quantity
	weightDelta := req.TargetWeight - store.Weight

	// If nothing to do, return early with current state.
	if quantityDelta == 0 && weightDelta == 0 {
		updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), storeID)
		utils.SuccessJson(w, http.StatusOK, "no adjustment needed", map[string]any{
			"adjustment":      nil,
			"store":           updatedStore,
			"bills_refreshed": map[string]int{},
		})
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	reason := "Pre-delivery adjustment"
	if req.Reason != nil && *req.Reason != "" {
		reason = *req.Reason
	}

	adjustment := &models.StoreAdjustment{
		UserID:         userID,
		StoreID:        storeID,
		AdjustmentType: "delta",
		InputWeight:    weightDelta,
		InputQuantity:  quantityDelta,
		WeightDelta:    weightDelta,
		QuantityDelta:  quantityDelta,
		Reason:         &reason,
		Notes:          nil,
		AdjustedAt:     time.Now(),
	}

	adjustmentID, err := h.app.Models.StoreAdjustment.Insert(r.Context(), adjustment)
	if err != nil {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: failed to insert adjustment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create adjustment")
		return
	}
	adjustment.ID = adjustmentID

	// Apply to store.
	if err := h.app.Models.Store.UpdateInventory(r.Context(), storeID, req.TargetQuantity, req.TargetWeight); err != nil {
		log.Printf("[STORES] PreDeliveryAdjustHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "adjustment saved but store update failed")
		return
	}

	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), storeID)

	// Refresh unbilled current-month bills.
	refreshed := h.refreshUnbilledStoreBills(r, storeID, updatedStore, lot)

	log.Printf("[STORES] PreDeliveryAdjustHandler SUCCESS: store=%d +qty=%.2f +wt=%.2f", storeID, quantityDelta, weightDelta)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "create",
		Description: "Pre-delivery adjustment #" + strconv.FormatInt(adjustmentID, 10) +
			" for store #" + strconv.FormatInt(storeID, 10) +
			" (+qty=" + strconv.FormatFloat(quantityDelta, 'f', 2, 64) +
			", +wt=" + strconv.FormatFloat(weightDelta, 'f', 2, 64) + ")",
		EntityType: "store_adjustments",
		EntityID:   adjustmentID,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "adjustment applied", map[string]any{
		"adjustment":      adjustment,
		"store":           updatedStore,
		"bills_refreshed": refreshed,
	})
}

// refreshUnbilledStoreBills refreshes current-month godown / customer_store /
// majhi bills for a store, skipping any bill that is attached to an invoice.
// Returns counts per bill type of how many were refreshed.
func (h *Handler) refreshUnbilledStoreBills(r *http.Request, storeID int64, store *models.Store, lot *models.Lot) map[string]int {
	ctx := r.Context()
	refreshed := map[string]int{}
	currentMonth := time.Now().Format("2006-01")

	// ---- Godown bills (current month, unbilled) ----
	godownBills, err := h.app.Models.GodownBill.GetByStoreID(ctx, storeID)
	if err == nil {
		for i := range godownBills {
			b := &godownBills[i]
			if b.MonthYear != currentMonth {
				continue
			}
			invoiced, _ := h.app.Models.InvoiceItem.ExistsForBill(ctx, "godown", b.ID)
			if invoiced {
				continue
			}
			b.WeightAtBilling = store.Weight
			b.QuantityAtBilling = store.Quantity
			b.WeightUnitAtBilling = &lot.WeightUnit
			b.QuantityUnitAtBilling = &lot.QuantityUnit
			b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
			if err := h.app.Models.GodownBill.Update(ctx, b); err != nil {
				log.Printf("[STORES] refreshUnbilledStoreBills ERROR: godown bill #%d - %v", b.ID, err)
				continue
			}
			refreshed["godown"]++
		}
	}

	// ---- Customer store bills (current month, unbilled) ----
	customerBills, err := h.app.Models.CustomerStoreBill.GetByStoreID(ctx, storeID)
	if err == nil {
		for i := range customerBills {
			b := &customerBills[i]
			if b.MonthYear != currentMonth {
				continue
			}
			invoiced, _ := h.app.Models.InvoiceItem.ExistsForBill(ctx, "customer_store", b.ID)
			if invoiced {
				continue
			}
			b.WeightAtBilling = store.Weight
			b.QuantityAtBilling = store.Quantity
			b.WeightUnitAtBilling = &lot.WeightUnit
			b.QuantityUnitAtBilling = &lot.QuantityUnit
			b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
			if err := h.app.Models.CustomerStoreBill.Update(ctx, b); err != nil {
				log.Printf("[STORES] refreshUnbilledStoreBills ERROR: customer store bill #%d - %v", b.ID, err)
				continue
			}
			refreshed["customer_store"]++
		}
	}

	// ---- Customer unload bills (one-time, unbilled) ----
	// Unload bills are one-time per store; only refresh if not yet invoiced.
	unloadBills, err := h.app.Models.CustomerUnloadBill.GetByStoreID(ctx, storeID)
	if err == nil {
		for i := range unloadBills {
			b := &unloadBills[i]
			invoiced, _ := h.app.Models.InvoiceItem.ExistsForBill(ctx, "customer_unload", b.ID)
			if invoiced {
				continue
			}
			b.WeightAtBilling = store.Weight
			b.QuantityAtBilling = store.Quantity
			b.WeightUnitAtBilling = &lot.WeightUnit
			b.QuantityUnitAtBilling = &lot.QuantityUnit
			b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
			if err := h.app.Models.CustomerUnloadBill.Update(ctx, b); err != nil {
				log.Printf("[STORES] refreshUnbilledStoreBills ERROR: customer unload bill #%d - %v", b.ID, err)
				continue
			}
			refreshed["customer_unload"]++
		}
	}

	// ---- Majhi bills (one-time per store, unbilled) ----
	majhiBills, err := h.app.Models.MajhiBill.GetByStoreID(ctx, storeID)
	if err == nil {
		for i := range majhiBills {
			b := &majhiBills[i]
			invoiced, _ := h.app.Models.InvoiceItem.ExistsForBill(ctx, "majhi", b.ID)
			if invoiced {
				continue
			}
			b.WeightAtBilling = store.Weight
			b.QuantityAtBilling = store.Quantity
			b.WeightUnitAtBilling = &lot.WeightUnit
			b.QuantityUnitAtBilling = &lot.QuantityUnit
			b.TotalAmount = billTotalFromSnapshot(b.BillType, b.Rate, b.WeightAtBilling, b.QuantityAtBilling)
			if err := h.app.Models.MajhiBill.Update(ctx, b); err != nil {
				log.Printf("[STORES] refreshUnbilledStoreBills ERROR: majhi bill #%d - %v", b.ID, err)
				continue
			}
			refreshed["majhi"]++
		}
	}

	return refreshed
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
//
// Bills already attached to an invoice are SKIPPED, since their amounts
// are locked in the invoice.
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

	// Reuse the same helper: refreshes current-month godown/customer_store
	// and all unbilled majhi/customer_unload bills, skipping invoiced ones.
	refreshed := h.refreshUnbilledStoreBills(r, id, store, lot)

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
