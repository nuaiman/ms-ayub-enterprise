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

type CustomerDeliveryBillResponse struct {
	models.CustomerDeliveryBill
	TotalPaid        float64    `json:"total_paid"`
	TotalPaidThrough *time.Time `json:"total_paid_through,omitempty"`
	Remaining        float64    `json:"remaining"`
}

// =============================================================================
// CREATE
// =============================================================================

func (h *Handler) CreateCustomerDeliveryBillHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID     int64   `json:"customer_id"`
		DeliveryItemID int64   `json:"delivery_item_id"`
		BillType       string  `json:"bill_type"`
		Rate           float64 `json:"rate"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_id is required")
		return
	}
	if req.DeliveryItemID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "delivery_item_id is required")
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
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !customerExists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	deliveryItem, err := h.app.Models.Delivery.GetItemByID(r.Context(), req.DeliveryItemID)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify delivery item")
		return
	}
	if deliveryItem == nil {
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	exists, err := h.app.Models.CustomerDeliveryBill.ExistsForDeliveryItem(r.Context(), req.DeliveryItemID)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to check existing bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify bill")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "customer delivery bill already exists for this delivery item")
		return
	}

	store, err := h.app.Models.Store.GetByID(r.Context(), deliveryItem.StoreID)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), store.LotID)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	total := billTotalFromSnapshot(req.BillType, req.Rate, deliveryItem.Weight, deliveryItem.Quantity)

	bill := &models.CustomerDeliveryBill{
		UserID:                userID,
		CustomerID:            req.CustomerID,
		DeliveryItemID:        req.DeliveryItemID,
		BillType:              req.BillType,
		Rate:                  req.Rate,
		WeightAtBilling:       deliveryItem.Weight,
		QuantityAtBilling:     deliveryItem.Quantity,
		WeightUnitAtBilling:   &lot.WeightUnit,
		QuantityUnitAtBilling: &lot.QuantityUnit,
		TotalAmount:           total,
	}

	id, err := h.app.Models.CustomerDeliveryBill.Insert(r.Context(), bill)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler ERROR: failed to insert bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer delivery bill")
		return
	}

	bill.ID = id

	log.Printf("[CUSTOMER_DELIVERY_BILLS] CreateCustomerDeliveryBillHandler SUCCESS: created bill ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer delivery bill #" + strconv.FormatInt(id, 10) + " for delivery item #" + strconv.FormatInt(req.DeliveryItemID, 10),
		EntityType:  "customer_delivery_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer delivery bill created successfully", bill)
}

// =============================================================================
// READ
// =============================================================================

func (h *Handler) GetCustomerDeliveryBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerDeliveryBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer delivery bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer delivery bill not found")
		return
	}

	responses, err := h.buildCustomerDeliveryBillResponses(r, []models.CustomerDeliveryBill{*bill})
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler ERROR: failed to build response - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] GetCustomerDeliveryBillHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "customer delivery bill details fetched", responses[0])
}

func (h *Handler) GetAllCustomerDeliveryBillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMER_DELIVERY_BILLS] GetAllCustomerDeliveryBillsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	customerID := r.URL.Query().Get("customer_id")
	deliveryItemID := r.URL.Query().Get("delivery_item_id")

	var bills []models.CustomerDeliveryBill
	var err error

	switch {
	case deliveryItemID != "":
		id, parseErr := strconv.ParseInt(deliveryItemID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid delivery_item_id")
			return
		}
		one, e := h.app.Models.CustomerDeliveryBill.GetByDeliveryItemID(r.Context(), id)
		if e != nil {
			err = e
		} else if one != nil {
			bills = []models.CustomerDeliveryBill{*one}
		}
	case customerID != "":
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		bills, err = h.app.Models.CustomerDeliveryBill.GetByCustomerID(r.Context(), id)
	default:
		bills, err = h.app.Models.CustomerDeliveryBill.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetAllCustomerDeliveryBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer delivery bills")
		return
	}

	responses, err := h.buildCustomerDeliveryBillResponses(r, bills)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] GetAllCustomerDeliveryBillsHandler ERROR: failed to build responses - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill totals")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] GetAllCustomerDeliveryBillsHandler SUCCESS: fetched %d bills", len(responses))
	utils.SuccessJson(w, http.StatusOK, "customer delivery bills fetched successfully", responses)
}

// =============================================================================
// UPDATE
// =============================================================================

func (h *Handler) UpdateCustomerDeliveryBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler called - ID: %d", id)

	type request struct {
		BillType          *string  `json:"bill_type,omitempty"`
		Rate              *float64 `json:"rate,omitempty"`
		WeightAtBilling   *float64 `json:"weight_at_billing,omitempty"`
		QuantityAtBilling *float64 `json:"quantity_at_billing,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bill, err := h.app.Models.CustomerDeliveryBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer delivery bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer delivery bill not found")
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

	if err := h.app.Models.CustomerDeliveryBill.Update(r.Context(), bill); err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler ERROR: failed to update bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer delivery bill")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] UpdateCustomerDeliveryBillHandler SUCCESS: updated bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer delivery bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_delivery_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.CustomerDeliveryBill.GetByID(r.Context(), id)
	responses, _ := h.buildCustomerDeliveryBillResponses(r, []models.CustomerDeliveryBill{*updated})
	utils.SuccessJson(w, http.StatusOK, "customer delivery bill updated successfully", responses[0])
}

