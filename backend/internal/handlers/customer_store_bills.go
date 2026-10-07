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

type CustomerStoreBillResponse struct {
	models.CustomerStoreBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

// CreateCustomerStoreBillHandler - POST /api/customer-store-bills
func (h *Handler) CreateCustomerStoreBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID int64   `json:"customer_id"`
		StoreID    int64   `json:"store_id"`
		MonthYear  string  `json:"month_year"` // optional — derived from store.start_date if empty
		BillType   string  `json:"bill_type"`
		Rate       float64 `json:"rate"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_id is required")
		return
	}
	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
		return
	}

	if req.BillType == "" {
		req.BillType = "quantity"
	}
	if req.BillType != "weight" && req.BillType != "quantity" {
		utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight' or 'quantity'")
		return
	}
	if req.Rate < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "rate cannot be negative")
		return
	}

	customerExists, err := h.app.Models.Customer.Exists(r.Context(), req.CustomerID)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !customerExists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: failed to fetch lot - %v", err)
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

	exists, err := h.app.Models.CustomerStoreBill.ExistsForCustomerStoreMonth(r.Context(), req.CustomerID, req.StoreID, req.MonthYear)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: failed to check existing bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify bill")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "customer store bill already exists for this customer, store and month")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	total := billTotalFromSnapshot(req.BillType, req.Rate, store.Weight, store.Quantity)

	bill := &models.CustomerStoreBill{
		UserID:                userID,
		CustomerID:            req.CustomerID,
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

	id, err := h.app.Models.CustomerStoreBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer store bill")
		return
	}

	bill.ID = id

	log.Printf("[CUSTOMER_STORE_BILLS] CreateCustomerStoreBillHandler SUCCESS: created bill ID=%d month=%s", id, req.MonthYear)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer store bill #" + strconv.FormatInt(id, 10) + " for " + req.MonthYear,
		EntityType:  "customer_store_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer store bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

// GetCustomerStoreBillHandler - GET /api/customer-store-bills/{id}
func (h *Handler) GetCustomerStoreBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerStoreBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer store bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer store bill not found")
		return
	}

	responses, err := h.buildCustomerStoreBillResponses(r, []models.CustomerStoreBill{*bill})
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] GetCustomerStoreBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "customer store bill details fetched", responses[0])
}

// GetAllCustomerStoreBillsHandler - GET /api/customer-store-bills
func (h *Handler) GetAllCustomerStoreBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_STORE_BILLS] GetAllCustomerStoreBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	storeID := r.URL.Query().Get("store_id")
	customerID := r.URL.Query().Get("customer_id")
	monthYear := r.URL.Query().Get("month_year")

	var bills []models.CustomerStoreBill
	var err error

	switch {
	case monthYear != "":
		bills, err = h.app.Models.CustomerStoreBill.GetByMonth(r.Context(), monthYear)
	case storeID != "":
		id, parseErr := strconv.ParseInt(storeID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid store_id")
			return
		}
		bills, err = h.app.Models.CustomerStoreBill.GetByStoreID(r.Context(), id)
	case customerID != "":
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		bills, err = h.app.Models.CustomerStoreBill.GetByCustomerID(r.Context(), id)
	default:
		bills, err = h.app.Models.CustomerStoreBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetAllCustomerStoreBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer store bills")
		return
	}

	responses, err := h.buildCustomerStoreBillResponses(r, bills)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetAllCustomerStoreBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] GetAllCustomerStoreBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "customer store bills fetched successfully", responses)
}

// GetCurrentMonthCustomerStoreBillsHandler - GET /api/customer-store-bills/current-month
func (h *Handler) GetCurrentMonthCustomerStoreBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	currentMonth := time.Now().Format("2006-01")
	ctx := r.Context()

	userID, ok := ctx.Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if err := h.app.Models.CustomerStoreBill.EnsureCurrentMonthBills(ctx, userID); err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler ERROR: failed to ensure bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to prepare bills")
		return
	}

	bills, err := h.app.Models.CustomerStoreBill.GetByMonth(ctx, currentMonth)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler ERROR: failed to fetch bills - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bills")
		return
	}

	responses, err := h.buildCustomerStoreBillResponses(r, bills)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] GetCurrentMonthCustomerStoreBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "current month customer store bills fetched", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateCustomerStoreBillHandler - PATCH /api/customer-store-bills/{id}
func (h *Handler) UpdateCustomerStoreBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler called - ID: %d", id)

	type request struct {
		BillType          *string  `json:"bill_type,omitempty"`
		Rate              *float64 `json:"rate,omitempty"`
		WeightAtBilling   *float64 `json:"weight_at_billing,omitempty"`
		QuantityAtBilling *float64 `json:"quantity_at_billing,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.CustomerStoreBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer store bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer store bill not found")
		return
	}

	if req.BillType != nil {
		if *req.BillType != "weight" && *req.BillType != "quantity" {
			utils.ErrorJson(w, http.StatusBadRequest, "bill_type must be 'weight' or 'quantity'")
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

	if err := h.app.Models.CustomerStoreBill.Update(r.Context(), bill); err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer store bill")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] UpdateCustomerStoreBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer store bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_store_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.CustomerStoreBill.GetByID(r.Context(), id)
	responses, _ := h.buildCustomerStoreBillResponses(r, []models.CustomerStoreBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "customer store bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteCustomerStoreBillHandler - DELETE /api/customer-store-bills/{id}
func (h *Handler) DeleteCustomerStoreBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerStoreBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer store bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer store bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "customer_store", id)
	if err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.CustomerStoreBill.Delete(r.Context(), id); err != nil {
		log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer store bill")
		return
	}

	log.Printf("[CUSTOMER_STORE_BILLS] DeleteCustomerStoreBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer store bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_store_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer store bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildCustomerStoreBillResponses(r *http.Request, bills []models.CustomerStoreBill) ([]CustomerStoreBillResponse, error) {
	if len(bills) == 0 {
		return []CustomerStoreBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "customer_store", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]CustomerStoreBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.TotalAmount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, CustomerStoreBillResponse{
			CustomerStoreBill: b,
			TotalPaid:         t.TotalPaid,
			TotalPaidThrough:  t.TotalPaidThrough,
			Remaining:         remaining,
		})
	}

	return out, nil
}
