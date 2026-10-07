package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// =============================================================================
// LOT MANAGEMENT
// =============================================================================

// CreateLotHandler - POST /api/lots
func (h *Handler) CreateLotHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOTS] CreateLotHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID   int64  `json:"customer_id"`
		LotNumber    string `json:"lot_number"`
		ProductName  string `json:"product_name"`
		WeightUnit   string `json:"weight_unit"`
		QuantityUnit string `json:"quantity_unit"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_id is required")
		return
	}

	req.LotNumber = strings.TrimSpace(req.LotNumber)
	if req.LotNumber == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "lot_number is required")
		return
	}

	req.ProductName = strings.TrimSpace(req.ProductName)
	if req.ProductName == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "product_name is required")
		return
	}

	if req.WeightUnit == "" {
		req.WeightUnit = "kg"
	}
	if req.QuantityUnit == "" {
		req.QuantityUnit = "units"
	}

	// Customer must exist
	exists, err := h.app.Models.Customer.Exists(r.Context(), req.CustomerID)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to verify customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if !exists {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	// lot_number must be unique per customer
	exists, err = h.app.Models.Lot.ExistsByCustomerAndLot(r.Context(), req.CustomerID, req.LotNumber)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to check existing lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "lot number already exists for this customer")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("[LOTS] CreateLotHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	lot := &models.Lot{
		UserID:       userID,
		CustomerID:   req.CustomerID,
		LotNumber:    req.LotNumber,
		ProductName:  req.ProductName,
		WeightUnit:   req.WeightUnit,
		QuantityUnit: req.QuantityUnit,
	}

	id, err := h.app.Models.Lot.Insert(r.Context(), lot)
	if err != nil {
		log.Printf("[LOTS] CreateLotHandler ERROR: failed to insert lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create lot")
		return
	}

	lot.ID = id

	log.Printf("[LOTS] CreateLotHandler SUCCESS: created lot ID=%d", id)

	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created lot " + req.LotNumber + " (id=" + strconv.FormatInt(id, 10) + ")",
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "lot created successfully", lot)
}

// GetLotHandler - GET /api/lots/{id}
func (h *Handler) GetLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] GetLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] GetLotHandler called - ID: %d", id)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] GetLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] GetLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	log.Printf("[LOTS] GetLotHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "lot details fetched", lot)
}

// GetAllLotsHandler - GET /api/lots
func (h *Handler) GetAllLotsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOTS] GetAllLotsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	customerID := r.URL.Query().Get("customer_id")

	var lots []models.Lot
	var err error

	switch {
	case search != "":
		lots, err = h.app.Models.Lot.Search(r.Context(), search)
	case customerID != "":
		id, parseErr := strconv.ParseInt(customerID, 10, 64)
		if parseErr != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		lots, err = h.app.Models.Lot.GetByCustomerID(r.Context(), id)
	default:
		lots, err = h.app.Models.Lot.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[LOTS] GetAllLotsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lots")
		return
	}

	if lots == nil {
		lots = []models.Lot{}
	}

	log.Printf("[LOTS] GetAllLotsHandler SUCCESS: fetched %d lots", len(lots))
	utils.SuccessJson(w, http.StatusOK, "lots fetched successfully", lots)
}

// UpdateLotHandler - PATCH /api/lots/{id}
func (h *Handler) UpdateLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] UpdateLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] UpdateLotHandler called - ID: %d", id)

	type request struct {
		CustomerID   *int64  `json:"customer_id,omitempty"`
		LotNumber    *string `json:"lot_number,omitempty"`
		ProductName  *string `json:"product_name,omitempty"`
		WeightUnit   *string `json:"weight_unit,omitempty"`
		QuantityUnit *string `json:"quantity_unit,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] UpdateLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	// Determine effective values for uniqueness check
	effectiveCustomerID := lot.CustomerID
	effectiveLotNumber := lot.LotNumber

	if req.CustomerID != nil {
		if *req.CustomerID == 0 {
			utils.ErrorJson(w, http.StatusBadRequest, "customer_id cannot be 0")
			return
		}
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[LOTS] UpdateLotHandler ERROR: failed to verify customer - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
		effectiveCustomerID = *req.CustomerID
	}

	if req.LotNumber != nil {
		trimmed := strings.TrimSpace(*req.LotNumber)
		if trimmed == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "lot_number cannot be empty")
			return
		}
		effectiveLotNumber = trimmed
	}

	if req.ProductName != nil {
		trimmed := strings.TrimSpace(*req.ProductName)
		if trimmed == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "product_name cannot be empty")
			return
		}
		lot.ProductName = trimmed
	}

	// If either customer or lot number is changing, re-check uniqueness
	if effectiveCustomerID != lot.CustomerID || effectiveLotNumber != lot.LotNumber {
		exists, err := h.app.Models.Lot.ExistsByCustomerAndLot(r.Context(), effectiveCustomerID, effectiveLotNumber)
		if err != nil {
			log.Printf("[LOTS] UpdateLotHandler ERROR: failed to check existing lot - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify lot")
			return
		}
		if exists {
			utils.ErrorJson(w, http.StatusConflict, "lot number already exists for this customer")
			return
		}
	}

	lot.CustomerID = effectiveCustomerID
	lot.LotNumber = effectiveLotNumber

	if req.WeightUnit != nil {
		if *req.WeightUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "weight_unit cannot be empty")
			return
		}
		lot.WeightUnit = *req.WeightUnit
	}
	if req.QuantityUnit != nil {
		if *req.QuantityUnit == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity_unit cannot be empty")
			return
		}
		lot.QuantityUnit = *req.QuantityUnit
	}

	if err := h.app.Models.Lot.Update(r.Context(), lot); err != nil {
		log.Printf("[LOTS] UpdateLotHandler ERROR: failed to update lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update lot")
		return
	}

	log.Printf("[LOTS] UpdateLotHandler SUCCESS: updated lot ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedLot, _ := h.app.Models.Lot.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "lot updated successfully", updatedLot)
}

// DeleteLotHandler - DELETE /api/lots/{id}
func (h *Handler) DeleteLotHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[LOTS] DeleteLotHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[LOTS] DeleteLotHandler called - ID: %d", id)

	lot, err := h.app.Models.Lot.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[LOTS] DeleteLotHandler ERROR: failed to fetch lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch lot")
		return
	}
	if lot == nil {
		log.Printf("[LOTS] DeleteLotHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "lot not found")
		return
	}

	if err := h.app.Models.Lot.Delete(r.Context(), id); err != nil {
		log.Printf("[LOTS] DeleteLotHandler ERROR: failed to delete lot - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete lot")
		return
	}

	log.Printf("[LOTS] DeleteLotHandler SUCCESS: deleted lot ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted lot #" + strconv.FormatInt(id, 10),
		EntityType:  "lots",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "lot deleted successfully", nil)
}
