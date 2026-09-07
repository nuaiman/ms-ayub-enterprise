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
// LOT MANAGEMENT
// =============================================================================

// CreateLotHandler - POST /api/lots
func (h *Handler) CreateLotHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOTS] CreateLotHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		ItemID                   int64      `json:"item_id"`
		LotNumber                int64      `json:"lot_number"`
		CustomerChargeType       string     `json:"customer_charge_type"`
		MajhiBillType            string     `json:"majhi_bill_type"`
		CustomerStorageRate      float64    `json:"customer_storage_rate"`
		UnloadRate               float64    `json:"unload_rate"`
		MajhiID                  *int64     `json:"majhi_id,omitempty"`
		MajhiCut                 float64    `json:"majhi_cut"`
		Notes                    *string    `json:"notes,omitempty"`
		CustomerLastPaidThrough  *time.Time `json:"customer_last_paid_through,omitempty"`
		CustomerLastPaidAmount   float64    `json:"customer_last_paid_amount"`
		CustomerPaidUnloadAmount float64    `json:"customer_paid_unload_amount,omitempty"`
		MajhiTotalPaid           float64    `json:"majhi_total_paid"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.ItemID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "item_id is required")
		return
	}
	if req.LotNumber <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "lot_number must be greater than 0")
		return
	}

	// Validate customer_charge_type
	if req.CustomerChargeType != "weight" && req.CustomerChargeType != "quantity" {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_charge_type must be 'weight' or 'quantity'")
		return
	}

	// Validate majhi_bill_type
	if req.MajhiBillType != "weight" && req.MajhiBillType != "quantity" && req.MajhiBillType != "job" {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_bill_type must be 'weight', 'quantity', or 'job'")
		return
	}

	// Validate amounts
	if req.CustomerLastPaidAmount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_last_paid_amount cannot be negative")
		return
	}
	if req.CustomerPaidUnloadAmount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_paid_unload_amount cannot be negative")
		return
	}
	if req.MajhiTotalPaid < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_total_paid cannot be negative")
		return
	}

	// Check if item exists
	exists, err := h.app.Models.Item.Exists(r.Context(), req.ItemID)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to verify item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify item")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	// Check if majhi exists if provided
	if req.MajhiID != nil {
		exists, err := h.app.Models.Majhi.Exists(r.Context(), *req.MajhiID)
		if err != nil {
			log.Printf("[LOTS] CreateLotHandler ERROR: failed to verify majhi - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
			return
		}
	}

	// Check if lot already exists for this item
	exists, err = h.app.Models.Lot.ExistsByItemAndLot(r.Context(), req.ItemID, req.LotNumber)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to check existing lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "lot already exists for this item")
		return
	}

	// Get current user ID for audit
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[LOTS] CreateLotHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	lot := &models.Lot{
		ItemID:                   req.ItemID,
		LotNumber:                req.LotNumber,
		CustomerChargeType:       req.CustomerChargeType,
		MajhiBillType:            req.MajhiBillType,
		CustomerStorageRate:      req.CustomerStorageRate,
		UnloadRate:               req.UnloadRate,
		MajhiID:                  req.MajhiID,
		MajhiCut:                 req.MajhiCut,
		IsActive:                 true,
		Notes:                    req.Notes,
		CustomerLastPaidThrough:  req.CustomerLastPaidThrough,
		CustomerLastPaidAmount:   req.CustomerLastPaidAmount,
		CustomerPaidUnloadAmount: req.CustomerPaidUnloadAmount,
		MajhiTotalPaid:           req.MajhiTotalPaid,
	}

	id, err := h.app.Models.Lot.Insert(r.Context(), lot)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to insert lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create lot")
		return
	}

	lot.ID = id

	log.Printf("[LOTS] CreateLotHandler SUCCESS: created lot ID=%d for item ID=%d", id, req.ItemID)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created lot #" + strconv.FormatInt(id, 10) + " for item #" + strconv.FormatInt(req.ItemID, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "lot created successfully", lot)
}

// GetLotHandler - GET /api/lots/{id}
func (h *Handler) GetLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] GetLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] GetLotHandler called - ID: %d", id)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] GetLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}

	if lot == nil {
		log.Printf("[LOTS] GetLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	log.Printf("[LOTS] GetLotHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "lot details fetched", lot)
}

// GetAllLotsHandler - GET /api/lots
func (h *Handler) GetAllLotsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOTS] GetAllLotsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	itemID := r.URL.Query().Get("item_id")
	activeOnly := r.URL.Query().Get("active") == "true"

	var lots []models.Lot
	var err error

	switch {
	case search != "":
		lots, err = h.app.Models.Lot.Search(r.Context(), search)
	case itemID != "":
		id, _ := strconv.ParseInt(itemID, 10, 64)
		lots, err = h.app.Models.Lot.GetByItemID(r.Context(), id)
	case activeOnly:
		lots, err = h.app.Models.Lot.GetActive(r.Context())
	default:
		lots, err = h.app.Models.Lot.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[LOTS] GetAllLotsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lots")
		return
	}

	if lots == nil {
		lots = []models.Lot{}
	}

	log.Printf("[LOTS] GetAllLotsHandler SUCCESS: fetched %d lots", len(lots))
	utils.SuccessJson(w, http.StatusOK, "lots fetched successfully", lots)
}

// UpdateLotHandler - PATCH /api/lots/{id}
func (h *Handler) UpdateLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] UpdateLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] UpdateLotHandler called - ID: %d", id)

	type request struct {
		CustomerChargeType       *string    `json:"customer_charge_type,omitempty"`
		MajhiBillType            *string    `json:"majhi_bill_type,omitempty"`
		CustomerStorageRate      *float64   `json:"customer_storage_rate,omitempty"`
		UnloadRate               *float64   `json:"unload_rate,omitempty"`
		MajhiID                  *int64     `json:"majhi_id,omitempty"`
		MajhiCut                 *float64   `json:"majhi_cut,omitempty"`
		Notes                    *string    `json:"notes,omitempty"`
		CustomerLastPaidThrough  *time.Time `json:"customer_last_paid_through,omitempty"`
		CustomerLastPaidAmount   *float64   `json:"customer_last_paid_amount,omitempty"`
		CustomerPaidUnloadAmount *float64   `json:"customer_paid_unload_amount,omitempty"`
		MajhiTotalPaid           *float64   `json:"majhi_total_paid,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing lot
	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] UpdateLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// Check if majhi exists if provided
	if req.MajhiID != nil {
		exists, err := h.app.Models.Majhi.Exists(r.Context(), *req.MajhiID)
		if err != nil {
			log.Printf("[LOTS] UpdateLotHandler ERROR: failed to verify majhi - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
			return
		}
	}

	// Apply updates
	if req.CustomerChargeType != nil {
		if *req.CustomerChargeType != "weight" && *req.CustomerChargeType != "quantity" {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_charge_type must be 'weight' or 'quantity'")
			return
		}
		lot.CustomerChargeType = *req.CustomerChargeType
	}
	if req.MajhiBillType != nil {
		if *req.MajhiBillType != "weight" && *req.MajhiBillType != "quantity" && *req.MajhiBillType != "job" {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_bill_type must be 'weight', 'quantity', or 'job'")
			return
		}
		lot.MajhiBillType = *req.MajhiBillType
	}
	if req.CustomerStorageRate != nil {
		if *req.CustomerStorageRate < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_storage_rate cannot be negative")
			return
		}
		lot.CustomerStorageRate = *req.CustomerStorageRate
	}
	if req.UnloadRate != nil {
		if *req.UnloadRate < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "unload_rate cannot be negative")
			return
		}
		lot.UnloadRate = *req.UnloadRate
	}
	if req.MajhiID != nil {
		lot.MajhiID = req.MajhiID
	}
	if req.MajhiCut != nil {
		if *req.MajhiCut < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_cut cannot be negative")
			return
		}
		lot.MajhiCut = *req.MajhiCut
	}
	if req.Notes != nil {
		lot.Notes = req.Notes
	}
	if req.CustomerLastPaidThrough != nil {
		lot.CustomerLastPaidThrough = req.CustomerLastPaidThrough
	}
	if req.CustomerLastPaidAmount != nil {
		if *req.CustomerLastPaidAmount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_last_paid_amount cannot be negative")
			return
		}
		lot.CustomerLastPaidAmount = *req.CustomerLastPaidAmount
	}
	if req.CustomerPaidUnloadAmount != nil {
		if *req.CustomerPaidUnloadAmount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_paid_unload_amount cannot be negative")
			return
		}
		lot.CustomerPaidUnloadAmount = *req.CustomerPaidUnloadAmount
	}
	if req.MajhiTotalPaid != nil {
		if *req.MajhiTotalPaid < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_total_paid cannot be negative")
			return
		}
		lot.MajhiTotalPaid = *req.MajhiTotalPaid
	}

	if err := h.app.Models.Lot.Update(r.Context(), lot); err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: failed to update lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update lot")
		return
	}

	log.Printf("[LOTS] UpdateLotHandler SUCCESS: updated lot ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated lot
	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "lot updated successfully", updatedLot)
}

