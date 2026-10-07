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
// DELIVERY (Master)
// =============================================================================

// CreateDeliveryHandler - POST /api/deliveries
func (h *Handler) CreateDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELIVERIES] CreateDeliveryHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID    *int64     `json:"customer_id,omitempty"`
		DeliveryDate  *time.Time `json:"delivery_date,omitempty"`
		ReceiverName  *string    `json:"receiver_name,omitempty"`
		ReceiverPhone *string    `json:"receiver_phone,omitempty"`
		FromLocation  *string    `json:"from_location,omitempty"`
		ToLocation    *string    `json:"to_location,omitempty"`
		Notes         *string    `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// If customer_id provided, verify it exists.
	if req.CustomerID != nil {
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: failed to verify customer - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	deliveryDate := time.Now()
	if req.DeliveryDate != nil {
		deliveryDate = *req.DeliveryDate
	}

	delivery := &models.Delivery{
		UserID:        userID,
		CustomerID:    req.CustomerID,
		DeliveryDate:  deliveryDate,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		FromLocation:  req.FromLocation,
		ToLocation:    req.ToLocation,
		Notes:         req.Notes,
	}

	id, err := h.app.Models.Delivery.Insert(r.Context(), delivery)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: failed to insert delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create delivery: "+err.Error())
		return
	}
	delivery.ID = id

	log.Printf("[DELIVERIES] CreateDeliveryHandler SUCCESS: created delivery ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created delivery #" + strconv.FormatInt(id, 10),
		EntityType:  "deliveries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "delivery created successfully", delivery)
}

// GetDeliveryHandler - GET /api/deliveries/{id}
func (h *Handler) GetDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] GetDeliveryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] GetDeliveryHandler called - ID: %d", id)

	delivery, err := h.app.Models.Delivery.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] GetDeliveryHandler ERROR: failed to fetch delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery")
		return
	}
	if delivery == nil {
		log.Printf("[DELIVERIES] GetDeliveryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	items, err := h.app.Models.Delivery.GetItemsByDeliveryID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] GetDeliveryHandler ERROR: failed to fetch items - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery items")
		return
	}
	if items == nil {
		items = []models.DeliveryItem{}
	}

	log.Printf("[DELIVERIES] GetDeliveryHandler SUCCESS: ID=%d, items=%d", id, len(items))
	utils.SuccessJson(w, http.StatusOK, "delivery details fetched", map[string]any{
		"delivery": delivery,
		"items":    items,
	})
}

// GetAllDeliveriesHandler - GET /api/deliveries
func (h *Handler) GetAllDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELIVERIES] GetAllDeliveriesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	customerID := r.URL.Query().Get("customer_id")

	var deliveries []models.Delivery
	var err error

	if customerID != "" {
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		deliveries, err = h.app.Models.Delivery.GetByCustomerID(r.Context(), id)
	} else {
		deliveries, err = h.app.Models.Delivery.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[DELIVERIES] GetAllDeliveriesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch deliveries")
		return
	}
	if deliveries == nil {
		deliveries = []models.Delivery{}
	}

	log.Printf("[DELIVERIES] GetAllDeliveriesHandler SUCCESS: fetched %d deliveries", len(deliveries))
	utils.SuccessJson(w, http.StatusOK, "deliveries fetched successfully", deliveries)
}

// UpdateDeliveryHandler - PATCH /api/deliveries/{id}
func (h *Handler) UpdateDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] UpdateDeliveryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] UpdateDeliveryHandler called - ID: %d", id)

	type request struct {
		CustomerID    *int64     `json:"customer_id,omitempty"`
		DeliveryDate  *time.Time `json:"delivery_date,omitempty"`
		ReceiverName  *string    `json:"receiver_name,omitempty"`
		ReceiverPhone *string    `json:"receiver_phone,omitempty"`
		FromLocation  *string    `json:"from_location,omitempty"`
		ToLocation    *string    `json:"to_location,omitempty"`
		Notes         *string    `json:"notes,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	delivery, err := h.app.Models.Delivery.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryHandler ERROR: failed to fetch delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery")
		return
	}
	if delivery == nil {
		log.Printf("[DELIVERIES] UpdateDeliveryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	if req.CustomerID != nil {
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[DELIVERIES] UpdateDeliveryHandler ERROR: failed to verify customer - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
		delivery.CustomerID = req.CustomerID
	}
	if req.DeliveryDate != nil {
		delivery.DeliveryDate = *req.DeliveryDate
	}
	if req.ReceiverName != nil {
		delivery.ReceiverName = req.ReceiverName
	}
	if req.ReceiverPhone != nil {
		delivery.ReceiverPhone = req.ReceiverPhone
	}
	if req.FromLocation != nil {
		delivery.FromLocation = req.FromLocation
	}
	if req.ToLocation != nil {
		delivery.ToLocation = req.ToLocation
	}
	if req.Notes != nil {
		delivery.Notes = req.Notes
	}

	if err := h.app.Models.Delivery.Update(r.Context(), delivery); err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryHandler ERROR: failed to update delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update delivery")
		return
	}

	log.Printf("[DELIVERIES] UpdateDeliveryHandler SUCCESS: updated delivery ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated delivery #" + strconv.FormatInt(id, 10),
		EntityType:  "deliveries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.Delivery.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "delivery updated successfully", updated)
}

