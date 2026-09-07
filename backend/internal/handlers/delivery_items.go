package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"strconv"
)

// =============================================================================
// DELIVERY ITEM MANAGEMENT
// =============================================================================

// CreateDeliveryItemHandler - POST /api/delivery-items
func (h *Handler) CreateDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		DeliveryID               int64    `json:"delivery_id"`
		StoreID                  int64    `json:"store_id"`
		MajhiID                  *int64   `json:"majhi_id,omitempty"`
		ItemID                   int64    `json:"item_id"`
		LotID                    int64    `json:"lot_id"`
		VehicleNumber            *string  `json:"vehicle_number,omitempty"`
		DriverNumber             *string  `json:"driver_number,omitempty"`
		Quantity                 float64  `json:"quantity"`
		QuantityUnit             string   `json:"quantity_unit"`
		Weight                   float64  `json:"weight"`
		WeightUnit               string   `json:"weight_unit"`
		LoadingRate              float64  `json:"loading_rate"`
		MajhiCut                 float64  `json:"majhi_cut"`
		Notes                    *string  `json:"notes,omitempty"`
		CustomerChargeType       string   `json:"customer_charge_type,omitempty"`
		CustomerPaidUnloadAmount *float64 `json:"customer_paid_unload_amount,omitempty"`
		MajhiBillType            string   `json:"majhi_bill_type,omitempty"`
		MajhiTotalPaid           *float64 `json:"majhi_total_paid,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.DeliveryID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "delivery_id is required")
		return
	}
	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
		return
	}
	if req.ItemID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "item_id is required")
		return
	}
	if req.LotID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "lot_id is required")
		return
	}
	if req.Quantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}
	if req.Weight < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
		return
	}
	if req.Quantity == 0 && req.Weight == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "either quantity or weight must be greater than 0")
		return
	}

	// Check if delivery exists
	exists, err := h.app.Models.Delivery.Exists(r.Context(), req.DeliveryID)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to verify delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify delivery")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	// Check if store exists and get its lot info
	store, err := h.app.Models.Store.GetByID(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to fetch store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
		return
	}
	if store == nil {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Check if item exists
	exists, err = h.app.Models.Item.Exists(r.Context(), req.ItemID)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to verify item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify item")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	// Check if lot exists and belongs to item
	lot, err := h.app.Models.Lot.GetByID(r.Context(), req.LotID)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}
	if lot.ItemID != req.ItemID {
		utils.ErrorJson(w, http.StatusBadRequest, "lot does not belong to the specified item")
		return
	}

	// Check if majhi exists if provided
	if req.MajhiID != nil {
		exists, err := h.app.Models.Majhi.Exists(r.Context(), *req.MajhiID)
		if err != nil {
			log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to verify majhi - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
			return
		}
	}

	// Check store inventory (triggers will also check)
	if store.Quantity < req.Quantity {
		utils.ErrorJson(w, http.StatusBadRequest, "insufficient quantity in store")
		return
	}
	if store.Weight < req.Weight {
		utils.ErrorJson(w, http.StatusBadRequest, "insufficient weight in store")
		return
	}

	// Set default units
	quantityUnit := req.QuantityUnit
	if quantityUnit == "" {
		quantityUnit = "units"
	}
	weightUnit := req.WeightUnit
	if weightUnit == "" {
		weightUnit = "kg"
	}

	// Set default values for new fields
	customerChargeType := req.CustomerChargeType
	if customerChargeType == "" {
		customerChargeType = lot.CustomerChargeType
		if customerChargeType == "" {
			customerChargeType = "quantity"
		}
	}

	majhiBillType := req.MajhiBillType
	if majhiBillType == "" {
		majhiBillType = lot.MajhiBillType
		if majhiBillType == "" {
			majhiBillType = "quantity"
		}
	}

	customerPaidUnloadAmount := float64(0)
	if req.CustomerPaidUnloadAmount != nil {
		customerPaidUnloadAmount = *req.CustomerPaidUnloadAmount
	}

	majhiTotalPaid := float64(0)
	if req.MajhiTotalPaid != nil {
		majhiTotalPaid = *req.MajhiTotalPaid
	}

	// Get current user ID for audit
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	deliveryItem := &models.DeliveryItem{
		DeliveryID:               req.DeliveryID,
		StoreID:                  req.StoreID,
		MajhiID:                  req.MajhiID,
		ItemID:                   req.ItemID,
		LotID:                    req.LotID,
		VehicleNumber:            req.VehicleNumber,
		DriverNumber:             req.DriverNumber,
		Quantity:                 req.Quantity,
		QuantityUnit:             quantityUnit,
		Weight:                   req.Weight,
		WeightUnit:               weightUnit,
		LoadingRate:              req.LoadingRate,
		MajhiCut:                 req.MajhiCut,
		Notes:                    req.Notes,
		CustomerChargeType:       customerChargeType,
		CustomerPaidUnloadAmount: customerPaidUnloadAmount,
		MajhiBillType:            majhiBillType,
		MajhiTotalPaid:           majhiTotalPaid,
	}

	id, err := h.app.Models.DeliveryItem.Insert(r.Context(), deliveryItem)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler ERROR: failed to insert delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create delivery item")
		return
	}

	deliveryItem.ID = id

	log.Printf("[DELIVERY_ITEMS] CreateDeliveryItemHandler SUCCESS: created delivery item ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created delivery item #" + strconv.FormatInt(id, 10) + " for delivery #" + strconv.FormatInt(req.DeliveryID, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "delivery item created successfully", deliveryItem)
}

// GetDeliveryItemHandler - GET /api/delivery-items/{id}
func (h *Handler) GetDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] GetDeliveryItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERY_ITEMS] GetDeliveryItemHandler called - ID: %d", id)

	item, err := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] GetDeliveryItemHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}

	if item == nil {
		log.Printf("[DELIVERY_ITEMS] GetDeliveryItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	log.Printf("[DELIVERY_ITEMS] GetDeliveryItemHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "delivery item details fetched", item)
}

// GetAllDeliveryItemsHandler - GET /api/delivery-items
func (h *Handler) GetAllDeliveryItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELIVERY_ITEMS] GetAllDeliveryItemsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	deliveryID := r.URL.Query().Get("delivery_id")
	storeID := r.URL.Query().Get("store_id")
	lotID := r.URL.Query().Get("lot_id")

	var items []models.DeliveryItem
	var err error

	switch {
	case deliveryID != "":
		id, _ := strconv.ParseInt(deliveryID, 10, 64)
		items, err = h.app.Models.DeliveryItem.GetByDeliveryID(r.Context(), id)
	case storeID != "":
		id, _ := strconv.ParseInt(storeID, 10, 64)
		items, err = h.app.Models.DeliveryItem.GetByStoreID(r.Context(), id)
	case lotID != "":
		id, _ := strconv.ParseInt(lotID, 10, 64)
		items, err = h.app.Models.DeliveryItem.GetByLotID(r.Context(), id)
	default:
		items, err = h.app.Models.DeliveryItem.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[DELIVERY_ITEMS] GetAllDeliveryItemsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery items")
		return
	}

	if items == nil {
		items = []models.DeliveryItem{}
	}

	log.Printf("[DELIVERY_ITEMS] GetAllDeliveryItemsHandler SUCCESS: fetched %d delivery items", len(items))
	utils.SuccessJson(w, http.StatusOK, "delivery items fetched successfully", items)
}

// UpdateDeliveryItemHandler - PATCH /api/delivery-items/{id}
func (h *Handler) UpdateDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler called - ID: %d", id)

	type request struct {
		StoreID                  *int64   `json:"store_id,omitempty"`
		MajhiID                  *int64   `json:"majhi_id,omitempty"`
		VehicleNumber            *string  `json:"vehicle_number,omitempty"`
		DriverNumber             *string  `json:"driver_number,omitempty"`
		Quantity                 *float64 `json:"quantity,omitempty"`
		QuantityUnit             *string  `json:"quantity_unit,omitempty"`
		Weight                   *float64 `json:"weight,omitempty"`
		WeightUnit               *string  `json:"weight_unit,omitempty"`
		LoadingRate              *float64 `json:"loading_rate,omitempty"`
		MajhiCut                 *float64 `json:"majhi_cut,omitempty"`
		Notes                    *string  `json:"notes,omitempty"`
		CustomerChargeType       *string  `json:"customer_charge_type,omitempty"`
		CustomerPaidUnloadAmount *float64 `json:"customer_paid_unload_amount,omitempty"`
		MajhiBillType            *string  `json:"majhi_bill_type,omitempty"`
		MajhiTotalPaid           *float64 `json:"majhi_total_paid,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing delivery item
	item, err := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	// Validate store if provided
	if req.StoreID != nil {
		if *req.StoreID != 0 {
			exists, err := h.app.Models.Store.Exists(r.Context(), *req.StoreID)
			if err != nil {
				log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: failed to verify store - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
				return
			}
			if !exists {
				utils.ErrorJson(w, http.StatusNotFound, "store not found")
				return
			}
		}
		item.StoreID = *req.StoreID
	}
	if req.MajhiID != nil {
		if *req.MajhiID != 0 {
			exists, err := h.app.Models.Majhi.Exists(r.Context(), *req.MajhiID)
			if err != nil {
				log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: failed to verify majhi - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
				return
			}
			if !exists {
				utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
				return
			}
		}
		item.MajhiID = req.MajhiID
	}
	if req.VehicleNumber != nil {
		item.VehicleNumber = req.VehicleNumber
	}
	if req.DriverNumber != nil {
		item.DriverNumber = req.DriverNumber
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity cannot be negative")
			return
		}
		item.Quantity = *req.Quantity
	}
	if req.QuantityUnit != nil {
		if *req.QuantityUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity_unit cannot be empty")
			return
		}
		item.QuantityUnit = *req.QuantityUnit
	}
	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		item.Weight = *req.Weight
	}
	if req.WeightUnit != nil {
		if *req.WeightUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "weight_unit cannot be empty")
			return
		}
		item.WeightUnit = *req.WeightUnit
	}
	if req.LoadingRate != nil {
		if *req.LoadingRate < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "loading_rate cannot be negative")
			return
		}
		item.LoadingRate = *req.LoadingRate
	}
	if req.MajhiCut != nil {
		if *req.MajhiCut < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_cut cannot be negative")
			return
		}
		item.MajhiCut = *req.MajhiCut
	}
	if req.Notes != nil {
		item.Notes = req.Notes
	}
	if req.CustomerChargeType != nil {
		if *req.CustomerChargeType != "weight" && *req.CustomerChargeType != "quantity" {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_charge_type must be 'weight' or 'quantity'")
			return
		}
		item.CustomerChargeType = *req.CustomerChargeType
	}
	if req.CustomerPaidUnloadAmount != nil {
		if *req.CustomerPaidUnloadAmount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_paid_unload_amount cannot be negative")
			return
		}
		item.CustomerPaidUnloadAmount = *req.CustomerPaidUnloadAmount
	}
	if req.MajhiBillType != nil {
		if *req.MajhiBillType != "weight" && *req.MajhiBillType != "quantity" && *req.MajhiBillType != "job" {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_bill_type must be 'weight', 'quantity', or 'job'")
			return
		}
		item.MajhiBillType = *req.MajhiBillType
	}
	if req.MajhiTotalPaid != nil {
		if *req.MajhiTotalPaid < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_total_paid cannot be negative")
			return
		}
		item.MajhiTotalPaid = *req.MajhiTotalPaid
	}

	if err := h.app.Models.DeliveryItem.Update(r.Context(), item); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler ERROR: failed to update delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update delivery item")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemHandler SUCCESS: updated delivery item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated delivery item #" + strconv.FormatInt(id, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated item
	updatedItem, _ := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "delivery item updated successfully", updatedItem)
}