// ToggleLotActiveHandler - PATCH /api/lots/{id}/toggle-active
func (h *Handler) ToggleLotActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] ToggleLotActiveHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] ToggleLotActiveHandler called - ID: %d", id)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] ToggleLotActiveHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] ToggleLotActiveHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	newStatus := !lot.IsActive

	if err := h.app.Models.Lot.ToggleActive(r.Context(), id, newStatus); err != nil {
		log.Printf("[LOTS] ToggleLotActiveHandler ERROR: failed to toggle status - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update lot status")
		return
	}

	log.Printf("[LOTS] ToggleLotActiveHandler SUCCESS: lot ID=%d now active=%v", id, newStatus)

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
		Description: "Lot " + action + ": #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	lot.IsActive = newStatus
	utils.SuccessJson(w, http.StatusOK, "lot status updated successfully", lot)
}

// DeleteLotHandler - DELETE /api/lots/{id}
func (h *Handler) DeleteLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] DeleteLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] DeleteLotHandler called - ID: %d", id)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] DeleteLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] DeleteLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if err := h.app.Models.Lot.Delete(r.Context(), id); err != nil {
		log.Printf("[LOTS] DeleteLotHandler ERROR: failed to delete lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete lot")
		return
	}

	log.Printf("[LOTS] DeleteLotHandler SUCCESS: deleted lot ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "lot deleted successfully", nil)
}

