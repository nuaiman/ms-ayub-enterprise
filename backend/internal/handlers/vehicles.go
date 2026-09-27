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
// VEHICLE MANAGEMENT
// =============================================================================

// CreateVehicleHandler - POST /api/vehicles
func (h *Handler) CreateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[VEHICLES] CreateVehicleHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		TransportID   int64   `json:"transport_id"`
		VehicleNumber string  `json:"vehicle_number"`
		BrokerID      int64   `json:"broker_id"`
		JomaCost      float64 `json:"joma_cost"`
		VehicleCost   float64 `json:"vehicle_cost"`
		OtherCost     float64 `json:"other_cost"`
		LabourCost    float64 `json:"labour_cost"`
		DemarageCost  float64 `json:"demarage_cost"`
		Notes         *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[VEHICLES] CreateVehicleHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TransportID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "transport_id is required")
		return
	}
	if req.VehicleNumber == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "vehicle_number is required")
		return
	}
	if req.BrokerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "broker_id is required")
		return
	}
	if req.JomaCost < 0 || req.VehicleCost < 0 ||
		req.OtherCost < 0 || req.LabourCost < 0 || req.DemarageCost < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "costs cannot be negative")
		return
	}

	transportExists, err := h.app.Models.Transport.Exists(r.Context(), req.TransportID)
	if err != nil {
		log.Printf("[VEHICLES] CreateVehicleHandler ERROR: failed to verify transport - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify transport")
		return
	}
	if !transportExists {
		utils.ErrorJson(w, http.StatusNotFound, "transport not found")
		return
	}

	brokerExists, err := h.app.Models.Broker.Exists(r.Context(), req.BrokerID)
	if err != nil {
		log.Printf("[VEHICLES] CreateVehicleHandler ERROR: failed to verify broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify broker")
		return
	}
	if !brokerExists {
		utils.ErrorJson(w, http.StatusNotFound, "broker not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[VEHICLES] CreateVehicleHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	vehicle := &models.Vehicle{
		UserID:            userID,
		TransportID:       req.TransportID,
		VehicleNumber:     req.VehicleNumber,
		BrokerID:          req.BrokerID,
		JomaCost:          req.JomaCost,
		VehicleCost:       req.VehicleCost,
		TotalPaidToBroker: 0,
		OtherCost:         req.OtherCost,
		LabourCost:        req.LabourCost,
		DemarageCost:      req.DemarageCost,
		Notes:             req.Notes,
	}

	id, err := h.app.Models.Vehicle.Insert(r.Context(), vehicle)
	if err != nil {
		log.Printf("[VEHICLES] CreateVehicleHandler ERROR: failed to insert vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create vehicle")
		return
	}

	vehicle.ID = id

	log.Printf("[VEHICLES] CreateVehicleHandler SUCCESS: created vehicle ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created vehicle #" + strconv.FormatInt(id, 10) + " for transport #" + strconv.FormatInt(req.TransportID, 10),
		EntityType:  "vehicles",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "vehicle created successfully", vehicle)
}

// GetVehicleHandler - GET /api/vehicles/{id}
func (h *Handler) GetVehicleHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[VEHICLES] GetVehicleHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[VEHICLES] GetVehicleHandler called - ID: %d", id)

	vehicle, err := h.app.Models.Vehicle.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[VEHICLES] GetVehicleHandler ERROR: failed to fetch vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch vehicle")
		return
	}
	if vehicle == nil {
		log.Printf("[VEHICLES] GetVehicleHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "vehicle not found")
		return
	}

	log.Printf("[VEHICLES] GetVehicleHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "vehicle details fetched", vehicle)
}

// GetAllVehiclesHandler - GET /api/vehicles
func (h *Handler) GetAllVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[VEHICLES] GetAllVehiclesHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	transportID := r.URL.Query().Get("transport_id")
	brokerID := r.URL.Query().Get("broker_id")
	vehicleNumber := r.URL.Query().Get("vehicle_number")

	var vehicles []models.Vehicle
	var err error

	switch {
	case transportID != "":
		id, _ := strconv.ParseInt(transportID, 10, 64)
		vehicles, err = h.app.Models.Vehicle.GetByTransportID(r.Context(), id)
	case brokerID != "":
		id, _ := strconv.ParseInt(brokerID, 10, 64)
		vehicles, err = h.app.Models.Vehicle.GetByBrokerID(r.Context(), id)
	case vehicleNumber != "":
		vehicles, err = h.app.Models.Vehicle.GetByVehicleNumber(r.Context(), vehicleNumber)
	default:
		vehicles, err = h.app.Models.Vehicle.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[VEHICLES] GetAllVehiclesHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch vehicles")
		return
	}

	if vehicles == nil {
		vehicles = []models.Vehicle{}
	}

	log.Printf("[VEHICLES] GetAllVehiclesHandler SUCCESS: fetched %d vehicles", len(vehicles))
	utils.SuccessJson(w, http.StatusOK, "vehicles fetched successfully", vehicles)
}

