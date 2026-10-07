package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// =============================================================================
// SHARED TYPES
// =============================================================================

type UnbilledBill struct {
	BillType  string  `json:"bill_type"`
	BillID    int64   `json:"bill_id"`
	Label     string  `json:"label"`
	Total     float64 `json:"total"`
	Paid      float64 `json:"paid"`
	Remaining float64 `json:"remaining"`
	Date      string  `json:"date"`
}

type discountInput struct {
	Type   string  `json:"type"`
	Value  float64 `json:"value"`
	Reason *string `json:"reason,omitempty"`
}

type billRef struct {
	BillType string `json:"bill_type"`
	BillID   int64  `json:"bill_id"`
}

type InvoiceDetailResponse struct {
	models.Invoice
	Items     []models.InvoiceItem `json:"items"`
	Discounts []models.Discount    `json:"discounts"`
}

var (
	errDiscountType          = errors.New("discount type must be 'flat' or 'percent'")
	errDiscountValue         = errors.New("discount value cannot be negative")
	errDiscountPercentTooBig = errors.New("percent discount cannot exceed 100")
	errBillNotAllowed        = errors.New("bill type is not allowed for this entity")
)

// billTypeAllowedForEntity reports whether the given bill_type can be
// attached to an invoice for the given entity_type.
func billTypeAllowedForEntity(entityType, billType string) bool {
	switch entityType {
	case "customer":
		return billType == "customer_store" ||
			billType == "customer_delivery" ||
			billType == "customer_additional"
	case "majhi":
		return billType == "majhi"
	case "godown":
		return billType == "godown"
	case "broker":
		// Brokers don't have bills yet. Kept as a valid entity for the
		// API surface, but no bill types are allowed.
		return false
	default:
		return false
	}
}

// =============================================================================
// CREATE
// =============================================================================

