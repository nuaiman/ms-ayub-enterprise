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
// STORE TRANSFERS
// =============================================================================

// CreateStoreTransferHandler - POST /api/stores/{id}/transfer
//
// Moves a store from its current godown to a new godown. Bills are not touched.
// The store must be active.
func (h *Handler) CreateStoreTransferHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: invalid store id")
		return
	}

	log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler called - StoreID: %d", storeID)

	type request struct {
		ToGodownID int64   `json:"to_godown_id"`
		Notes      *string `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ToGodownID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "to_godown_id is required")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	if !store.IsActive {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot transfer an inactive store")
		return
	}

	if store.GodownID == req.ToGodownID {
		utils.ErrorJson(w, http.StatusBadRequest, "store is already at this godown")
		return
	}

	toGodown, err := h.app.Models.Godown.GetByID(r.Context(), req.ToGodownID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: failed to verify target godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify target godown")
		return
	}
	if toGodown == nil {
		utils.ErrorJson(w, http.StatusNotFound, "target godown not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	transferredAt := time.Now()
	fromGodownID := store.GodownID

	transfer := &models.StoreTransfer{
		UserID:             userID,
		StoreID:            storeID,
		FromGodownID:       fromGodownID,
		ToGodownID:         req.ToGodownID,
		WeightAtTransfer:   store.Weight,
		QuantityAtTransfer: store.Quantity,
		Notes:              req.Notes,
		TransferredAt:      transferredAt,
	}

	id, err := h.app.Models.StoreTransfer.Insert(r.Context(), transfer)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: failed to insert transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create transfer")
		return
	}
	transfer.ID = id

	// Mutate store godown.
	store.GodownID = req.ToGodownID
	if err := h.app.Models.Store.UpdateGodown(r.Context(), store.ID, store.GodownID); err != nil {
		log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler ERROR: failed to update store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "transfer recorded but store update failed")
		return
	}

	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), storeID)

	log.Printf("[STORE_TRANSFERS] CreateStoreTransferHandler SUCCESS: transfer ID=%d store=%d %d→%d",
		id, storeID, fromGodownID, req.ToGodownID)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "create",
		Description: "Transferred store #" + strconv.FormatInt(storeID, 10) +
			" from godown #" + strconv.FormatInt(fromGodownID, 10) +
			" to godown #" + strconv.FormatInt(req.ToGodownID, 10),
		EntityType: "store_transfers",
		EntityID:   id,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "store transferred successfully", map[string]any{
		"transfer": transfer,
		"store":    updatedStore,
	})
}

// =============================================================================
// READ
// =============================================================================

// GetStoreTransferHandler - GET /api/store-transfers/{id}
func (h *Handler) GetStoreTransferHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_TRANSFERS] GetStoreTransferHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORE_TRANSFERS] GetStoreTransferHandler called - ID: %d", id)

	transfer, err := h.app.Models.StoreTransfer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] GetStoreTransferHandler ERROR: failed to fetch transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfer")
		return
	}
	if transfer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "transfer not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "transfer details fetched", transfer)
}

// GetStoreTransfersHandler - GET /api/stores/{id}/transfers
func (h *Handler) GetStoreTransfersHandler(w http.ResponseWriter, r *http.Request) {
	storeID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_TRANSFERS] GetStoreTransfersHandler ERROR: invalid store id")
		return
	}

	log.Printf("[STORE_TRANSFERS] GetStoreTransfersHandler called - StoreID: %d", storeID)

	store, err := h.app.Models.Store.GetByID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] GetStoreTransfersHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	transfers, err := h.app.Models.StoreTransfer.GetByStoreID(r.Context(), storeID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] GetStoreTransfersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfers")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "transfers fetched successfully", transfers)
}

// GetAllStoreTransfersHandler - GET /api/store-transfers
func (h *Handler) GetAllStoreTransfersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[STORE_TRANSFERS] GetAllStoreTransfersHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	userIDParam := r.URL.Query().Get("user_id")

	var transfers []models.StoreTransfer
	var err error

	if userIDParam != "" {
		id, parseErr := strconv.ParseInt(userIDParam, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		transfers, err = h.app.Models.StoreTransfer.GetByUserID(r.Context(), id)
	} else {
		transfers, err = h.app.Models.StoreTransfer.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[STORE_TRANSFERS] GetAllStoreTransfersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfers")
		return
	}

	if transfers == nil {
		transfers = []models.StoreTransfer{}
	}

	utils.SuccessJson(w, http.StatusOK, "transfers fetched successfully", transfers)
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteStoreTransferHandler - DELETE /api/store-transfers/{id}
//
// Reverses the transfer: moves the store back to from_godown_id, then deletes
// the transfer row. Rejected if a newer transfer exists for the same store,
// since reverting out-of-order would corrupt the transfer chain.
func (h *Handler) DeleteStoreTransferHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler called - ID: %d", id)

	transfer, err := h.app.Models.StoreTransfer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: failed to fetch transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfer")
		return
	}
	if transfer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "transfer not found")
		return
	}

	// Only the most recent transfer for a store may be deleted.
	latest, err := h.app.Models.StoreTransfer.GetLatestForStore(r.Context(), transfer.StoreID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: failed to fetch latest transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify transfer chain")
		return
	}
	if latest == nil || latest.ID != transfer.ID {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete; a later transfer exists for this store")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), transfer.StoreID)
	if err != nil {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Move store back first.
	store.GodownID = transfer.FromGodownID
	if err := h.app.Models.Store.UpdateGodown(r.Context(), store.ID, store.GodownID); err != nil {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: failed to move store back - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to reverse transfer on store")
		return
	}

	if err := h.app.Models.StoreTransfer.Delete(r.Context(), id); err != nil {
		log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler ERROR: failed to delete transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete transfer")
		return
	}

	updatedStore, _ := h.app.Models.Store.GetByID(r.Context(), transfer.StoreID)

	log.Printf("[STORE_TRANSFERS] DeleteStoreTransferHandler SUCCESS: deleted transfer ID=%d store=%d → back to godown %d",
		id, transfer.StoreID, transfer.FromGodownID)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "delete",
		Description: "Deleted store transfer #" + strconv.FormatInt(id, 10) +
			" (store #" + strconv.FormatInt(transfer.StoreID, 10) +
			" moved back to godown #" + strconv.FormatInt(transfer.FromGodownID, 10) + ")",
		EntityType: "store_transfers",
		EntityID:   id,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "store transfer deleted successfully", map[string]any{
		"store": updatedStore,
	})
}
