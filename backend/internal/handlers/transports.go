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
// TRANSPORT MANAGEMENT
// =============================================================================

// CreateTransportHandler - POST /api/transports
func (h *Handler) CreateTransportHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[TRANSPORTS] CreateTransportHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID             *int64     `json:"customer_id,omitempty"`
		FromLocation           string     `json:"from_location"`
		ToLocation             *string    `json:"to_location,omitempty"`
		VehicleQuantity        float64    `json:"vehicle_quantity"`
		DeliveryType           *string    `json:"delivery_type,omitempty"`
		Notes                  *string    `json:"notes,omitempty"`
		TransportDate          *time.Time `json:"transport_date,omitempty"`
		OfficeCommissionAmount float64    `json:"office_commission_amount"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[TRANSPORTS] CreateTransportHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.FromLocation == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "from_location is required")
		return
	}
	if req.VehicleQuantity < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "vehicle_quantity cannot be negative")
		return
	}

	// Validate delivery_type if provided
	if req.DeliveryType != nil && *req.DeliveryType != "" {
		if *req.DeliveryType != "local" && *req.DeliveryType != "district" {
			utils.ErrorJson(w, http.StatusBadRequest, "delivery_type must be 'local' or 'district'")
			return
		}
	}

	// Check if customer exists if provided
	if req.CustomerID != nil {
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[TRANSPORTS] CreateTransportHandler ERROR: failed to verify customer - %v", err)
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
		log.Printf("[TRANSPORTS] CreateTransportHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default transport date if not provided
	transportDate := req.TransportDate
	if transportDate == nil {
		now := time.Now()
		transportDate = &now
	}

	transport := &models.Transport{
		UserID:                 userID,
		CustomerID:             req.CustomerID,
		FromLocation:           req.FromLocation,
		ToLocation:             req.ToLocation,
		VehicleQuantity:        req.VehicleQuantity,
		DeliveryType:           req.DeliveryType,
		Notes:                  req.Notes,
		TransportDate:          *transportDate,
		OfficeCommissionAmount: req.OfficeCommissionAmount,
		CustomerTotalPaid:      0, // NEW - default to 0 on creation
	}

	id, err := h.app.Models.Transport.Insert(r.Context(), transport)
	if err != nil {
		log.Printf("[TRANSPORTS] CreateTransportHandler ERROR: failed to insert transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create transport")
		return
	}

	transport.ID = id

	log.Printf("[TRANSPORTS] CreateTransportHandler SUCCESS: created transport ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created transport #" + strconv.FormatInt(id, 10),
		EntityType:  "transports",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "transport created successfully", transport)
}

// GetTransportHandler - GET /api/transports/{id}
func (h *Handler) GetTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[TRANSPORTS] GetTransportHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[TRANSPORTS] GetTransportHandler called - ID: %d", id)

	transport, err := h.app.Models.Transport.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[TRANSPORTS] GetTransportHandler ERROR: failed to fetch transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transport")
		return
	}

	if transport == nil {
		log.Printf("[TRANSPORTS] GetTransportHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "transport not found")
		return
	}

	log.Printf("[TRANSPORTS] GetTransportHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "transport details fetched", transport)
}

// GetAllTransportsHandler - GET /api/transports
func (h *Handler) GetAllTransportsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[TRANSPORTS] GetAllTransportsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	customerID := r.URL.Query().Get("customer_id")
	userID := r.URL.Query().Get("user_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var transports []models.Transport
	var err error

	switch {
	case search != "":
		transports, err = h.app.Models.Transport.Search(r.Context(), search)
	case customerID != "":
		id, _ := strconv.ParseInt(customerID, 10, 64)
		transports, err = h.app.Models.Transport.GetByCustomerID(r.Context(), id)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		transports, err = h.app.Models.Transport.GetByUserID(r.Context(), id)
	case startDate != "" && endDate != "":
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		transports, err = h.app.Models.Transport.GetByDateRange(r.Context(), start, end)
	default:
		transports, err = h.app.Models.Transport.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[TRANSPORTS] GetAllTransportsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transports")
		return
	}

	if transports == nil {
		transports = []models.Transport{}
	}

	log.Printf("[TRANSPORTS] GetAllTransportsHandler SUCCESS: fetched %d transports", len(transports))
	utils.SuccessJson(w, http.StatusOK, "transports fetched successfully", transports)
}

// UpdateTransportHandler - PATCH /api/transports/{id}
func (h *Handler) UpdateTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[TRANSPORTS] UpdateTransportHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[TRANSPORTS] UpdateTransportHandler called - ID: %d", id)

	type request struct {
		CustomerID             *int64     `json:"customer_id,omitempty"`
		FromLocation           *string    `json:"from_location,omitempty"`
		ToLocation             *string    `json:"to_location,omitempty"`
		VehicleQuantity        *float64   `json:"vehicle_quantity,omitempty"`
		DeliveryType           *string    `json:"delivery_type,omitempty"`
		Notes                  *string    `json:"notes,omitempty"`
		TransportDate          *time.Time `json:"transport_date,omitempty"`
		OfficeCommissionAmount *float64   `json:"office_commission_amount,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[TRANSPORTS] UpdateTransportHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing transport
	transport, err := h.app.Models.Transport.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[TRANSPORTS] UpdateTransportHandler ERROR: failed to fetch transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transport")
		return
	}
	if transport == nil {
		log.Printf("[TRANSPORTS] UpdateTransportHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "transport not found")
		return
	}

	// Validate customer if provided
	if req.CustomerID != nil {
		if *req.CustomerID != 0 {
			exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
			if err != nil {
				log.Printf("[TRANSPORTS] UpdateTransportHandler ERROR: failed to verify customer - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
				return
			}
			if !exists {
				utils.ErrorJson(w, http.StatusNotFound, "customer not found")
				return
			}
		}
		transport.CustomerID = req.CustomerID
	}
	if req.FromLocation != nil {
		if *req.FromLocation == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "from_location cannot be empty")
			return
		}
		transport.FromLocation = *req.FromLocation
	}
	if req.ToLocation != nil {
		transport.ToLocation = req.ToLocation
	}
	if req.VehicleQuantity != nil {
		if *req.VehicleQuantity < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "vehicle_quantity cannot be negative")
			return
		}
		transport.VehicleQuantity = *req.VehicleQuantity
	}
	if req.DeliveryType != nil {
		if *req.DeliveryType != "" && *req.DeliveryType != "local" && *req.DeliveryType != "district" {
			utils.ErrorJson(w, http.StatusBadRequest, "delivery_type must be 'local' or 'district'")
			return
		}
		transport.DeliveryType = req.DeliveryType
	}
	if req.Notes != nil {
		transport.Notes = req.Notes
	}
	if req.TransportDate != nil {
		transport.TransportDate = *req.TransportDate
	}
	if req.OfficeCommissionAmount != nil {
		if *req.OfficeCommissionAmount < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "office_commission_amount cannot be negative")
			return
		}
		transport.OfficeCommissionAmount = *req.OfficeCommissionAmount
	}

	if err := h.app.Models.Transport.Update(r.Context(), transport); err != nil {
		log.Printf("[TRANSPORTS] UpdateTransportHandler ERROR: failed to update transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update transport")
		return
	}

	log.Printf("[TRANSPORTS] UpdateTransportHandler SUCCESS: updated transport ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated transport #" + strconv.FormatInt(id, 10),
		EntityType:  "transports",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated transport
	updatedTransport, _ := h.app.Models.Transport.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "transport updated successfully", updatedTransport)
}

