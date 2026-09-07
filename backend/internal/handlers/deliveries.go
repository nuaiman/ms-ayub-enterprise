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
// DELIVERY MANAGEMENT
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

	// Validate customer exists if provided
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

	// Get current user ID
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default delivery date if not provided
	deliveryDate := req.DeliveryDate
	if deliveryDate == nil {
		now := time.Now()
		deliveryDate = &now
	}

	delivery := &models.Delivery{
		UserID:        userID,
		CustomerID:    req.CustomerID,
		DeliveryDate:  *deliveryDate,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		FromLocation:  req.FromLocation,
		ToLocation:    req.ToLocation,
		Notes:         req.Notes,
	}

	id, err := h.app.Models.Delivery.Insert(r.Context(), delivery)
	if err != nil {
		log.Printf("[DELIVERIES] CreateDeliveryHandler ERROR: failed to insert delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create delivery")
		return
	}

	delivery.ID = id

	log.Printf("[DELIVERIES] CreateDeliveryHandler SUCCESS: created delivery ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
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

	log.Printf("[DELIVERIES] GetDeliveryHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "delivery details fetched", delivery)
}

// GetAllDeliveriesHandler - GET /api/deliveries
func (h *Handler) GetAllDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELIVERIES] GetAllDeliveriesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	customerID := r.URL.Query().Get("customer_id")
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var deliveries []models.Delivery
	var err error

	switch {
	case search != "":
		deliveries, err = h.app.Models.Delivery.Search(r.Context(), search)
	case customerID != "":
		id, _ := strconv.ParseInt(customerID, 10, 64)
		deliveries, err = h.app.Models.Delivery.GetByCustomerID(r.Context(), id)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		deliveries, err = h.app.Models.Delivery.GetByUserID(r.Context(), id)
	case startDate != "" && endDate != "":
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		deliveries, err = h.app.Models.Delivery.GetByDateRange(r.Context(), start, end)
	default:
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

	// Get existing delivery
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

	// Validate customer if provided
	if req.CustomerID != nil {
		if *req.CustomerID != 0 {
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

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated delivery #" + strconv.FormatInt(id, 10),
		EntityType:  "deliveries",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated delivery
	updatedDelivery, _ := h.app.Models.Delivery.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "delivery updated successfully", updatedDelivery)
}

// DeleteDeliveryHandler - DELETE /api/deliveries/{id}
func (h *Handler) DeleteDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryHandler called - ID: %d", id)

	// Check if delivery exists
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

	// Check if delivery has delivery items
	// The ON DELETE CASCADE will handle this, but we should warn the user
	// We'll let the foreign key constraint handle it

	if err := h.app.Models.Delivery.Delete(r.Context(), id); err != nil {
		log.Printf("[DELIVERIES] DeleteDeliveryHandler ERROR: failed to delete delivery - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete delivery")
		return
	}

	log.Printf("[DELIVERIES] DeleteDeliveryHandler SUCCESS: deleted delivery ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
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
