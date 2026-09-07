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
// DAMAGE MANAGEMENT
// =============================================================================

// CreateDamageHandler - POST /api/damages
func (h *Handler) CreateDamageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DAMAGES] CreateDamageHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		StoreID      int64      `json:"store_id"`
		Quantity     float64    `json:"quantity"`
		QuantityUnit string     `json:"quantity_unit"`
		Weight       float64    `json:"weight"`
		WeightUnit   string     `json:"weight_unit"`
		DamageDate   *time.Time `json:"damage_date,omitempty"`
		Reason       string     `json:"reason"`
		Amount       float64    `json:"amount"`
		Notes        *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
		return
	}
	if req.Reason == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "reason is required")
		return
	}
	if req.Quantity < 0 || req.Weight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "quantity and weight cannot be negative")
		return
	}
	if req.Quantity == 0 && req.Weight == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "either quantity or weight must be greater than 0")
		return
	}
	if req.Amount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
		return
	}

	// Check if store exists
	exists, err := h.app.Models.Store.Exists(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to verify store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Get store to check if we have enough inventory
	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Check if we have enough quantity/weight
	if req.Quantity > 0 && store.Quantity < req.Quantity {
		utils.ErrorJson(w, http.StatusBadRequest, "insufficient quantity in store")
		return
	}
	if req.Weight > 0 && store.Weight < req.Weight {
		utils.ErrorJson(w, http.StatusBadRequest, "insufficient weight in store")
		return
	}

	// Get current user ID
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default values
	damageDate := req.DamageDate
	if damageDate == nil {
		now := time.Now()
		damageDate = &now
	}

	quantityUnit := req.QuantityUnit
	if quantityUnit == "" {
		quantityUnit = "units"
	}

	weightUnit := req.WeightUnit
	if weightUnit == "" {
		weightUnit = "kg"
	}

	damage := &models.Damage{
		StoreID:      req.StoreID,
		UserID:       userID,
		Quantity:     req.Quantity,
		QuantityUnit: quantityUnit,
		Weight:       req.Weight,
		WeightUnit:   weightUnit,
		DamageDate:   *damageDate,
		Reason:       req.Reason,
		Amount:       req.Amount,
		Notes:        req.Notes,
	}

	id, err := h.app.Models.Damage.Insert(r.Context(), damage)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to insert damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create damage record")
		return
	}

	damage.ID = id

	// Update store inventory
	newQuantity := store.Quantity - req.Quantity
	newWeight := store.Weight - req.Weight
	if err := h.app.Models.Store.UpdateInventory(r.Context(), req.StoreID, newQuantity, newWeight); err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to update store inventory - %v", err)
		// Note: We should probably rollback the damage insertion here
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update store inventory")
		return
	}

	log.Printf("[DAMAGES] CreateDamageHandler SUCCESS: created damage ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created damage record #" + strconv.FormatInt(id, 10) + " for store #" + strconv.FormatInt(req.StoreID, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "damage record created successfully", damage)
}

// GetDamageHandler - GET /api/damages/{id}
func (h *Handler) GetDamageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] GetDamageHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DAMAGES] GetDamageHandler called - ID: %d", id)

	damage, err := h.app.Models.Damage.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DAMAGES] GetDamageHandler ERROR: failed to fetch damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage record")
		return
	}

	if damage == nil {
		log.Printf("[DAMAGES] GetDamageHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "damage record not found")
		return
	}

	log.Printf("[DAMAGES] GetDamageHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "damage details fetched", damage)
}

// GetAllDamagesHandler - GET /api/damages
func (h *Handler) GetAllDamagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DAMAGES] GetAllDamagesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	storeID := r.URL.Query().Get("store_id")
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var damages []models.Damage
	var err error

	switch {
	case storeID != "":
		id, _ := strconv.ParseInt(storeID, 10, 64)
		damages, err = h.app.Models.Damage.GetByStoreID(r.Context(), id)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		damages, err = h.app.Models.Damage.GetByUserID(r.Context(), id)
	case startDate != "" && endDate != "":
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		damages, err = h.app.Models.Damage.GetByDateRange(r.Context(), start, end)
	default:
		damages, err = h.app.Models.Damage.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[DAMAGES] GetAllDamagesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damages")
		return
	}

	if damages == nil {
		damages = []models.Damage{}
	}

	log.Printf("[DAMAGES] GetAllDamagesHandler SUCCESS: fetched %d damages", len(damages))
	utils.SuccessJson(w, http.StatusOK, "damages fetched successfully", damages)
}