// NEW: UpdateCustomerPaymentHandler - PATCH /api/transports/{id}/customer-payment
// Updates the customer total paid amount for a transport
func (h *Handler) UpdateCustomerPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler called - ID: %d", id)

	type request struct {
		CustomerTotalPaid float64 `json:"customer_total_paid"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerTotalPaid < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_total_paid cannot be negative")
		return
	}

	// Check if transport exists
	transport, err := h.app.Models.Transport.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler ERROR: failed to fetch transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transport")
		return
	}
	if transport == nil {
		log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "transport not found")
		return
	}

	if err := h.app.Models.Transport.UpdateCustomerPayment(r.Context(), id, req.CustomerTotalPaid); err != nil {
		log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler ERROR: failed to update customer payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer payment")
		return
	}

	log.Printf("[TRANSPORTS] UpdateCustomerPaymentHandler SUCCESS: updated customer payment for transport ID=%d to %.2f", id, req.CustomerTotalPaid)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer payment for transport #" + strconv.FormatInt(id, 10) + " to " + strconv.FormatFloat(req.CustomerTotalPaid, 'f', 2, 64),
		EntityType:  "transports",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated transport
	updatedTransport, _ := h.app.Models.Transport.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer payment updated successfully", updatedTransport)
}

// DeleteTransportHandler - DELETE /api/transports/{id}
func (h *Handler) DeleteTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[TRANSPORTS] DeleteTransportHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[TRANSPORTS] DeleteTransportHandler called - ID: %d", id)

	// Check if transport exists
	transport, err := h.app.Models.Transport.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[TRANSPORTS] DeleteTransportHandler ERROR: failed to fetch transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch transport")
		return
	}
	if transport == nil {
		log.Printf("[TRANSPORTS] DeleteTransportHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "transport not found")
		return
	}

	if err := h.app.Models.Transport.Delete(r.Context(), id); err != nil {
		log.Printf("[TRANSPORTS] DeleteTransportHandler ERROR: failed to delete transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete transport")
		return
	}

	log.Printf("[TRANSPORTS] DeleteTransportHandler SUCCESS: deleted transport ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted transport #" + strconv.FormatInt(id, 10),
		EntityType:  "transports",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "transport deleted successfully", nil)
}
