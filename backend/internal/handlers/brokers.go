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
// BROKER MANAGEMENT
// =============================================================================

// CreateBrokerHandler - POST /api/brokers
func (h *Handler) CreateBrokerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[BROKERS] CreateBrokerHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Name  string  `json:"name"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[BROKERS] CreateBrokerHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	// Check if broker name already exists
	exists, err := h.app.Models.Broker.ExistsByName(r.Context(), req.Name)
	if err != nil {
		log.Printf("[BROKERS] CreateBrokerHandler ERROR: failed to check existing broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify broker")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "broker with this name already exists")
		return
	}

	broker := &models.Broker{
		Name:  req.Name,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	id, err := h.app.Models.Broker.Insert(r.Context(), broker)
	if err != nil {
		log.Printf("[BROKERS] CreateBrokerHandler ERROR: failed to insert broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create broker")
		return
	}

	broker.ID = id

	log.Printf("[BROKERS] CreateBrokerHandler SUCCESS: created broker ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created broker: " + req.Name,
		EntityType:  "brokers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "broker created successfully", broker)
}

// GetBrokerHandler - GET /api/brokers/{id}
func (h *Handler) GetBrokerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[BROKERS] GetBrokerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[BROKERS] GetBrokerHandler called - ID: %d", id)

	broker, err := h.app.Models.Broker.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[BROKERS] GetBrokerHandler ERROR: failed to fetch broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch broker")
		return
	}

	if broker == nil {
		log.Printf("[BROKERS] GetBrokerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "broker not found")
		return
	}

	log.Printf("[BROKERS] GetBrokerHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "broker details fetched", broker)
}

// GetAllBrokersHandler - GET /api/brokers
func (h *Handler) GetAllBrokersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[BROKERS] GetAllBrokersHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")

	var brokers []models.Broker
	var err error

	if search != "" {
		brokers, err = h.app.Models.Broker.Search(r.Context(), search)
	} else {
		brokers, err = h.app.Models.Broker.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[BROKERS] GetAllBrokersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch brokers")
		return
	}

	if brokers == nil {
		brokers = []models.Broker{}
	}

	log.Printf("[BROKERS] GetAllBrokersHandler SUCCESS: fetched %d brokers", len(brokers))
	utils.SuccessJson(w, http.StatusOK, "brokers fetched successfully", brokers)
}

// UpdateBrokerHandler - PATCH /api/brokers/{id}
func (h *Handler) UpdateBrokerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[BROKERS] UpdateBrokerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[BROKERS] UpdateBrokerHandler called - ID: %d", id)

	type request struct {
		Name  *string `json:"name,omitempty"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[BROKERS] UpdateBrokerHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing broker
	broker, err := h.app.Models.Broker.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[BROKERS] UpdateBrokerHandler ERROR: failed to fetch broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch broker")
		return
	}
	if broker == nil {
		log.Printf("[BROKERS] UpdateBrokerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "broker not found")
		return
	}

	// Apply updates
	if req.Name != nil {
		if *req.Name == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		// Check if new name conflicts with another broker
		if *req.Name != broker.Name {
			exists, err := h.app.Models.Broker.ExistsByName(r.Context(), *req.Name)
			if err != nil {
				log.Printf("[BROKERS] UpdateBrokerHandler ERROR: failed to check existing broker - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify broker")
				return
			}
			if exists {
				utils.ErrorJson(w, http.StatusConflict, "broker with this name already exists")
				return
			}
		}
		broker.Name = *req.Name
	}
	if req.Phone != nil {
		broker.Phone = req.Phone
	}
	if req.Notes != nil {
		broker.Notes = req.Notes
	}

	if err := h.app.Models.Broker.Update(r.Context(), broker); err != nil {
		log.Printf("[BROKERS] UpdateBrokerHandler ERROR: failed to update broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update broker")
		return
	}

	log.Printf("[BROKERS] UpdateBrokerHandler SUCCESS: updated broker ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated broker #" + strconv.FormatInt(id, 10),
		EntityType:  "brokers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated broker
	updatedBroker, _ := h.app.Models.Broker.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "broker updated successfully", updatedBroker)
}

// DeleteBrokerHandler - DELETE /api/brokers/{id}
func (h *Handler) DeleteBrokerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[BROKERS] DeleteBrokerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[BROKERS] DeleteBrokerHandler called - ID: %d", id)

	// Check if broker exists
	broker, err := h.app.Models.Broker.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[BROKERS] DeleteBrokerHandler ERROR: failed to fetch broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch broker")
		return
	}
	if broker == nil {
		log.Printf("[BROKERS] DeleteBrokerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "broker not found")
		return
	}

	// Check if broker is referenced by any vehicles
	// You can add this check if needed - depends on your business rules

	if err := h.app.Models.Broker.Delete(r.Context(), id); err != nil {
		log.Printf("[BROKERS] DeleteBrokerHandler ERROR: failed to delete broker - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete broker")
		return
	}

	log.Printf("[BROKERS] DeleteBrokerHandler SUCCESS: deleted broker ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted broker: " + broker.Name,
		EntityType:  "brokers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "broker deleted successfully", nil)
}