// CreateInvoiceHandler - POST /api/invoices
//
// Body:
//
//	{
//	  "entity_type": "customer" | "majhi" | "godown" | "broker",
//	  "entity_id": 12,
//	  "bills": [ { "bill_type": "customer_store", "bill_id": 3 }, ... ],
//	  "discounts": [ { "type": "flat", "value": 200, "reason": "..." }, ... ],
//	  "notes": "optional"
//	}
//
// The handler verifies each bill exists, is unpaid, and belongs to the
// entity, snapshots the remaining amount into invoice_items, and stores the
// invoice. The individual bills are NOT mutated.
func (h *Handler) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INVOICES] CreateInvoiceHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		EntityType string          `json:"entity_type"`
		EntityID   int64           `json:"entity_id"`
		Bills      []billRef       `json:"bills"`
		Discounts  []discountInput `json:"discounts,omitempty"`
		Notes      *string         `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[INVOICES] CreateInvoiceHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validEntityTypes := map[string]bool{
		"customer": true, "majhi": true, "godown": true, "broker": true,
	}
	if !validEntityTypes[req.EntityType] {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_type must be 'customer', 'majhi', 'godown', or 'broker'")
		return
	}
	if req.EntityID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_id is required")
		return
	}
	if len(req.Bills) == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "at least one bill must be selected")
		return
	}

	if !h.entityExists(r, req.EntityType, req.EntityID) {
		utils.ErrorJson(w, http.StatusNotFound, "entity not found")
		return
	}

	// ---- Verify each selected bill and compute subtotal. ----
	subtotal := 0.0
	invoiceItems := make([]*models.InvoiceItem, 0, len(req.Bills))

	for _, ref := range req.Bills {
		if !billTypeAllowedForEntity(req.EntityType, ref.BillType) {
			utils.ErrorJson(w, http.StatusBadRequest, errBillNotAllowed.Error())
			return
		}
		if ref.BillID == 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "bill_id is required on every bill")
			return
		}

		remaining, ok, err := h.billRemaining(r, ref.BillType, ref.BillID)
		if err != nil {
			log.Printf("[INVOICES] CreateInvoiceHandler ERROR: failed to lookup bill %s#%d - %v", ref.BillType, ref.BillID, err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to lookup bill")
			return
		}
		if !ok {
			utils.ErrorJson(w, http.StatusNotFound, "bill not found: "+ref.BillType+"#"+strconv.FormatInt(ref.BillID, 10))
			return
		}
		if remaining <= 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "bill "+ref.BillType+"#"+strconv.FormatInt(ref.BillID, 10)+" is already fully paid")
			return
		}

		// Verify the bill belongs to the requested entity.
		if !h.billBelongsToEntity(r, ref.BillType, ref.BillID, req.EntityType, req.EntityID) {
			utils.ErrorJson(w, http.StatusBadRequest, "bill "+ref.BillType+"#"+strconv.FormatInt(ref.BillID, 10)+" does not belong to the selected entity")
			return
		}

		subtotal += remaining
		invoiceItems = append(invoiceItems, &models.InvoiceItem{
			BillType: ref.BillType,
			BillID:   ref.BillID,
			Amount:   remaining,
		})
	}

	discountTotal, discountRows, err := buildDiscountsFromInputs(req.Discounts, subtotal)
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	total := subtotal - discountTotal
	if total < 0 {
		total = 0
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	invoice := &models.Invoice{
		UserID:         userID,
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		Subtotal:       subtotal,
		DiscountAmount: discountTotal,
		Total:          total,
		Notes:          req.Notes,
	}

	invoiceID, err := h.app.Models.Invoice.Insert(r.Context(), invoice)
	if err != nil {
		log.Printf("[INVOICES] CreateInvoiceHandler ERROR: failed to insert invoice - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create invoice")
		return
	}
	invoice.ID = invoiceID

	// Attach invoice_id to all items and insert in one transaction.
	for _, it := range invoiceItems {
		it.InvoiceID = invoiceID
	}
	if err := h.app.Models.InvoiceItem.InsertMany(r.Context(), invoiceItems); err != nil {
		log.Printf("[INVOICES] CreateInvoiceHandler ERROR: failed to insert invoice items - %v", err)
		// Roll back the invoice row.
		_ = h.app.Models.Invoice.Delete(r.Context(), invoiceID)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create invoice items")
		return
	}

	for _, d := range discountRows {
		d.InvoiceID = invoiceID
		if _, err := h.app.Models.Discount.Insert(r.Context(), d); err != nil {
			log.Printf("[INVOICES] CreateInvoiceHandler ERROR: failed to insert discount - %v", err)
		}
	}

	log.Printf("[INVOICES] CreateInvoiceHandler SUCCESS: invoice ID=%d entity=%s:%d bills=%d subtotal=%.2f discount=%.2f total=%.2f",
		invoiceID, req.EntityType, req.EntityID, len(req.Bills), subtotal, discountTotal, total)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "create",
		Description: "Created invoice #" + strconv.FormatInt(invoiceID, 10) +
			" for " + req.EntityType + " #" + strconv.FormatInt(req.EntityID, 10),
		EntityType: "invoices",
		EntityID:   invoiceID,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	detail, _ := h.buildInvoiceDetail(r, invoiceID)
	utils.SuccessJson(w, http.StatusCreated, "invoice created successfully", detail)
}

// =============================================================================
// READ
// =============================================================================

// GetInvoiceHandler - GET /api/invoices/{id}
func (h *Handler) GetInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	detail, err := h.buildInvoiceDetail(r, id)
	if err != nil {
		log.Printf("[INVOICES] GetInvoiceHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch invoice")
		return
	}
	if detail == nil {
		utils.ErrorJson(w, http.StatusNotFound, "invoice not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "invoice fetched successfully", detail)
}

// GetAllInvoicesHandler - GET /api/invoices
func (h *Handler) GetAllInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("entity_type")
	entityIDParam := r.URL.Query().Get("entity_id")

	var (
		invoices []models.Invoice
		err      error
	)

	if entityType != "" && entityIDParam != "" {
		entityID, parseErr := strconv.ParseInt(entityIDParam, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid entity_id")
			return
		}
		invoices, err = h.app.Models.Invoice.GetByEntity(r.Context(), entityType, entityID)
	} else {
		invoices, err = h.app.Models.Invoice.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[INVOICES] GetAllInvoicesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch invoices")
		return
	}

	if invoices == nil {
		invoices = []models.Invoice{}
	}

	utils.SuccessJson(w, http.StatusOK, "invoices fetched successfully", invoices)
}

// GetUnbilledBillsHandler - GET /api/invoices/unbilled?entity_type=customer&entity_id=3
//
// Returns unpaid bills for the given entity so the frontend can let the
// user pick which ones to include.
func (h *Handler) GetUnbilledBillsHandler(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("entity_type")
	entityIDParam := r.URL.Query().Get("entity_id")

	if entityType == "" || entityIDParam == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "entity_type and entity_id are required")
		return
	}

	validEntityTypes := map[string]bool{
		"customer": true, "majhi": true, "godown": true, "broker": true,
	}
	if !validEntityTypes[entityType] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid entity_type")
		return
	}

	entityID, err := strconv.ParseInt(entityIDParam, 10, 64)
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid entity_id")
		return
	}

	if !h.entityExists(r, entityType, entityID) {
		utils.ErrorJson(w, http.StatusNotFound, "entity not found")
		return
	}

	candidates, err := h.collectUnpaidBills(r, entityType, entityID)
	if err != nil {
		log.Printf("[INVOICES] GetUnbilledBillsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch unpaid bills")
		return
	}

	if candidates == nil {
		candidates = []UnbilledBill{}
	}

	utils.SuccessJson(w, http.StatusOK, "unbilled bills fetched successfully", candidates)
}

// =============================================================================
// UPDATE
// =============================================================================

// UpdateInvoiceHandler - PATCH /api/invoices/{id}
//
// Body:
//
//	{
//	  "discounts": [ { "type": "flat", "value": 100, "reason": "..." }, ... ],
//	  "notes": "..."
//	}
//
// Only discounts and notes are editable. Items and subtotal are frozen at
// creation time. If "discounts" is present (even as []), all existing
// discounts are replaced.
func (h *Handler) UpdateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	type request struct {
		Discounts *[]discountInput `json:"discounts,omitempty"`
		Notes     *string          `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[INVOICES] UpdateInvoiceHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	invoice, err := h.app.Models.Invoice.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[INVOICES] UpdateInvoiceHandler ERROR: failed to fetch invoice - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch invoice")
		return
	}
	if invoice == nil {
		utils.ErrorJson(w, http.StatusNotFound, "invoice not found")
		return
	}

	if req.Notes != nil {
		invoice.Notes = req.Notes
	}

	if req.Discounts != nil {
		discountTotal, discountRows, derr := buildDiscountsFromInputs(*req.Discounts, invoice.Subtotal)
		if derr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, derr.Error())
			return
		}

		if err := h.app.Models.Discount.DeleteByInvoiceID(r.Context(), invoice.ID); err != nil {
			log.Printf("[INVOICES] UpdateInvoiceHandler ERROR: failed to clear discounts - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update discounts")
			return
		}

		for _, d := range discountRows {
			d.InvoiceID = invoice.ID
			if _, err := h.app.Models.Discount.Insert(r.Context(), d); err != nil {
				log.Printf("[INVOICES] UpdateInvoiceHandler ERROR: failed to insert discount - %v", err)
			}
		}
		invoice.DiscountAmount = discountTotal
	}

	total := invoice.Subtotal - invoice.DiscountAmount
	if total < 0 {
		total = 0
	}
	invoice.Total = total

	if err := h.app.Models.Invoice.Update(r.Context(), invoice); err != nil {
		log.Printf("[INVOICES] UpdateInvoiceHandler ERROR: failed to update invoice - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update invoice")
		return
	}

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated invoice #" + strconv.FormatInt(id, 10),
		EntityType:  "invoices",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	detail, _ := h.buildInvoiceDetail(r, id)
	utils.SuccessJson(w, http.StatusOK, "invoice updated successfully", detail)
}