// DeleteDeliveryHandler - DELETE /api/deliveries/{id}
func (h *Handler) DeleteDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryHandler called - ID: %d", id)

	delivery, err := h.app.Models.Delivery.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler ERROR: failed to fetch delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery")
		return
	}
	if delivery == nil {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	if err := h.app.Models.Delivery.Delete(r.Context(), id); err != nil {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler ERROR: failed to delete delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete delivery: "+err.Error())
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryHandler SUCCESS: deleted delivery ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted delivery #" + strconv.FormatInt(id, 10),
		EntityType:  "deliveries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "delivery deleted successfully", nil)
}

// =============================================================================
// DELIVERY ITEMS
// =============================================================================

// CreateDeliveryItemHandler - POST /api/deliveries/{id}/items
func (h *Handler) CreateDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	deliveryID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: invalid delivery id")
		return
	}

	log.Printf("[DELIVERIES] CreateDeliveryItemHandler called - DeliveryID: %d", deliveryID)

	type request struct {
		StoreID       int64   `json:"store_id"`
		MajhiID       int64   `json:"majhi_id"`
		VehicleNumber *string `json:"vehicle_number,omitempty"`
		DriverNumber  *string `json:"driver_number,omitempty"`
		Quantity      float64 `json:"quantity"`
		Weight        float64 `json:"weight"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StoreID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "store_id is required")
		return
	}
	if req.MajhiID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "majhi_id is required")
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

	// Delivery must exist.
	deliveryExists, err := h.app.Models.Delivery.Exists(r.Context(), deliveryID)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: failed to verify delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify delivery")
		return
	}
	if !deliveryExists {
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	// Store must exist.
	storeExists, err := h.app.Models.Store.Exists(r.Context(), req.StoreID)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: failed to verify store - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
		return
	}
	if !storeExists {
		utils.ErrorJson(w, http.StatusNotFound, "store not found")
		return
	}

	// Majhi must exist.
	majhiExists, err := h.app.Models.Majhi.Exists(r.Context(), req.MajhiID)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: failed to verify majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
		return
	}
	if !majhiExists {
		utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
		return
	}

	item := &models.DeliveryItem{
		DeliveryID:    deliveryID,
		StoreID:       req.StoreID,
		MajhiID:       req.MajhiID,
		VehicleNumber: req.VehicleNumber,
		DriverNumber:  req.DriverNumber,
		Quantity:      req.Quantity,
		Weight:        req.Weight,
	}

	id, err := h.app.Models.Delivery.InsertItem(r.Context(), item)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryItemHandler ERROR: failed to insert item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create delivery item: "+err.Error())
		return
	}
	item.ID = id

	log.Printf("[DELIVERIES] CreateDeliveryItemHandler SUCCESS: created item ID=%d delivery=%d", id, deliveryID)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID: userID,
		Action: "create",
		Description: "Created delivery item #" + strconv.FormatInt(id, 10) +
			" for delivery #" + strconv.FormatInt(deliveryID, 10),
		EntityType: "delivery_items",
		EntityID:   id,
		IPAddress:  &ip,
		UserAgent:  &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "delivery item created successfully", item)
}

// GetDeliveryItemsHandler - GET /api/deliveries/{id}/items
func (h *Handler) GetDeliveryItemsHandler(w http.ResponseWriter, r *http.Request) {
	deliveryID, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] GetDeliveryItemsHandler ERROR: invalid delivery id")
		return
	}

	log.Printf("[DELIVERIES] GetDeliveryItemsHandler called - DeliveryID: %d", deliveryID)

	deliveryExists, err := h.app.Models.Delivery.Exists(r.Context(), deliveryID)
	if err != nil {
		log.Printf("[DELIVERIES] GetDeliveryItemsHandler ERROR: failed to verify delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify delivery")
		return
	}
	if !deliveryExists {
		utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
		return
	}

	items, err := h.app.Models.Delivery.GetItemsByDeliveryID(r.Context(), deliveryID)
	if err != nil {
		log.Printf("[DELIVERIES] GetDeliveryItemsHandler ERROR: failed to fetch items - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery items")
		return
	}
	if items == nil {
		items = []models.DeliveryItem{}
	}

	log.Printf("[DELIVERIES] GetDeliveryItemsHandler SUCCESS: fetched %d items", len(items))
	utils.SuccessJson(w, http.StatusOK, "delivery items fetched successfully", items)
}

// UpdateDeliveryItemHandler - PATCH /api/delivery-items/{id}
func (h *Handler) UpdateDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] UpdateDeliveryItemHandler called - ID: %d", id)

	type request struct {
		StoreID       *int64   `json:"store_id,omitempty"`
		MajhiID       *int64   `json:"majhi_id,omitempty"`
		VehicleNumber *string  `json:"vehicle_number,omitempty"`
		DriverNumber  *string  `json:"driver_number,omitempty"`
		Quantity      *float64 `json:"quantity,omitempty"`
		Weight        *float64 `json:"weight,omitempty"`
	}

	var req request
	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.app.Models.Delivery.GetItemByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERIES] UpdateDeliveryItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	if req.StoreID != nil {
		if *req.StoreID == 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "store_id cannot be 0")
			return
		}
		storeExists, err := h.app.Models.Store.Exists(r.Context(), *req.StoreID)
		if err != nil {
			log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: failed to verify store - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify store")
			return
		}
		if !storeExists {
			utils.ErrorJson(w, http.StatusNotFound, "store not found")
			return
		}
		item.StoreID = *req.StoreID
	}

	if req.MajhiID != nil {
		if *req.MajhiID == 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "majhi_id cannot be 0")
			return
		}
		majhiExists, err := h.app.Models.Majhi.Exists(r.Context(), *req.MajhiID)
		if err != nil {
			log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: failed to verify majhi - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
			return
		}
		if !majhiExists {
			utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
			return
		}
		item.MajhiID = *req.MajhiID
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
	if req.Weight != nil {
		if *req.Weight < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "weight cannot be negative")
			return
		}
		item.Weight = *req.Weight
	}

	if err := h.app.Models.Delivery.UpdateItem(r.Context(), item); err != nil {
		log.Printf("[DELIVERIES] UpdateDeliveryItemHandler ERROR: failed to update item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update delivery item: "+err.Error())
		return
	}

	log.Printf("[DELIVERIES] UpdateDeliveryItemHandler SUCCESS: updated item ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated delivery item #" + strconv.FormatInt(id, 10),
		EntityType:  "delivery_items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updated, _ := h.app.Models.Delivery.GetItemByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "delivery item updated successfully", updated)
}

// DeleteDeliveryItemHandler - DELETE /api/delivery-items/{id}
func (h *Handler) DeleteDeliveryItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] DeleteDeliveryItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryItemHandler called - ID: %d", id)

	item, err := h.app.Models.Delivery.GetItemByID(r.Context(), id)
	if err != nil {
		log.Printf("[DELIVERIES] DeleteDeliveryItemHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery item")
		return
	}
	if item == nil {
		log.Printf("[DELIVERIES] DeleteDeliveryItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "delivery item not found")
		return
	}

	if err := h.app.Models.Delivery.DeleteItem(r.Context(), id); err != nil {
		log.Printf("[DELIVERIES] DeleteDeliveryItemHandler ERROR: failed to delete item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete delivery item: "+err.Error())
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryItemHandler SUCCESS: deleted item ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
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