// DeleteDeliveryItemHandler - DELETE /api/delivery-items/{id}
func (h *Handler) DeleteDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler called - ID: %d", id)

	// Check if delivery item exists
	item, err := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	if err := h.app.Models.DeliveryItem.Delete(r.Context(), id); err != nil {
		log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler ERROR: failed to delete delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete delivery item")
		return
	}

	log.Printf("[DELIVERY_ITEMS] DeleteDeliveryItemHandler SUCCESS: deleted delivery item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted delivery item #" + strconv.FormatInt(id, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "delivery item deleted successfully", nil)
}

// UpdateDeliveryItemCustomerUnloadPaymentHandler - PATCH /api/delivery-items/{id}/customer-unload-payment
// Updates the customer unload payment amount for a specific delivery item
func (h *Handler) UpdateDeliveryItemCustomerUnloadPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler called - ID: %d", id)

	type request struct {
		CustomerPaidUnloadAmount float64 `json:"customer_paid_unload_amount"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerPaidUnloadAmount < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_paid_unload_amount cannot be negative")
		return
	}

	// Check if delivery item exists
	item, err := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	if err := h.app.Models.DeliveryItem.UpdateCustomerUnloadPayment(r.Context(), id, req.CustomerPaidUnloadAmount); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler ERROR: failed to update payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer unload payment")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemCustomerUnloadPaymentHandler SUCCESS: updated customer unload payment for delivery item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer unload payment for delivery item #" + strconv.FormatInt(id, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedItem, _ := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer unload payment updated successfully", updatedItem)
}

// UpdateDeliveryItemMajhiPaymentHandler - PATCH /api/delivery-items/{id}/majhi-payment
// Updates the majhi payment amount for a specific delivery item
func (h *Handler) UpdateDeliveryItemMajhiPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler called - ID: %d", id)

	type request struct {
		MajhiTotalPaid float64 `json:"majhi_total_paid"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.MajhiTotalPaid < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_total_paid cannot be negative")
		return
	}

	// Check if delivery item exists
	item, err := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler ERROR: failed to fetch delivery item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	if err := h.app.Models.DeliveryItem.UpdateMajhiPayment(r.Context(), id, req.MajhiTotalPaid); err != nil {
		log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler ERROR: failed to update payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update majhi payment")
		return
	}

	log.Printf("[DELIVERY_ITEMS] UpdateDeliveryItemMajhiPaymentHandler SUCCESS: updated majhi payment for delivery item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated majhi payment for delivery item #" + strconv.FormatInt(id, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedItem, _ := h.app.Models.DeliveryItem.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "majhi payment updated successfully", updatedItem)
}