// =============================================================================
// DELETE
// =============================================================================

// DeleteInvoiceHandler - DELETE /api/invoices/{id}
func (h *Handler) DeleteInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		return
	}

	invoice, err := h.app.Models.Invoice.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[INVOICES] DeleteInvoiceHandler ERROR: failed to fetch invoice - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch invoice")
		return
	}
	if invoice == nil {
		utils.ErrorJson(w, http.StatusNotFound, "invoice not found")
		return
	}

	if err := h.app.Models.Invoice.Delete(r.Context(), id); err != nil {
		log.Printf("[INVOICES] DeleteInvoiceHandler ERROR: failed to delete invoice - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete invoice")
		return
	}

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted invoice #" + strconv.FormatInt(id, 10),
		EntityType:  "invoices",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "invoice deleted successfully", nil)
}

// =============================================================================
// INTERNAL HELPERS
// =============================================================================

func (h *Handler) entityExists(r *http.Request, entityType string, entityID int64) bool {
	ctx := r.Context()
	switch entityType {
	case "customer":
		ok, err := h.app.Models.Customer.Exists(ctx, entityID)
		return err == nil && ok
	case "majhi":
		ok, err := h.app.Models.Majhi.Exists(ctx, entityID)
		return err == nil && ok
	case "godown":
		ok, err := h.app.Models.Godown.Exists(ctx, entityID)
		return err == nil && ok
	case "broker":
		ok, err := h.app.Models.Broker.Exists(ctx, entityID)
		return err == nil && ok
	}
	return false
}

