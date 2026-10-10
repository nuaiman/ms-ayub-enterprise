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

type CustomerUnloadBillResponse struct {
	models.CustomerUnloadBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

// CreateCustomerUnloadBillHandler - POST /api/customer-unload-bills
func (h *Handler) CreateCustomerUnloadBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID int64   `json:"customer_id"`
		StoreID    int64   `json:"store_id"`
		BillType   string  `json:"bill_type"`
		Rate       float64 `json:"rate"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: invalid request body - %v", err)
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
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !customerExists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	total := billTotalFromSnapshot(req.BillType, req.Rate, store.Weight, store.Quantity)

	bill := &models.CustomerUnloadBill{
		UserID:                userID,
		CustomerID:            req.CustomerID,
		StoreID:               req.StoreID,
		BillType:              req.BillType,
		Rate:                  req.Rate,
		WeightAtBilling:       store.Weight,
		QuantityAtBilling:     store.Quantity,
		WeightUnitAtBilling:   &lot.WeightUnit,
		QuantityUnitAtBilling: &lot.QuantityUnit,
		TotalAmount:           total,
	}

	id, err := h.app.Models.CustomerUnloadBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer unload bill")
		return
	}

	bill.ID = id

	log.Printf("[CUSTOMER_UNLOAD_BILLS] CreateCustomerUnloadBillHandler SUCCESS: created bill ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer unload bill #" + strconv.FormatInt(id, 10) + " for store #" + strconv.FormatInt(req.StoreID, 10),
		EntityType:  "customer_unload_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer unload bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

// GetCustomerUnloadBillHandler - GET /api/customer-unload-bills/{id}
func (h *Handler) GetCustomerUnloadBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerUnloadBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer unload bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer unload bill not found")
		return
	}

	responses, err := h.buildCustomerUnloadBillResponses(r, []models.CustomerUnloadBill{*bill})
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] GetCustomerUnloadBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "customer unload bill details fetched", responses[0])
}

// GetAllCustomerUnloadBillsHandler - GET /api/customer-unload-bills
func (h *Handler) GetAllCustomerUnloadBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_UNLOAD_BILLS] GetAllCustomerUnloadBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	customerID := r.URL.Query().Get("customer_id")
	storeID := r.URL.Query().Get("store_id")

	var bills []models.CustomerUnloadBill
	var err error

	switch {
	case customerID != "":
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		bills, err = h.app.Models.CustomerUnloadBill.GetByCustomerID(r.Context(), id)
	case storeID != "":
		id, parseErr := strconv.ParseInt(storeID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid store_id")
			return
		}
		bills, err = h.app.Models.CustomerUnloadBill.GetByStoreID(r.Context(), id)
	default:
		bills, err = h.app.Models.CustomerUnloadBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetAllCustomerUnloadBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer unload bills")
		return
	}

	responses, err := h.buildCustomerUnloadBillResponses(r, bills)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] GetAllCustomerUnloadBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] GetAllCustomerUnloadBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "customer unload bills fetched successfully", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateCustomerUnloadBillHandler - PATCH /api/customer-unload-bills/{id}
func (h *Handler) UpdateCustomerUnloadBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler called - ID: %d", id)

	type request struct {
		BillType          *string  `json:"bill_type,omitempty"`
		Rate              *float64 `json:"rate,omitempty"`
		WeightAtBilling   *float64 `json:"weight_at_billing,omitempty"`
		QuantityAtBilling *float64 `json:"quantity_at_billing,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.CustomerUnloadBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer unload bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer unload bill not found")
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

	if err := h.app.Models.CustomerUnloadBill.Update(r.Context(), bill); err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer unload bill")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] UpdateCustomerUnloadBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer unload bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_unload_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.CustomerUnloadBill.GetByID(r.Context(), id)
	responses, _ := h.buildCustomerUnloadBillResponses(r, []models.CustomerUnloadBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "customer unload bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteCustomerUnloadBillHandler - DELETE /api/customer-unload-bills/{id}
func (h *Handler) DeleteCustomerUnloadBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerUnloadBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer unload bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer unload bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "customer_unload", id)
	if err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.CustomerUnloadBill.Delete(r.Context(), id); err != nil {
		log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer unload bill")
		return
	}

	log.Printf("[CUSTOMER_UNLOAD_BILLS] DeleteCustomerUnloadBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer unload bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_unload_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer unload bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildCustomerUnloadBillResponses(r *http.Request, bills []models.CustomerUnloadBill) ([]CustomerUnloadBillResponse, error) {
	if len(bills) == 0 {
		return []CustomerUnloadBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "customer_unload", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]CustomerUnloadBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.TotalAmount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, CustomerUnloadBillResponse{
			CustomerUnloadBill: b,
			TotalPaid:          t.TotalPaid,
			TotalPaidThrough:   t.TotalPaidThrough,
			Remaining:          remaining,
		})
	}

	return out, nil
}