// =============================================================================
// DELETE
// =============================================================================

func (h *Handler) DeleteCustomerDeliveryBillHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler called - ID: %d", id)

	bill, err := h.app.Models.CustomerDeliveryBill.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler ERROR: failed to fetch bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer delivery bill")
		return
	}
	if bill == nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer delivery bill not found")
		return
	}

	// Block deletion if this bill is attached to any invoice.
	invoiced, err := h.app.Models.InvoiceItem.ExistsForBill(r.Context(), "customer_delivery", id)
	if err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler ERROR: failed to check invoice linkage - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify invoice linkage")
		return
	}
	if invoiced {
		utils.ErrorJson(w, http.StatusConflict, "cannot delete bill: it is attached to an invoice")
		return
	}

	if err := h.app.Models.CustomerDeliveryBill.Delete(r.Context(), id); err != nil {
		log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler ERROR: failed to delete bill - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer delivery bill")
		return
	}

	log.Printf("[CUSTOMER_DELIVERY_BILLS] DeleteCustomerDeliveryBillHandler SUCCESS: deleted bill ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer delivery bill #" + strconv.FormatInt(id, 10),
		EntityType:  "customer_delivery_bills",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer delivery bill deleted successfully", nil)
}

// =============================================================================
// RESPONSE BUILDERS
// =============================================================================

func (h *Handler) buildCustomerDeliveryBillResponses(r *http.Request, bills []models.CustomerDeliveryBill) ([]CustomerDeliveryBillResponse, error) {
	if len(bills) == 0 {
		return []CustomerDeliveryBillResponse{}, nil
	}

	billIDs := make([]int64, 0, len(bills))
	for _, b := range bills {
		billIDs = append(billIDs, b.ID)
	}

	totals, err := h.app.Models.BillPayment.GetTotalsForBills(r.Context(), "customer_delivery", billIDs)
	if err != nil {
		return nil, err
	}

	out := make([]CustomerDeliveryBillResponse, 0, len(bills))
	for _, b := range bills {
		t := totals[b.ID]

		remaining := b.TotalAmount - t.TotalPaid
		if remaining < 0 {
			remaining = 0
		}

		out = append(out, CustomerDeliveryBillResponse{
			CustomerDeliveryBill: b,
			TotalPaid:            t.TotalPaid,
			TotalPaidThrough:     t.TotalPaidThrough,
			Remaining:            remaining,
		})
	}

	return out, nil
}