// billRemaining returns the outstanding balance for a single bill.
// Returns (remaining, found, error).
func (h *Handler) billRemaining(r *http.Request, billType string, billID int64) (float64, bool, error) {
	ctx := r.Context()

	switch billType {
	case "godown":
		b, err := h.app.Models.GodownBill.GetByID(ctx, billID)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "godown", billID)
		return b.TotalAmount - paid, true, nil

	case "majhi":
		b, err := h.app.Models.MajhiBill.GetByID(ctx, billID)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "majhi", billID)
		return b.TotalAmount - paid, true, nil

	case "customer_store":
		b, err := h.app.Models.CustomerStoreBill.GetByID(ctx, billID)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_store", billID)
		return b.TotalAmount - paid, true, nil

	case "customer_delivery":
		b, err := h.app.Models.CustomerDeliveryBill.GetByID(ctx, billID)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_delivery", billID)
		return b.TotalAmount - paid, true, nil

	case "customer_additional":
		b, err := h.app.Models.CustomerAdditionalBill.GetByID(ctx, billID)
		if err != nil {
			return 0, false, err
		}
		if b == nil {
			return 0, false, nil
		}
		paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_additional", billID)
		return b.Amount - paid, true, nil
	}

	return 0, false, nil
}

// billBelongsToEntity verifies a bill's customer/entity matches the
// entity the invoice is being created for.
func (h *Handler) billBelongsToEntity(r *http.Request, billType string, billID int64, entityType string, entityID int64) bool {
	ctx := r.Context()

	switch billType {
	case "customer_store":
		b, err := h.app.Models.CustomerStoreBill.GetByID(ctx, billID)
		if err != nil || b == nil {
			return false
		}
		return entityType == "customer" && b.CustomerID == entityID

	case "customer_delivery":
		b, err := h.app.Models.CustomerDeliveryBill.GetByID(ctx, billID)
		if err != nil || b == nil {
			return false
		}
		return entityType == "customer" && b.CustomerID == entityID

	case "customer_additional":
		b, err := h.app.Models.CustomerAdditionalBill.GetByID(ctx, billID)
		if err != nil || b == nil {
			return false
		}
		return entityType == "customer" && b.CustomerID == entityID

	case "majhi":
		b, err := h.app.Models.MajhiBill.GetByID(ctx, billID)
		if err != nil || b == nil {
			return false
		}
		return entityType == "majhi" && b.MajhiID == entityID

	case "godown":
		b, err := h.app.Models.GodownBill.GetByID(ctx, billID)
		if err != nil || b == nil {
			return false
		}
		return entityType == "godown" && b.GodownID == entityID
	}

	return false
}

