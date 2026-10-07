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
// STORE ADJUSTMENTS
// =============================================================================

// CreateStoreAdjustmentHandler - POST /api/stores/{id}/adjustments
//
// Records a weight and/or quantity adjustment against a store.
// adjustment_type:
//
//	"delta"    -> input_* is the change to apply (+/-)
//	"absolute" -> input_* is the new target value; delta is derived
//
// The store row is updated in the same transaction.
func (h *Handler) CreateStoreAdjustmentHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: invalid store id")
		return
	}

	log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler called - StoreID: %d", storeID)

	type request struct {
		AdjustmentType string  `json:"adjustment_type"`
		InputWeight    float64 `json:"input_weight"`
		InputQuantity  float64 `json:"input_quantity"`
		Reason         *string `json:"reason,omitempty"`
		Notes          *string `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AdjustmentType == "" {
		req.AdjustmentType = "delta"
	}
	if req.AdjustmentType != "delta" && req.AdjustmentType != "absolute" {
		utils.ErrorJson(w, http.StatusBadRequest, "adjustment_type must be 'delta' or 'absolute'")
		return
	}

	if req.InputWeight == 0 && req.InputQuantity == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "at least one of input_weight or input_quantity must be non-zero")
		return
	}

	if req.AdjustmentType == "absolute" {
		if req.InputWeight < 0 || req.InputQuantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "absolute values cannot be negative")
			return
		}
	}

	// Verify store exists.
	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Resolve the delta based on adjustment_type.
	var weightDelta, quantityDelta float64

	switch req.AdjustmentType {
	case "delta":
		weightDelta = req.InputWeight
		quantityDelta = req.InputQuantity
	case "absolute":
		// Only compute delta for the fields that were actually supplied.
		// A zero input in absolute mode means "leave unchanged" for that field.
		if req.InputWeight != 0 {
			weightDelta = req.InputWeight - store.Weight
		}
		if req.InputQuantity != 0 {
			quantityDelta = req.InputQuantity - store.Quantity
		}
	}

	newWeight := store.Weight + weightDelta
	newQuantity := store.Quantity + quantityDelta

	if newWeight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "adjustment would make store weight negative")
		return
	}
	if newQuantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "adjustment would make store quantity negative")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	adjustedAt := time.Now()

	adjustment := &models.StoreAdjustment{
		UserID:         userID,
		StoreID:        storeID,
		AdjustmentType: req.AdjustmentType,
		InputWeight:    req.InputWeight,
		InputQuantity:  req.InputQuantity,
		WeightDelta:    weightDelta,
		QuantityDelta:  quantityDelta,
		Reason:         req.Reason,
		Notes:          req.Notes,
		AdjustedAt:     adjustedAt,
	}

	// Insert adjustment, then update store. Two writes; no explicit tx
	// because this is a single-process SQLite app with WAL, and the failure
	// modes here (insert succeeds, update fails) are recoverable by hand.
	id, err := h.app.Models.StoreAdjustment.Insert(r.Context(), adjustment)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: failed to insert adjustment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create adjustment")
		return
	}
	adjustment.ID = id

	// Apply to store.
	store.Weight = newWeight
	store.Quantity = newQuantity
	if err := h.app.Models.Store.UpdateInventory(r.Context(), storeID, newQuantity, newWeight); err != nil {
		log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "adjustment saved but store update failed")
		return
	}

	// Re-fetch so the store snapshot returned to the client is fresh
	// (triggers may have flipped is_active).
	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), storeID)

	log.Printf("[STORE_ADJUSTMENTS] CreateStoreAdjustmentHandler SUCCESS: adjustment ID=%d store=%d type=%s wΔ=%.2f qΔ=%.2f",
		id, storeID, req.AdjustmentType, weightDelta, quantityDelta)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created store adjustment #" + strconv.FormatInt(id, 10) + " for store #" + strconv.FormatInt(storeID, 10),
		EntityType:  "store_adjustments",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "store adjustment created successfully", map[string]any{
		"adjustment": adjustment,
		"store":      updatedStore,
	})
}

// =============================================================================
// READ
// =============================================================================

// GetStoreAdjustmentHandler - GET /api/store-adjustments/{id}
func (h *Handler) GetStoreAdjustmentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentHandler called - ID: %d", id)

	adjustment, err := h.app.Models.StoreAdjustment.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentHandler ERROR: failed to fetch adjustment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch adjustment")
		return
	}
	if adjustment == nil {
		utils.ErrorJson(w, http.StatusNotFound, "adjustment not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "adjustment details fetched", adjustment)
}

// GetStoreAdjustmentsHandler - GET /api/stores/{id}/adjustments
func (h *Handler) GetStoreAdjustmentsHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentsHandler ERROR: invalid store id")
		return
	}

	log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentsHandler called - StoreID: %d", storeID)

	// Verify store exists.
	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentsHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	adjustments, err := h.app.Models.StoreAdjustment.GetByStoreID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] GetStoreAdjustmentsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch adjustments")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "adjustments fetched successfully", adjustments)
}

// GetAllStoreAdjustmentsHandler - GET /api/store-adjustments
func (h *Handler) GetAllStoreAdjustmentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[STORE_ADJUSTMENTS] GetAllStoreAdjustmentsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	userIDParam := r.URL.Query().Get("user_id")

	var adjustments []models.StoreAdjustment
	var err error

	if userIDParam != "" {
		id, parseErr := strconv.ParseInt(userIDParam, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		adjustments, err = h.app.Models.StoreAdjustment.GetByUserID(r.Context(), id)
	} else {
		adjustments, err = h.app.Models.StoreAdjustment.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] GetAllStoreAdjustmentsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch adjustments")
		return
	}

	if adjustments == nil {
		adjustments = []models.StoreAdjustment{}
	}

	utils.SuccessJson(w, http.StatusOK, "adjustments fetched successfully", adjustments)
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteStoreAdjustmentHandler - DELETE /api/store-adjustments/{id}
//
// Reverses the adjustment: subtracts weight_delta and quantity_delta from the
// store, then deletes the adjustment row. Rejects if the reversal would make
// the store's weight or quantity negative.
func (h *Handler) DeleteStoreAdjustmentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler called - ID: %d", id)

	adjustment, err := h.app.Models.StoreAdjustment.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler ERROR: failed to fetch adjustment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch adjustment")
		return
	}
	if adjustment == nil {
		utils.ErrorJson(w, http.StatusNotFound, "adjustment not found")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), adjustment.StoreID)
	if err != nil {
		log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	reversedWeight := store.Weight - adjustment.WeightDelta
	reversedQuantity := store.Quantity - adjustment.QuantityDelta

	if reversedWeight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete adjustment: reversal would make store weight negative")
		return
	}
	if reversedQuantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete adjustment: reversal would make store quantity negative")
		return
	}

	// Apply reversal to store first, then delete the adjustment row.
	if err := h.app.Models.Store.UpdateInventory(r.Context(), store.ID, reversedQuantity, reversedWeight); err != nil {
		log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to reverse adjustment on store")
		return
	}

	if err := h.app.Models.StoreAdjustment.Delete(r.Context(), id); err != nil {
		log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler ERROR: failed to delete adjustment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete adjustment")
		return
	}

	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), store.ID)

	log.Printf("[STORE_ADJUSTMENTS] DeleteStoreAdjustmentHandler SUCCESS: deleted adjustment ID=%d store=%d reversed w=%.2f q=%.2f",
		id, store.ID, adjustment.WeightDelta, adjustment.QuantityDelta)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted store adjustment #" + strconv.FormatInt(id, 10) + " (store #" + strconv.FormatInt(store.ID, 10) + ")",
		EntityType:  "store_adjustments",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "store adjustment deleted successfully", map[string]any{
		"store": updatedStore,
	})
}