// UpdateVehicleHandler - PATCH /api/vehicles/{id}
func (h *Handler) UpdateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[VEHICLES] UpdateVehicleHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[VEHICLES] UpdateVehicleHandler called - ID: %d", id)

	type request struct {
		VehicleNumber *string  `json:"vehicle_number,omitempty"`
		BrokerID      *int64   `json:"broker_id,omitempty"`
		JomaCost      *float64 `json:"joma_cost,omitempty"`
		VehicleCost   *float64 `json:"vehicle_cost,omitempty"`
		OtherCost     *float64 `json:"other_cost,omitempty"`
		LabourCost    *float64 `json:"labour_cost,omitempty"`
		DemarageCost  *float64 `json:"demarage_cost,omitempty"`
		Notes         *string  `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[VEHICLES] UpdateVehicleHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vehicle, err := h.app.Models.Vehicle.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[VEHICLES] UpdateVehicleHandler ERROR: failed to fetch vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch vehicle")
		return
	}
	if vehicle == nil {
		log.Printf("[VEHICLES] UpdateVehicleHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "vehicle not found")
		return
	}

	if req.VehicleNumber != nil {
		if *req.VehicleNumber == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "vehicle_number cannot be empty")
			return
		}
		vehicle.VehicleNumber = *req.VehicleNumber
	}
	if req.BrokerID != nil {
		if *req.BrokerID == 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "broker_id is required")
			return
		}
		exists, err := h.app.Models.Broker.Exists(r.Context(), *req.BrokerID)
		if err != nil {
			log.Printf("[VEHICLES] UpdateVehicleHandler ERROR: failed to verify broker - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify broker")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "broker not found")
			return
		}
		vehicle.BrokerID = *req.BrokerID
	}
	if req.JomaCost != nil {
		if *req.JomaCost < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "joma_cost cannot be negative")
			return
		}
		vehicle.JomaCost = *req.JomaCost
	}
	if req.VehicleCost != nil {
		if *req.VehicleCost < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "vehicle_cost cannot be negative")
			return
		}
		vehicle.VehicleCost = *req.VehicleCost
	}
	if req.OtherCost != nil {
		if *req.OtherCost < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "other_cost cannot be negative")
			return
		}
		vehicle.OtherCost = *req.OtherCost
	}
	if req.LabourCost != nil {
		if *req.LabourCost < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "labour_cost cannot be negative")
			return
		}
		vehicle.LabourCost = *req.LabourCost
	}
	if req.DemarageCost != nil {
		if *req.DemarageCost < 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "demarage_cost cannot be negative")
			return
		}
		vehicle.DemarageCost = *req.DemarageCost
	}
	if req.Notes != nil {
		vehicle.Notes = req.Notes
	}

	if err := h.app.Models.Vehicle.Update(r.Context(), vehicle); err != nil {
		log.Printf("[VEHICLES] UpdateVehicleHandler ERROR: failed to update vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update vehicle")
		return
	}

	log.Printf("[VEHICLES] UpdateVehicleHandler SUCCESS: updated vehicle ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated vehicle #" + strconv.FormatInt(id, 10),
		EntityType:  "vehicles",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedVehicle, _ := h.app.Models.Vehicle.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "vehicle updated successfully", updatedVehicle)
}

// UpdateVehicleBrokerPaymentHandler - PATCH /api/vehicles/{id}/broker-payment
func (h *Handler) UpdateVehicleBrokerPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler called - ID: %d", id)

	type request struct {
		TotalPaidToBroker float64 `json:"total_paid_to_broker"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TotalPaidToBroker < 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "total_paid_to_broker cannot be negative")
		return
	}

	vehicle, err := h.app.Models.Vehicle.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler ERROR: failed to fetch vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch vehicle")
		return
	}
	if vehicle == nil {
		log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "vehicle not found")
		return
	}

	if err := h.app.Models.Vehicle.UpdateBrokerPayment(r.Context(), id, req.TotalPaidToBroker); err != nil {
		log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler ERROR: failed to update broker payment - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update broker payment")
		return
	}

	log.Printf("[VEHICLES] UpdateVehicleBrokerPaymentHandler SUCCESS: updated broker payment for vehicle ID=%d to %.2f", id, req.TotalPaidToBroker)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated broker payment for vehicle #" + strconv.FormatInt(id, 10) + " to " + strconv.FormatFloat(req.TotalPaidToBroker, 'f', 2, 64),
		EntityType:  "vehicles",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedVehicle, _ := h.app.Models.Vehicle.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "broker payment updated successfully", updatedVehicle)
}

// DeleteVehicleHandler - DELETE /api/vehicles/{id}
func (h *Handler) DeleteVehicleHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[VEHICLES] DeleteVehicleHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[VEHICLES] DeleteVehicleHandler called - ID: %d", id)

	vehicle, err := h.app.Models.Vehicle.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[VEHICLES] DeleteVehicleHandler ERROR: failed to fetch vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch vehicle")
		return
	}
	if vehicle == nil {
		log.Printf("[VEHICLES] DeleteVehicleHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "vehicle not found")
		return
	}

	if err := h.app.Models.Vehicle.Delete(r.Context(), id); err != nil {
		log.Printf("[VEHICLES] DeleteVehicleHandler ERROR: failed to delete vehicle - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete vehicle")
		return
	}

	log.Printf("[VEHICLES] DeleteVehicleHandler SUCCESS: deleted vehicle ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted vehicle #" + strconv.FormatInt(id, 10),
		EntityType:  "vehicles",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "vehicle deleted successfully", nil)
}
