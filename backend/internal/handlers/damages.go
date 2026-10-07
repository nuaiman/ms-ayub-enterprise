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
// DAMAGES
// =============================================================================

// CreateDamageHandler - POST /api/damages
//
// store_id is required.
//
// quantity_unit and weight_unit are ignored if sent by the client - they are
// always derived from the parent lot of the store.
func (h *Handler) CreateDamageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DAMAGES] CreateDamageHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		StoreID    int64      `json:"store_id"`
		Quantity   float64    `json:"quantity"`
		Weight     float64    `json:"weight"`
		DamageDate *time.Time `json:"damage_date,omitempty"`
		Reason     string     `json:"reason"`
		Amount     float64    `json:"amount"`
		Notes      *string    `json:"notes,omitempty"`
		ImageURL   *string    `json:"image_url,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
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
	if req.Quantity == 0 && req.Weight == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "at least one of quantity or weight must be greater than 0")
		return
	}

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "reason is required")
		return
	}

	if req.Amount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount cannot be negative")
		return
	}

	// Store must exist.
	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Lot must exist (for unit derivation).
	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	damageDate := time.Now()
	if req.DamageDate != nil {
		damageDate = *req.DamageDate
	}

	damage := &models.Damage{
		UserID:       userID,
		StoreID:      req.StoreID,
		Quantity:     req.Quantity,
		QuantityUnit: lot.QuantityUnit,
		Weight:       req.Weight,
		WeightUnit:   lot.WeightUnit,
		DamageDate:   damageDate,
		Reason:       reason,
		Amount:       req.Amount,
		Notes:        req.Notes,
		ImageURL:     req.ImageURL,
	}

	id, err := h.app.Models.Damage.Insert(r.Context(), damage)
	if err != nil {
		log.Printf("[DAMAGES] CreateDamageHandler ERROR: failed to insert damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create damage")
		return
	}
	damage.ID = id

	log.Printf("[DAMAGES] CreateDamageHandler SUCCESS: created damage ID=%d store=%d", id, req.StoreID)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created damage #" + strconv.FormatInt(id, 10) + " for store #" + strconv.FormatInt(req.StoreID, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "damage created successfully", damage)
}

// =============================================================================
// READ
// =============================================================================

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
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage")
		return
	}
	if damage == nil {
		log.Printf("[DAMAGES] GetDamageHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "damage not found")
		return
	}

	log.Printf("[DAMAGES] GetDamageHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "damage details fetched", damage)
}

// GetAllDamagesHandler - GET /api/damages
func (h *Handler) GetAllDamagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DAMAGES] GetAllDamagesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	storeIDParam := r.URL.Query().Get("store_id")

	var damages []models.Damage
	var err error

	if storeIDParam != "" {
		id, parseErr := strconv.ParseInt(storeIDParam, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid store_id")
			return
		}
		damages, err = h.app.Models.Damage.GetByStoreID(r.Context(), id)
	} else {
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

// GetStoreDamagesHandler - GET /api/stores/{id}/damages
func (h *Handler) GetStoreDamagesHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] GetStoreDamagesHandler ERROR: invalid store id")
		return
	}

	log.Printf("[DAMAGES] GetStoreDamagesHandler called - StoreID: %d", storeID)

	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[DAMAGES] GetStoreDamagesHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	damages, err := h.app.Models.Damage.GetByStoreID(r.Context(), storeID)
	if err != nil {
		log.Printf("[DAMAGES] GetStoreDamagesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damages")
		return
	}

	if damages == nil {
		damages = []models.Damage{}
	}

	utils.SuccessJson(w, http.StatusOK, "damages fetched successfully", damages)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateDamageHandler - PATCH /api/damages/{id}
//
// store_id and user_id are immutable. If the source is wrong, delete and
// recreate.
func (h *Handler) UpdateDamageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DAMAGES] UpdateDamageHandler called - ID: %d", id)

	type request struct {
		Quantity   *float64   `json:"quantity,omitempty"`
		Weight     *float64   `json:"weight,omitempty"`
		DamageDate *time.Time `json:"damage_date,omitempty"`
		Reason     *string    `json:"reason,omitempty"`
		Amount     *float64   `json:"amount,omitempty"`
		Notes      *string    `json:"notes,omitempty"`
		ImageURL   *string    `json:"image_url,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

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

	if req.Quantity != nil {
		if *req.Quantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
			return
		}
		damage.Quantity = *req.Quantity
	}
	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		damage.Weight = *req.Weight
	}
	if req.DamageDate != nil {
		damage.DamageDate = *req.DamageDate
	}
	if req.Reason != nil {
		trimmed := strings.TrimSpace(*req.Reason)
		if trimmed == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "reason cannot be empty")
			return
		}
		damage.Reason = trimmed
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
	if req.ImageURL != nil {
		damage.ImageURL = req.ImageURL
	}

	if damage.Quantity == 0 && damage.Weight == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "at least one of quantity or weight must be greater than 0")
		return
	}

	if err := h.app.Models.Damage.Update(r.Context(), damage); err != nil {
		log.Printf("[DAMAGES] UpdateDamageHandler ERROR: failed to update damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update damage")
		return
	}

	log.Printf("[DAMAGES] UpdateDamageHandler SUCCESS: updated damage ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated damage #" + strconv.FormatInt(id, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.Damage.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "damage updated successfully", updated)
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteDamageHandler - DELETE /api/damages/{id}
func (h *Handler) DeleteDamageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DAMAGES] DeleteDamageHandler called - ID: %d", id)

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

	if err := h.app.Models.Damage.Delete(r.Context(), id); err != nil {
		log.Printf("[DAMAGES] DeleteDamageHandler ERROR: failed to delete damage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete damage")
		return
	}

	log.Printf("[DAMAGES] DeleteDamageHandler SUCCESS: deleted damage ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted damage #" + strconv.FormatInt(id, 10),
		EntityType:  "damages",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "damage deleted successfully", nil)
}
