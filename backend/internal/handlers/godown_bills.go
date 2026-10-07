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

type GodownBillResponse struct {
	models.GodownBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

// CreateGodownBillHandler - POST /api/godown-bills
func (h *Handler) CreateGodownBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GODOWN_BILLS] CreateGodownBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		GodownID  int64   `json:"godown_id"`
		StoreID   int64   `json:"store_id"`
		MonthYear string  `json:"month_year"` // optional — derived from store.start_date if empty
		BillType  string  `json:"bill_type"`
		Rate      float64 `json:"rate"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.GodownID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "godown_id is required")
		return
	}
	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
		return
	}

	if req.BillType == "" {
		req.BillType = "fixed"
	}
	if req.BillType != "weight" && req.BillType != "quantity" && req.BillType != "fixed" {
		utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight', 'quantity', or 'fixed'")
		return
	}
	if req.Rate < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "rate cannot be negative")
		return
	}

	godownExists, err := h.app.Models.Godown.Exists(r.Context(), req.GodownID)
	if err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: failed to verify godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
		return
	}
	if !godownExists {
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// ---- Derive or validate month_year ----
	startMonth := store.StartDate.Format("2006-01")
	currentMonth := time.Now().Format("2006-01")

	if req.MonthYear == "" {
		if currentMonth >= startMonth {
			req.MonthYear = currentMonth
		} else {
			req.MonthYear = startMonth
		}
	} else {
		if len(req.MonthYear) != 7 || req.MonthYear[4] != '-' {
			utils.ErrorJson(w, http.StatusBadRequest, "month_year must be in YYYY-MM format")
			return
		}
		if req.MonthYear < startMonth {
			utils.ErrorJson(w, http.StatusBadRequest, "month_year cannot be earlier than the store start month")
			return
		}
		if req.MonthYear > currentMonth {
			utils.ErrorJson(w, http.StatusBadRequest, "month_year cannot be in the future")
			return
		}
	}

	exists, err := h.app.Models.GodownBill.ExistsForGodownStoreMonth(r.Context(), req.GodownID, req.StoreID, req.MonthYear)
	if err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: failed to check existing bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify bill")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "godown bill already exists for this godown, store and month")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	total := billTotalFromSnapshot(req.BillType, req.Rate, store.Weight, store.Quantity)

	bill := &models.GodownBill{
		UserID:                userID,
		GodownID:              req.GodownID,
		StoreID:               req.StoreID,
		MonthYear:             req.MonthYear,
		BillType:              req.BillType,
		Rate:                  req.Rate,
		WeightAtBilling:       store.Weight,
		QuantityAtBilling:     store.Quantity,
		WeightUnitAtBilling:   &lot.WeightUnit,
		QuantityUnitAtBilling: &lot.QuantityUnit,
		TotalAmount:           total,
	}

	id, err := h.app.Models.GodownBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[GODOWN_BILLS] CreateGodownBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create godown bill")
		return
	}

	bill.ID = id

	log.Printf("[GODOWN_BILLS] CreateGodownBillHandler SUCCESS: created bill ID=%d month=%s", id, req.MonthYear)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created godown bill #" + strconv.FormatInt(id, 10) + " for " + req.MonthYear,
		EntityType:  "godown_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "godown bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

// GetGodownBillHandler - GET /api/godown-bills/{id}
func (h *Handler) GetGodownBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWN_BILLS] GetGodownBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWN_BILLS] GetGodownBillHandler called - ID: %d", id)

	bill, err := h.app.Models.GodownBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWN_BILLS] GetGodownBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown bill")
		return
	}
	if bill == nil {
		log.Printf("[GODOWN_BILLS] GetGodownBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown bill not found")
		return
	}

	responses, err := h.buildGodownBillResponses(r, []models.GodownBill{*bill})
	if err != nil {
		log.Printf("[GODOWN_BILLS] GetGodownBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[GODOWN_BILLS] GetGodownBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "godown bill details fetched", responses[0])
}

// GetAllGodownBillsHandler - GET /api/godown-bills
func (h *Handler) GetAllGodownBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GODOWN_BILLS] GetAllGodownBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	storeID := r.URL.Query().Get("store_id")
	godownID := r.URL.Query().Get("godown_id")
	monthYear := r.URL.Query().Get("month_year")

	var bills []models.GodownBill
	var err error

	switch {
	case monthYear != "":
		bills, err = h.app.Models.GodownBill.GetByMonth(r.Context(), monthYear)
	case storeID != "":
		id, parseErr := strconv.ParseInt(storeID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid store_id")
			return
		}
		bills, err = h.app.Models.GodownBill.GetByStoreID(r.Context(), id)
	case godownID != "":
		id, parseErr := strconv.ParseInt(godownID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid godown_id")
			return
		}
		bills, err = h.app.Models.GodownBill.GetByGodownID(r.Context(), id)
	default:
		bills, err = h.app.Models.GodownBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[GODOWN_BILLS] GetAllGodownBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown bills")
		return
	}

	responses, err := h.buildGodownBillResponses(r, bills)
	if err != nil {
		log.Printf("[GODOWN_BILLS] GetAllGodownBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[GODOWN_BILLS] GetAllGodownBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "godown bills fetched successfully", responses)
}

// GetCurrentMonthGodownBillsHandler - GET /api/godown-bills/current-month
func (h *Handler) GetCurrentMonthGodownBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	currentMonth := time.Now().Format("2006-01")
	ctx := r.Context()

	userID, ok := ctx.Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if err := h.app.Models.GodownBill.EnsureCurrentMonthBills(ctx, userID); err != nil {
		log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler ERROR: failed to ensure bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to prepare bills")
		return
	}

	bills, err := h.app.Models.GodownBill.GetByMonth(ctx, currentMonth)
	if err != nil {
		log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler ERROR: failed to fetch bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bills")
		return
	}

	responses, err := h.buildGodownBillResponses(r, bills)
	if err != nil {
		log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[GODOWN_BILLS] GetCurrentMonthGodownBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "current month godown bills fetched", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateGodownBillHandler - PATCH /api/godown-bills/{id}
func (h *Handler) UpdateGodownBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler called - ID: %d", id)

	type request struct {
		BillType          *string  `json:"bill_type,omitempty"`
		Rate              *float64 `json:"rate,omitempty"`
		WeightAtBilling   *float64 `json:"weight_at_billing,omitempty"`
		QuantityAtBilling *float64 `json:"quantity_at_billing,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.GodownBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown bill")
		return
	}
	if bill == nil {
		log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown bill not found")
		return
	}

	if req.BillType != nil {
		if *req.BillType != "weight" && *req.BillType != "quantity" && *req.BillType != "fixed" {
			utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight', 'quantity', or 'fixed'")
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

	// Always recompute the total from the current values.
	bill.TotalAmount = billTotalFromSnapshot(bill.BillType, bill.Rate, bill.WeightAtBilling, bill.QuantityAtBilling)

	if err := h.app.Models.GodownBill.Update(r.Context(), bill); err != nil {
		log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update godown bill")
		return
	}

	log.Printf("[GODOWN_BILLS] UpdateGodownBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated godown bill #" + strconv.FormatInt(id, 10),
		EntityType:  "godown_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.GodownBill.GetByID(r.Context(), id)
	responses, _ := h.buildGodownBillResponses(r, []models.GodownBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "godown bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteGodownBillHandler - DELETE /api/godown-bills/{id}
func (h *Handler) DeleteGodownBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler called - ID: %d", id)

	bill, err := h.app.Models.GodownBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown bill")
		return
	}
	if bill == nil {
		log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "godown", id)
	if err != nil {
		log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.GodownBill.Delete(r.Context(), id); err != nil {
		log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete godown bill")
		return
	}

	log.Printf("[GODOWN_BILLS] DeleteGodownBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted godown bill #" + strconv.FormatInt(id, 10),
		EntityType:  "godown_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "godown bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildGodownBillResponses(r *http.Request, bills []models.GodownBill) ([]GodownBillResponse, error) {
	if len(bills) == 0 {
		return []GodownBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "godown", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]GodownBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.TotalAmount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, GodownBillResponse{
			GodownBill:       b,
			TotalPaid:        t.TotalPaid,
			TotalPaidThrough: t.TotalPaidThrough,
			Remaining:        remaining,
		})
	}

	return out, nil
}
