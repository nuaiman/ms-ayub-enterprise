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
// RESPONSE SHAPES
// =============================================================================

type MajhiBillResponse struct {
	models.MajhiBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

// CreateMajhiBillHandler - POST /api/majhi-bills
//
// Requires exactly one of store_id or delivery_item_id.
func (h *Handler) CreateMajhiBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		MajhiID        int64   `json:"majhi_id"`
		StoreID        *int64  `json:"store_id,omitempty"`
		DeliveryItemID *int64  `json:"delivery_item_id,omitempty"`
		BillType       string  `json:"bill_type"`
		Rate           float64 `json:"rate"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.MajhiID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_id is required")
		return
	}

	hasStore := req.StoreID != nil && *req.StoreID != 0
	hasDeliveryItem := req.DeliveryItemID != nil && *req.DeliveryItemID != 0

	if hasStore == hasDeliveryItem {
		utils.ErrorJson(w, http.StatusBadRequest, "exactly one of store_id or delivery_item_id is required")
		return
	}

	if req.BillType == "" {
		req.BillType = "quantity"
	}
	if req.BillType != "weight" && req.BillType != "quantity" && req.BillType != "job" {
		utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight', 'quantity', or 'job'")
		return
	}

	if req.Rate < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "rate cannot be negative")
		return
	}

	majhiExists, err := h.app.Models.Majhi.Exists(r.Context(), req.MajhiID)
	if err != nil {
		log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to verify majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
		return
	}
	if !majhiExists {
		utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
		return
	}

	var (
		storeID        *int64
		deliveryItemID *int64
		weightAtBill   float64
		quantityAtBill float64
		weightUnit     *string
		quantityUnit   *string
	)

	if hasStore {
		store, err := h.app.Models.Store.GetByID(r.Context(), *req.StoreID)
		if err != nil {
			log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to verify store - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
			return
		}
		if store == nil {
			utils.ErrorJson(w, http.StatusNotFound, "store not found")
			return
		}

		lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
		if err != nil {
			log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to fetch lot - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
			return
		}
		if lot == nil {
			utils.ErrorJson(w, http.StatusNotFound, "lot not found")
			return
		}

		storeID = req.StoreID
		weightAtBill = store.Weight
		quantityAtBill = store.Quantity
		weightUnit = &lot.WeightUnit
		quantityUnit = &lot.QuantityUnit
	} else {
		item, err := h.app.Models.Delivery.GetItemByID(r.Context(), *req.DeliveryItemID)
		if err != nil {
			log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to verify delivery item - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify delivery item")
			return
		}
		if item == nil {
			utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
			return
		}

		store, err := h.app.Models.Store.GetByID(r.Context(), item.StoreID)
		if err != nil {
			log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to fetch store for delivery item - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
			return
		}
		if store == nil {
			utils.ErrorJson(w, http.StatusNotFound, "store not found for delivery item")
			return
		}

		lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
		if err != nil {
			log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to fetch lot - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
			return
		}
		if lot == nil {
			utils.ErrorJson(w, http.StatusNotFound, "lot not found")
			return
		}

		deliveryItemID = req.DeliveryItemID
		weightAtBill = item.Weight
		quantityAtBill = item.Quantity
		weightUnit = &lot.WeightUnit
		quantityUnit = &lot.QuantityUnit
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	total := billTotalFromSnapshot(req.BillType, req.Rate, weightAtBill, quantityAtBill)

	bill := &models.MajhiBill{
		UserID:                userID,
		MajhiID:               req.MajhiID,
		StoreID:               storeID,
		DeliveryItemID:        deliveryItemID,
		BillType:              req.BillType,
		Rate:                  req.Rate,
		WeightAtBilling:       weightAtBill,
		QuantityAtBilling:     quantityAtBill,
		WeightUnitAtBilling:   weightUnit,
		QuantityUnitAtBilling: quantityUnit,
		TotalAmount:           total,
	}

	id, err := h.app.Models.MajhiBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create majhi bill")
		return
	}

	bill.ID = id

	log.Printf("[MAJHI_BILLS] CreateMajhiBillHandler SUCCESS: created bill ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created majhi bill #" + strconv.FormatInt(id, 10),
		EntityType:  "majhi_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "majhi bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

func (h *Handler) GetMajhiBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHI_BILLS] GetMajhiBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHI_BILLS] GetMajhiBillHandler called - ID: %d", id)

	bill, err := h.app.Models.MajhiBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHI_BILLS] GetMajhiBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi bill")
		return
	}
	if bill == nil {
		log.Printf("[MAJHI_BILLS] GetMajhiBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi bill not found")
		return
	}

	responses, err := h.buildMajhiBillResponses(r, []models.MajhiBill{*bill})
	if err != nil {
		log.Printf("[MAJHI_BILLS] GetMajhiBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[MAJHI_BILLS] GetMajhiBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "majhi bill details fetched", responses[0])
}

func (h *Handler) GetAllMajhiBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MAJHI_BILLS] GetAllMajhiBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	storeID := r.URL.Query().Get("store_id")
	majhiID := r.URL.Query().Get("majhi_id")
	deliveryItemID := r.URL.Query().Get("delivery_item_id")

	var bills []models.MajhiBill
	var err error

	switch {
	case storeID != "":
		id, parseErr := strconv.ParseInt(storeID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid store_id")
			return
		}
		bills, err = h.app.Models.MajhiBill.GetByStoreID(r.Context(), id)
	case deliveryItemID != "":
		id, parseErr := strconv.ParseInt(deliveryItemID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid delivery_item_id")
			return
		}
		bills, err = h.app.Models.MajhiBill.GetByDeliveryItemID(r.Context(), id)
	case majhiID != "":
		id, parseErr := strconv.ParseInt(majhiID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid majhi_id")
			return
		}
		bills, err = h.app.Models.MajhiBill.GetByMajhiID(r.Context(), id)
	default:
		bills, err = h.app.Models.MajhiBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[MAJHI_BILLS] GetAllMajhiBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi bills")
		return
	}

	responses, err := h.buildMajhiBillResponses(r, bills)
	if err != nil {
		log.Printf("[MAJHI_BILLS] GetAllMajhiBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[MAJHI_BILLS] GetAllMajhiBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "majhi bills fetched successfully", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

func (h *Handler) UpdateMajhiBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler called - ID: %d", id)

	type request struct {
		BillType          *string  `json:"bill_type,omitempty"`
		Rate              *float64 `json:"rate,omitempty"`
		WeightAtBilling   *float64 `json:"weight_at_billing,omitempty"`
		QuantityAtBilling *float64 `json:"quantity_at_billing,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.MajhiBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi bill")
		return
	}
	if bill == nil {
		log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi bill not found")
		return
	}

	if req.BillType != nil {
		if *req.BillType != "weight" && *req.BillType != "quantity" && *req.BillType != "job" {
			utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight', 'quantity', or 'job'")
			return
		}
		bill.BillType = *req.BillType
	}
	if req.Rate != nil {
		if *req.Rate < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "rate cannot be negative")
			return
		}
		bill.Rate = *req.Rate
	}
	if req.WeightAtBilling != nil {
		if *req.WeightAtBilling < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight_at_billing cannot be negative")
			return
		}
		bill.WeightAtBilling = *req.WeightAtBilling
	}
	if req.QuantityAtBilling != nil {
		if *req.QuantityAtBilling < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity_at_billing cannot be negative")
			return
		}
		bill.QuantityAtBilling = *req.QuantityAtBilling
	}

	bill.TotalAmount = billTotalFromSnapshot(bill.BillType, bill.Rate, bill.WeightAtBilling, bill.QuantityAtBilling)

	if err := h.app.Models.MajhiBill.Update(r.Context(), bill); err != nil {
		log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update majhi bill")
		return
	}

	log.Printf("[MAJHI_BILLS] UpdateMajhiBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated majhi bill #" + strconv.FormatInt(id, 10),
		EntityType:  "majhi_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.MajhiBill.GetByID(r.Context(), id)
	responses, _ := h.buildMajhiBillResponses(r, []models.MajhiBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "majhi bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

func (h *Handler) DeleteMajhiBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler called - ID: %d", id)

	bill, err := h.app.Models.MajhiBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi bill")
		return
	}
	if bill == nil {
		log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "majhi", id)
	if err != nil {
		log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.MajhiBill.Delete(r.Context(), id); err != nil {
		log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete majhi bill")
		return
	}

	log.Printf("[MAJHI_BILLS] DeleteMajhiBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted majhi bill #" + strconv.FormatInt(id, 10),
		EntityType:  "majhi_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "majhi bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildMajhiBillResponses(r *http.Request, bills []models.MajhiBill) ([]MajhiBillResponse, error) {
	if len(bills) == 0 {
		return []MajhiBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "majhi", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]MajhiBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.TotalAmount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, MajhiBillResponse{
			MajhiBill:        b,
			TotalPaid:        t.TotalPaid,
			TotalPaidThrough: t.TotalPaidThrough,
			Remaining:        remaining,
		})
	}

	return out, nil
}