// UpdateDamageHandler - PATCH /api/damages/{id}
func (h *Handler) UpdateDamageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DAMAGES] UpdateDamageHandler called - ID: %d", id)

	type request struct {
		Quantity     *float64   `json:"quantity,omitempty"`
		QuantityUnit *string    `json:"quantity_unit,omitempty"`
		Weight       *float64   `json:"weight,omitempty"`
		WeightUnit   *string    `json:"weight_unit,omitempty"`
		DamageDate   *time.Time `json:"damage_date,omitempty"`
		Reason       *string    `json:"reason,omitempty"`
		Amount       *float64   `json:"amount,omitempty"`
		Notes        *string    `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing damage
	damage, err := h.app.Models.Damage.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: failed to fetch damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage")
		return
	}
	if damage == nil {
		log.Printf("[DAMAGES] UpdateDamageHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "damage not found")
		return
	}

	// Check if store exists
	store, err := h.app.Models.Store.GetByID(r.Context(), damage.StoreID)
	if err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}

	// Track old values for inventory adjustment
	oldQuantity := damage.Quantity
	oldWeight := damage.Weight

	// Apply updates
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
			return
		}
		damage.Quantity = *req.Quantity
	}
	if req.QuantityUnit != nil {
		if *req.QuantityUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity_unit cannot be empty")
			return
		}
		damage.QuantityUnit = *req.QuantityUnit
	}
	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		damage.Weight = *req.Weight
	}
	if req.WeightUnit != nil {
		if *req.WeightUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "weight_unit cannot be empty")
			return
		}
		damage.WeightUnit = *req.WeightUnit
	}
	if req.DamageDate != nil {
		damage.DamageDate = *req.DamageDate
	}
	if req.Reason != nil {
		if *req.Reason == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "reason cannot be empty")
			return
		}
		damage.Reason = *req.Reason
	}
	if req.Amount != nil {
		if *req.Amount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
			return
		}
		damage.Amount = *req.Amount
	}
	if req.Notes != nil {
		damage.Notes = req.Notes
	}

	if err := h.app.Models.Damage.Update(r.Context(), damage); err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: failed to update damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update damage")
		return
	}

	// Update store inventory - restore old values and deduct new ones
	if store != nil {
		newQuantity := store.Quantity + oldQuantity - damage.Quantity
		newWeight := store.Weight + oldWeight - damage.Weight
		if err := h.app.Models.Store.UpdateInventory(r.Context(), store.ID, newQuantity, newWeight); err != nil {
			log.Printf("[DAMAGES] UpdateDamageHandler ERROR: failed to update store inventory - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update store inventory")
			return
		}
	}

	log.Printf("[DAMAGES] UpdateDamageHandler SUCCESS: updated damage ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated damage record #" + strconv.FormatInt(id, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated damage
	updatedDamage, _ := h.app.Models.Damage.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "damage updated successfully", updatedDamage)
}

// DeleteDamageHandler - DELETE /api/damages/{id}
func (h *Handler) DeleteDamageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DAMAGES] DeleteDamageHandler called - ID: %d", id)

	// Get existing damage
	damage, err := h.app.Models.Damage.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: failed to fetch damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage")
		return
	}
	if damage == nil {
		log.Printf("[DAMAGES] DeleteDamageHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "damage not found")
		return
	}

	// Restore inventory to store
	store, err := h.app.Models.Store.GetByID(r.Context(), damage.StoreID)
	if err != nil {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store != nil {
		newQuantity := store.Quantity + damage.Quantity
		newWeight := store.Weight + damage.Weight
		if err := h.app.Models.Store.UpdateInventory(r.Context(), store.ID, newQuantity, newWeight); err != nil {
			log.Printf("[DAMAGES] DeleteDamageHandler ERROR: failed to restore store inventory - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to restore store inventory")
			return
		}
	}

	if err := h.app.Models.Damage.Delete(r.Context(), id); err != nil {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: failed to delete damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete damage")
		return
	}

	log.Printf("[DAMAGES] DeleteDamageHandler SUCCESS: deleted damage ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted damage record #" + strconv.FormatInt(id, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "damage deleted successfully", nil)
}
