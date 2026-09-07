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
// MAJHI MANAGEMENT
// =============================================================================

// CreateMajhiHandler - POST /api/majhis
func (h *Handler) CreateMajhiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MAJHIS] CreateMajhiHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Name  string  `json:"name"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[MAJHIS] CreateMajhiHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	// Check if majhi name already exists
	exists, err := h.app.Models.Majhi.ExistsByName(r.Context(), req.Name)
	if err != nil {
		log.Printf("[MAJHIS] CreateMajhiHandler ERROR: failed to check existing majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "majhi with this name already exists")
		return
	}

	majhi := &models.Majhi{
		Name:  req.Name,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	id, err := h.app.Models.Majhi.Insert(r.Context(), majhi)
	if err != nil {
		log.Printf("[MAJHIS] CreateMajhiHandler ERROR: failed to insert majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create majhi")
		return
	}

	majhi.ID = id

	log.Printf("[MAJHIS] CreateMajhiHandler SUCCESS: created majhi ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created majhi: " + req.Name,
		EntityType:  "majhis",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "majhi created successfully", majhi)
}

// GetMajhiHandler - GET /api/majhis/{id}
func (h *Handler) GetMajhiHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHIS] GetMajhiHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHIS] GetMajhiHandler called - ID: %d", id)

	majhi, err := h.app.Models.Majhi.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHIS] GetMajhiHandler ERROR: failed to fetch majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi")
		return
	}

	if majhi == nil {
		log.Printf("[MAJHIS] GetMajhiHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
		return
	}

	log.Printf("[MAJHIS] GetMajhiHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "majhi details fetched", majhi)
}

// GetAllMajhisHandler - GET /api/majhis
func (h *Handler) GetAllMajhisHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MAJHIS] GetAllMajhisHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")

	var majhis []models.Majhi
	var err error

	if search != "" {
		majhis, err = h.app.Models.Majhi.Search(r.Context(), search)
	} else {
		majhis, err = h.app.Models.Majhi.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[MAJHIS] GetAllMajhisHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhis")
		return
	}

	if majhis == nil {
		majhis = []models.Majhi{}
	}

	log.Printf("[MAJHIS] GetAllMajhisHandler SUCCESS: fetched %d majhis", len(majhis))
	utils.SuccessJson(w, http.StatusOK, "majhis fetched successfully", majhis)
}

// UpdateMajhiHandler - PATCH /api/majhis/{id}
func (h *Handler) UpdateMajhiHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHIS] UpdateMajhiHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHIS] UpdateMajhiHandler called - ID: %d", id)

	type request struct {
		Name  *string `json:"name,omitempty"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[MAJHIS] UpdateMajhiHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing majhi
	majhi, err := h.app.Models.Majhi.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHIS] UpdateMajhiHandler ERROR: failed to fetch majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi")
		return
	}
	if majhi == nil {
		log.Printf("[MAJHIS] UpdateMajhiHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
		return
	}

	// Apply updates
	if req.Name != nil {
		if *req.Name == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		// Check if new name conflicts with another majhi
		if *req.Name != majhi.Name {
			exists, err := h.app.Models.Majhi.ExistsByName(r.Context(), *req.Name)
			if err != nil {
				log.Printf("[MAJHIS] UpdateMajhiHandler ERROR: failed to check existing majhi - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify majhi")
				return
			}
			if exists {
				utils.ErrorJson(w, http.StatusConflict, "majhi with this name already exists")
				return
			}
		}
		majhi.Name = *req.Name
	}
	if req.Phone != nil {
		majhi.Phone = req.Phone
	}
	if req.Notes != nil {
		majhi.Notes = req.Notes
	}

	if err := h.app.Models.Majhi.Update(r.Context(), majhi); err != nil {
		log.Printf("[MAJHIS] UpdateMajhiHandler ERROR: failed to update majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update majhi")
		return
	}

	log.Printf("[MAJHIS] UpdateMajhiHandler SUCCESS: updated majhi ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated majhi #" + strconv.FormatInt(id, 10),
		EntityType:  "majhis",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated majhi
	updatedMajhi, _ := h.app.Models.Majhi.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "majhi updated successfully", updatedMajhi)
}

// DeleteMajhiHandler - DELETE /api/majhis/{id}
func (h *Handler) DeleteMajhiHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[MAJHIS] DeleteMajhiHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[MAJHIS] DeleteMajhiHandler called - ID: %d", id)

	// Check if majhi exists
	majhi, err := h.app.Models.Majhi.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[MAJHIS] DeleteMajhiHandler ERROR: failed to fetch majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch majhi")
		return
	}
	if majhi == nil {
		log.Printf("[MAJHIS] DeleteMajhiHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "majhi not found")
		return
	}

	// Check if majhi is referenced by any delivery_items
	// You can add this check if needed - depends on your business rules

	if err := h.app.Models.Majhi.Delete(r.Context(), id); err != nil {
		log.Printf("[MAJHIS] DeleteMajhiHandler ERROR: failed to delete majhi - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete majhi")
		return
	}

	log.Printf("[MAJHIS] DeleteMajhiHandler SUCCESS: deleted majhi ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted majhi: " + majhi.Name,
		EntityType:  "majhis",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "majhi deleted successfully", nil)
}
