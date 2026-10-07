package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"time"
)

// =============================================================================
// BILL PAYMENTS
// =============================================================================

func billTotalFromSnapshot(billType string, rate, weightAtBilling, quantityAtBilling float64) float64 {
	switch billType {
	case "weight":
		return rate * weightAtBilling
	case "quantity":
		return rate * quantityAtBilling
	case "job", "fixed":
		return rate
	default:
		return 0
	}
}

// =============================================================================
// CREATE PAYMENT
// =============================================================================

func (h *Handler) CreateGodownBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	h.createBillPayment(w, r, "godown")
}

func (h *Handler) CreateMajhiBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	h.createBillPayment(w, r, "majhi")
}

func (h *Handler) CreateCustomerStoreBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	h.createBillPayment(w, r, "customer_store")
}

func (h *Handler) CreateCustomerDeliveryBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	h.createBillPayment(w, r, "customer_delivery")
}

func (h *Handler) CreateCustomerAdditionalBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	h.createBillPayment(w, r, "customer_additional")
}

func (h *Handler) createBillPayment(w http.ResponseWriter, r *http.Request, billType string) {
	log.Printf("[BILL_PAYMENTS] CreateBillPayment called - type=%s method=%s path=%s", billType, r.Method, r.URL.Path)

	type request struct {
		Amount          float64    `json:"amount"`
		PaymentDate     *time.Time `json:"payment_date,omitempty"`
		PaymentMethod   *string    `json:"payment_method,omitempty"`
		ReferenceNumber *string    `json:"reference_number,omitempty"`
		Notes           *string    `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Amount <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	if req.PaymentMethod != nil && *req.PaymentMethod != "" {
		valid := map[string]bool{
			"cash": true, "bank_transfer": true, "check": true, "mobile_banking": true,
		}
		if !valid[*req.PaymentMethod] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid payment_method")
			return
		}
	}

	billID, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	total, found, err := h.lookupBillTotal(r, billType, billID)
	if err != nil {
		log.Printf("[BILL_PAYMENTS] lookup failed: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill")
		return
	}
	if !found {
		utils.ErrorJson(w, http.StatusNotFound, "bill not found")
		return
	}

	alreadyPaid, err := h.app.Models.BillPayment.GetTotalPaid(r.Context(), billType, billID)
	if err != nil {
		log.Printf("[BILL_PAYMENTS] failed to sum payments: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify payments")
		return
	}

	if alreadyPaid+req.Amount > total+0.0001 {
		utils.ErrorJson(w, http.StatusBadRequest, "payment exceeds remaining balance")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	paymentDate := time.Now()
	if req.PaymentDate != nil {
		paymentDate = *req.PaymentDate
	}

	payment := &models.BillPayment{
		UserID:          userID,
		BillType:        billType,
		BillID:          billID,
		Amount:          req.Amount,
		PaymentDate:     paymentDate,
		PaymentMethod:   req.PaymentMethod,
		ReferenceNumber: req.ReferenceNumber,
		Notes:           req.Notes,
	}

	id, err := h.app.Models.BillPayment.Insert(r.Context(), payment)
	if err != nil {
		log.Printf("[BILL_PAYMENTS] insert failed: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to record payment")
		return
	}
	payment.ID = id

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Recorded payment for " + billType + " bill #" + utils.FormatInt64(billID),
		EntityType:  "bill_payments",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	newTotalPaid, _ := h.app.Models.BillPayment.GetTotalPaid(r.Context(), billType, billID)
	newPaidThrough, _ := h.app.Models.BillPayment.GetTotalPaidThrough(r.Context(), billType, billID)

	remaining := total - newTotalPaid
	if remaining < 0 {
		remaining = 0
	}

	utils.SuccessJson(w, http.StatusCreated, "payment recorded successfully", map[string]any{
		"payment":            payment,
		"bill_total":         total,
		"total_paid":         newTotalPaid,
		"remaining":          remaining,
		"total_paid_through": newPaidThrough,
	})
}

// =============================================================================
// LIST PAYMENTS
// =============================================================================

func (h *Handler) GetGodownBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	h.getBillPayments(w, r, "godown")
}

func (h *Handler) GetMajhiBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	h.getBillPayments(w, r, "majhi")
}

func (h *Handler) GetCustomerStoreBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	h.getBillPayments(w, r, "customer_store")
}

func (h *Handler) GetCustomerDeliveryBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	h.getBillPayments(w, r, "customer_delivery")
}

func (h *Handler) GetCustomerAdditionalBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	h.getBillPayments(w, r, "customer_additional")
}

func (h *Handler) getBillPayments(w http.ResponseWriter, r *http.Request, billType string) {
	billID, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	_, found, err := h.lookupBillTotal(r, billType, billID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch bill")
		return
	}
	if !found {
		utils.ErrorJson(w, http.StatusNotFound, "bill not found")
		return
	}

	payments, err := h.app.Models.BillPayment.GetByBill(r.Context(), billType, billID)
	if err != nil {
		log.Printf("[BILL_PAYMENTS] list failed: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch payments")
		return
	}
	if payments == nil {
		payments = []models.BillPayment{}
	}

	totalPaid, _ := h.app.Models.BillPayment.GetTotalPaid(r.Context(), billType, billID)
	paidThrough, _ := h.app.Models.BillPayment.GetTotalPaidThrough(r.Context(), billType, billID)

	utils.SuccessJson(w, http.StatusOK, "payments fetched successfully", map[string]any{
		"payments":           payments,
		"total_paid":         totalPaid,
		"total_paid_through": paidThrough,
	})
}

// =============================================================================
// DELETE PAYMENT
// =============================================================================

func (h *Handler) DeleteBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	payment, err := h.app.Models.BillPayment.GetByID(r.Context(), id)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch payment")
		return
	}
	if payment == nil {
		utils.ErrorJson(w, http.StatusNotFound, "payment not found")
		return
	}

	if err := h.app.Models.BillPayment.Delete(r.Context(), id); err != nil {
		log.Printf("[BILL_PAYMENTS] delete failed: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete payment")
		return
	}

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted bill payment #" + utils.FormatInt64(id),
		EntityType:  "bill_payments",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "payment deleted successfully", nil)
}

// =============================================================================
// INTERNAL LOOKUPS
// =============================================================================

func (h *Handler) lookupBillTotal(r *http.Request, billType string, id int64) (float64, bool, error) {
	switch billType {
	case "godown":
		b, err := h.app.Models.GodownBill.GetByID(r.Context(), id)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		return b.TotalAmount, true, nil

	case "majhi":
		b, err := h.app.Models.MajhiBill.GetByID(r.Context(), id)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		return b.TotalAmount, true, nil

	case "customer_store":
		b, err := h.app.Models.CustomerStoreBill.GetByID(r.Context(), id)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		return b.TotalAmount, true, nil

	case "customer_delivery":
		b, err := h.app.Models.CustomerDeliveryBill.GetByID(r.Context(), id)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		return b.TotalAmount, true, nil

	case "customer_additional":
		b, err := h.app.Models.CustomerAdditionalBill.GetByID(r.Context(), id)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		return b.Amount, true, nil
	}
	return 0, false, nil
}