// =============================================================================
// PAYMENT HANDLERS
// =============================================================================

// UpdateLotCustomerPaymentHandler - PATCH /api/lots/{id}/customer-payment
// Updates customer payment tracking (storage bills)
func (h *Handler) UpdateLotCustomerPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] UpdateLotCustomerPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] UpdateLotCustomerPaymentHandler called - ID: %d", id)

	type request struct {
		CustomerLastPaidThrough *time.Time `json:"customer_last_paid_through,omitempty"`
		CustomerLastPaidAmount  float64    `json:"customer_last_paid_amount"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] UpdateLotCustomerPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerLastPaidAmount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_last_paid_amount cannot be negative")
		return
	}

	// Get existing lot
	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] UpdateLotCustomerPaymentHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] UpdateLotCustomerPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if err := h.app.Models.Lot.UpdateCustomerPayment(r.Context(), id, req.CustomerLastPaidThrough, req.CustomerLastPaidAmount); err != nil {
		log.Printf("[LOTS] UpdateLotCustomerPaymentHandler ERROR: failed to update payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update payment")
		return
	}

	log.Printf("[LOTS] UpdateLotCustomerPaymentHandler SUCCESS: updated customer payment for lot ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer payment for lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer payment updated successfully", updatedLot)
}

// UpdateLotCustomerUnloadPaymentHandler - PATCH /api/lots/{id}/customer-unload-payment
// Updates customer unload payment amount (for one-time unload bills)
func (h *Handler) UpdateLotCustomerUnloadPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler called - ID: %d", id)

	type request struct {
		CustomerPaidUnloadAmount float64    `json:"customer_paid_unload_amount"`
		PaidThrough              *time.Time `json:"paid_through,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerPaidUnloadAmount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_paid_unload_amount cannot be negative")
		return
	}

	// Get existing lot
	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if err := h.app.Models.Lot.UpdateCustomerUnloadPayment(r.Context(), id, req.CustomerPaidUnloadAmount); err != nil {
		log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler ERROR: failed to update unload payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update unload payment")
		return
	}

	// Also update customer_last_paid_through if provided
	if req.PaidThrough != nil {
		// Get current lot data to preserve customer_last_paid_amount
		current, err := h.app.Models.Lot.GetByID(r.Context(), id)
		if err == nil && current != nil {
			_ = h.app.Models.Lot.UpdateCustomerPayment(r.Context(), id, req.PaidThrough, current.CustomerLastPaidAmount)
		}
	}

	log.Printf("[LOTS] UpdateLotCustomerUnloadPaymentHandler SUCCESS: updated customer unload payment for lot ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer unload payment for lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer unload payment updated successfully", updatedLot)
}

// UpdateLotMajhiPaymentHandler - PATCH /api/lots/{id}/majhi-payment
// Updates majhi payment tracking
func (h *Handler) UpdateLotMajhiPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] UpdateLotMajhiPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] UpdateLotMajhiPaymentHandler called - ID: %d", id)

	type request struct {
		MajhiTotalPaid float64 `json:"majhi_total_paid"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] UpdateLotMajhiPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.MajhiTotalPaid < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_total_paid cannot be negative")
		return
	}

	// Get existing lot
	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] UpdateLotMajhiPaymentHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] UpdateLotMajhiPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if err := h.app.Models.Lot.UpdateMajhiPayment(r.Context(), id, req.MajhiTotalPaid); err != nil {
		log.Printf("[LOTS] UpdateLotMajhiPaymentHandler ERROR: failed to update majhi payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update majhi payment")
		return
	}

	log.Printf("[LOTS] UpdateLotMajhiPaymentHandler SUCCESS: updated majhi payment for lot ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated majhi payment for lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "majhi payment updated successfully", updatedLot)
}
