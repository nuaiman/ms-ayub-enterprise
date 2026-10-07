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
// LOT TRANSFERS (ACCOUNT TRANSFERS)
// =============================================================================

// CreateLotTransferHandler - POST /api/lots/{id}/transfer
//
// Reassigns a lot to a different customer. Optionally creates a customer
// store bill for the current month for each active store under the lot, using
// the bill_type and rate provided in the request.
//
// Old customer's existing bills are not touched.
func (h *Handler) CreateLotTransferHandler(w http.ResponseWriter, r *http.Request) {
	lotID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: invalid lot id")
		return
	}

	log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler called - LotID: %d", lotID)

	type request struct {
		ToCustomerID int64   `json:"to_customer_id"`
		Notes        *string `json:"notes,omitempty"`

		// Optional billing block. If BillType is set and Rate > 0, we create
		// a customer_store_bill for the current month for each active store
		// under the lot, under the new customer.
		BillType *string  `json:"bill_type,omitempty"`
		Rate     *float64 `json:"rate,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ToCustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "to_customer_id is required")
		return
	}

	// Validate optional billing block.
	wantsBilling := false
	billType := "quantity"
	var rate float64

	if req.BillType != nil && *req.BillType != "" {
		if *req.BillType != "weight" && *req.BillType != "quantity" {
			utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight' or 'quantity'")
			return
		}
		billType = *req.BillType
	}
	if req.Rate != nil {
		if *req.Rate < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "rate cannot be negative")
			return
		}
		rate = *req.Rate
	}
	// Billing is only applied if a positive rate is provided.
	if rate > 0 {
		wantsBilling = true
	}

	// Load lot.
	lot, err := h.app.Models.Lot.GetByID(r.Context(), lotID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if lot.CustomerID == req.ToCustomerID {
		utils.ErrorJson(w, http.StatusBadRequest, "lot already belongs to this customer")
		return
	}

	// Verify target customer.
	toCustomer, err := h.app.Models.Customer.GetByID(r.Context(), req.ToCustomerID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to verify target customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify target customer")
		return
	}
	if toCustomer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "target customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	transferredAt := time.Now()
	fromCustomerID := lot.CustomerID

	// 1. Insert transfer record.
	transfer := &models.LotTransfer{
		UserID:         userID,
		LotID:          lotID,
		FromCustomerID: fromCustomerID,
		ToCustomerID:   req.ToCustomerID,
		Notes:          req.Notes,
		TransferredAt:  transferredAt,
	}

	id, err := h.app.Models.LotTransfer.Insert(r.Context(), transfer)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to insert transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create transfer")
		return
	}
	transfer.ID = id

	// 2. Move lot to new customer.
	lot.CustomerID = req.ToCustomerID
	if err := h.app.Models.Lot.Update(r.Context(), lot); err != nil {
		log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to update lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "transfer recorded but lot update failed")
		return
	}

	// 3. Optionally create customer store bills for the current month for
	//    each active store under the lot.
	createdBills := 0
	if wantsBilling {
		currentMonth := time.Now().Format("2006-01")

		stores, err := h.app.Models.Store.GetByLotID(r.Context(), lotID)
		if err != nil {
			log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to fetch stores - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "transfer succeeded but failed to fetch stores for billing")
			return
		}

		for _, store := range stores {
			if !store.IsActive {
				continue
			}

			// Skip if a bill already exists for this customer+store+month.
			exists, err := h.app.Models.CustomerStoreBill.ExistsForCustomerStoreMonth(
				r.Context(), req.ToCustomerID, store.ID, currentMonth,
			)
			if err != nil {
				log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to check existing bill - %v", err)
				continue
			}
			if exists {
				continue
			}

			total := billTotalFromSnapshot(billType, rate, store.Weight, store.Quantity)

			newBill := &models.CustomerStoreBill{
				UserID:                userID,
				CustomerID:            req.ToCustomerID,
				StoreID:               store.ID,
				MonthYear:             currentMonth,
				BillType:              billType,
				Rate:                  rate,
				WeightAtBilling:       store.Weight,
				QuantityAtBilling:     store.Quantity,
				WeightUnitAtBilling:   &lot.WeightUnit,
				QuantityUnitAtBilling: &lot.QuantityUnit,
				TotalAmount:           total,
			}

			if _, err := h.app.Models.CustomerStoreBill.Insert(r.Context(), newBill); err != nil {
				log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler ERROR: failed to create bill for store %d - %v", store.ID, err)
				continue
			}

			createdBills++
		}
	}

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), lotID)

	log.Printf("[LOT_TRANSFERS] CreateLotTransferHandler SUCCESS: transfer ID=%d lot=%d %d→%d bills_created=%d",
		id, lotID, fromCustomerID, req.ToCustomerID, createdBills)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "create",
		Description: "Transferred lot #" + strconv.FormatInt(lotID, 10) +
			" from customer #" + strconv.FormatInt(fromCustomerID, 10) +
			" to customer #" + strconv.FormatInt(req.ToCustomerID, 10),
		EntityType: "lot_transfers",
		EntityID:   id,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "lot transferred successfully", map[string]any{
		"transfer":      transfer,
		"lot":           updatedLot,
		"bills_created": createdBills,
	})
}

// =============================================================================
// READ
// =============================================================================

// GetLotTransferHandler - GET /api/lot-transfers/{id}
func (h *Handler) GetLotTransferHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOT_TRANSFERS] GetLotTransferHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOT_TRANSFERS] GetLotTransferHandler called - ID: %d", id)

	transfer, err := h.app.Models.LotTransfer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] GetLotTransferHandler ERROR: failed to fetch transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfer")
		return
	}
	if transfer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "transfer not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "transfer details fetched", transfer)
}

// GetLotTransfersHandler - GET /api/lots/{id}/transfers
func (h *Handler) GetLotTransfersHandler(w http.ResponseWriter, r *http.Request) {
	lotID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOT_TRANSFERS] GetLotTransfersHandler ERROR: invalid lot id")
		return
	}

	log.Printf("[LOT_TRANSFERS] GetLotTransfersHandler called - LotID: %d", lotID)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), lotID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] GetLotTransfersHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	transfers, err := h.app.Models.LotTransfer.GetByLotID(r.Context(), lotID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] GetLotTransfersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfers")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "transfers fetched successfully", transfers)
}

// GetAllLotTransfersHandler - GET /api/lot-transfers
func (h *Handler) GetAllLotTransfersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOT_TRANSFERS] GetAllLotTransfersHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	userIDParam := r.URL.Query().Get("user_id")

	var transfers []models.LotTransfer
	var err error

	if userIDParam != "" {
		id, parseErr := strconv.ParseInt(userIDParam, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		transfers, err = h.app.Models.LotTransfer.GetByUserID(r.Context(), id)
	} else {
		transfers, err = h.app.Models.LotTransfer.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[LOT_TRANSFERS] GetAllLotTransfersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfers")
		return
	}

	if transfers == nil {
		transfers = []models.LotTransfer{}
	}

	utils.SuccessJson(w, http.StatusOK, "transfers fetched successfully", transfers)
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteLotTransferHandler - DELETE /api/lot-transfers/{id}
//
// Reverses the transfer: sets the lot's customer_id back to from_customer_id,
// then deletes the transfer row. Rejected if a newer transfer exists for the
// same lot, since reverting out-of-order would corrupt the transfer chain.
//
// Note: any customer store bills created by the transfer are NOT removed here.
// Delete them manually from the Customer Store Bills page if desired.
func (h *Handler) DeleteLotTransferHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler called - ID: %d", id)

	transfer, err := h.app.Models.LotTransfer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: failed to fetch transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transfer")
		return
	}
	if transfer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "transfer not found")
		return
	}

	// Only the most recent transfer for a lot may be deleted.
	latest, err := h.app.Models.LotTransfer.GetLatestForLot(r.Context(), transfer.LotID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: failed to fetch latest transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify transfer chain")
		return
	}
	if latest == nil || latest.ID != transfer.ID {
		utils.ErrorJson(w, http.StatusBadRequest, "cannot delete; a later transfer exists for this lot")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), transfer.LotID)
	if err != nil {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// Move lot back to its previous customer.
	lot.CustomerID = transfer.FromCustomerID
	if err := h.app.Models.Lot.Update(r.Context(), lot); err != nil {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: failed to move lot back - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to reverse transfer on lot")
		return
	}

	if err := h.app.Models.LotTransfer.Delete(r.Context(), id); err != nil {
		log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler ERROR: failed to delete transfer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete transfer")
		return
	}

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), transfer.LotID)

	log.Printf("[LOT_TRANSFERS] DeleteLotTransferHandler SUCCESS: deleted transfer ID=%d lot=%d → back to customer %d",
		id, transfer.LotID, transfer.FromCustomerID)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "delete",
		Description: "Deleted lot transfer #" + strconv.FormatInt(id, 10) +
			" (lot #" + strconv.FormatInt(transfer.LotID, 10) +
			" moved back to customer #" + strconv.FormatInt(transfer.FromCustomerID, 10) + ")",
		EntityType: "lot_transfers",
		EntityID:   id,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "lot transfer deleted successfully", map[string]any{
		"lot": updatedLot,
	})
}