// collectUnpaidBills gathers unpaid (remaining > 0) bills for the entity.
func (h *Handler) collectUnpaidBills(r *http.Request, entityType string, entityID int64) ([]UnbilledBill, error) {
	ctx := r.Context()
	out := []UnbilledBill{}

	switch entityType {

	case "customer":
		rows, err := h.app.Models.CustomerStoreBill.GetByCustomerID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		for _, b := range rows {
			paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_store", b.ID)
			remaining := b.TotalAmount - paid
			if remaining <= 0 {
				continue
			}
			out = append(out, UnbilledBill{
				BillType:  "customer_store",
				BillID:    b.ID,
				Label:     "Customer Store Bill #" + strconv.FormatInt(b.ID, 10) + " (" + b.MonthYear + ")",
				Total:     b.TotalAmount,
				Paid:      paid,
				Remaining: remaining,
				Date:      b.CreatedAt.Format(time.RFC3339),
			})
		}

		rows2, err := h.app.Models.CustomerDeliveryBill.GetByCustomerID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		for _, b := range rows2 {
			paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_delivery", b.ID)
			remaining := b.TotalAmount - paid
			if remaining <= 0 {
				continue
			}
			out = append(out, UnbilledBill{
				BillType:  "customer_delivery",
				BillID:    b.ID,
				Label:     "Customer Delivery Bill #" + strconv.FormatInt(b.ID, 10),
				Total:     b.TotalAmount,
				Paid:      paid,
				Remaining: remaining,
				Date:      b.CreatedAt.Format(time.RFC3339),
			})
		}

		rows3, err := h.app.Models.CustomerAdditionalBill.GetByCustomerID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		for _, b := range rows3 {
			paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "customer_additional", b.ID)
			remaining := b.Amount - paid
			if remaining <= 0 {
				continue
			}
			out = append(out, UnbilledBill{
				BillType:  "customer_additional",
				BillID:    b.ID,
				Label:     "Customer Additional Bill #" + strconv.FormatInt(b.ID, 10) + " — " + b.Description,
				Total:     b.Amount,
				Paid:      paid,
				Remaining: remaining,
				Date:      b.CreatedAt.Format(time.RFC3339),
			})
		}

	case "majhi":
		rows, err := h.app.Models.MajhiBill.GetByMajhiID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		for _, b := range rows {
			paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "majhi", b.ID)
			remaining := b.TotalAmount - paid
			if remaining <= 0 {
				continue
			}
			out = append(out, UnbilledBill{
				BillType:  "majhi",
				BillID:    b.ID,
				Label:     "Majhi Bill #" + strconv.FormatInt(b.ID, 10),
				Total:     b.TotalAmount,
				Paid:      paid,
				Remaining: remaining,
				Date:      b.CreatedAt.Format(time.RFC3339),
			})
		}

	case "godown":
		rows, err := h.app.Models.GodownBill.GetByGodownID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		for _, b := range rows {
			paid, _ := h.app.Models.BillPayment.GetTotalPaid(ctx, "godown", b.ID)
			remaining := b.TotalAmount - paid
			if remaining <= 0 {
				continue
			}
			out = append(out, UnbilledBill{
				BillType:  "godown",
				BillID:    b.ID,
				Label:     "Godown Bill #" + strconv.FormatInt(b.ID, 10) + " (" + b.MonthYear + ")",
				Total:     b.TotalAmount,
				Paid:      paid,
				Remaining: remaining,
				Date:      b.CreatedAt.Format(time.RFC3339),
			})
		}

	case "broker":
		return out, nil
	}

	return out, nil
}

func (h *Handler) buildInvoiceDetail(r *http.Request, invoiceID int64) (*InvoiceDetailResponse, error) {
	ctx := r.Context()

	invoice, err := h.app.Models.Invoice.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, nil
	}

	items, err := h.app.Models.InvoiceItem.GetByInvoiceID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.InvoiceItem{}
	}

	discounts, err := h.app.Models.Discount.GetByInvoiceID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if discounts == nil {
		discounts = []models.Discount{}
	}

	return &InvoiceDetailResponse{
		Invoice:   *invoice,
		Items:     items,
		Discounts: discounts,
	}, nil
}

func buildDiscountsFromInputs(reqs []discountInput, subtotal float64) (float64, []*models.Discount, error) {
	total := 0.0
	rows := make([]*models.Discount, 0, len(reqs))

	for _, d := range reqs {
		t := strings.ToLower(strings.TrimSpace(d.Type))
		if t != "flat" && t != "percent" {
			return 0, nil, errDiscountType
		}
		if d.Value < 0 {
			return 0, nil, errDiscountValue
		}

		var computed float64
		if t == "flat" {
			computed = d.Value
		} else {
			if d.Value > 100 {
				return 0, nil, errDiscountPercentTooBig
			}
			computed = subtotal * d.Value / 100.0
		}

		total += computed
		rows = append(rows, &models.Discount{
			Type:           t,
			Value:          d.Value,
			ComputedAmount: computed,
			Reason:         d.Reason,
		})
	}

	if total > subtotal {
		total = subtotal
	}

	return total, rows, nil
}
