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
// GODOWN MANAGEMENT
// =============================================================================

// CreateGodownHandler - POST /api/godowns
func (h *Handler) CreateGodownHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GODOWNS] CreateGodownHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		Name  string  `json:"name"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[GODOWNS] CreateGodownHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	exists, err := h.app.Models.Godown.ExistsByName(r.Context(), req.Name)
	if err != nil {
		log.Printf("[GODOWNS] CreateGodownHandler ERROR: failed to check existing godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "godown with this name already exists")
		return
	}

	godown := &models.Godown{
		Name:  req.Name,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	id, err := h.app.Models.Godown.Insert(r.Context(), godown)
	if err != nil {
		log.Printf("[GODOWNS] CreateGodownHandler ERROR: failed to insert godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create godown")
		return
	}

	godown.ID = id

	log.Printf("[GODOWNS] CreateGodownHandler SUCCESS: created godown ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created godown: " + req.Name,
		EntityType:  "godowns",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "godown created successfully", godown)
}

// GetGodownHandler - GET /api/godowns/{id}
func (h *Handler) GetGodownHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWNS] GetGodownHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWNS] GetGodownHandler called - ID: %d", id)

	godown, err := h.app.Models.Godown.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWNS] GetGodownHandler ERROR: failed to fetch godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown")
		return
	}

	if godown == nil {
		log.Printf("[GODOWNS] GetGodownHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	log.Printf("[GODOWNS] GetGodownHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "godown details fetched", godown)
}

// GetAllGodownsHandler - GET /api/godowns
func (h *Handler) GetAllGodownsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GODOWNS] GetAllGodownsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")

	var godowns []models.Godown
	var err error

	if search != "" {
		godowns, err = h.app.Models.Godown.Search(r.Context(), search)
	} else {
		godowns, err = h.app.Models.Godown.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[GODOWNS] GetAllGodownsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godowns")
		return
	}

	if godowns == nil {
		godowns = []models.Godown{}
	}

	log.Printf("[GODOWNS] GetAllGodownsHandler SUCCESS: fetched %d godowns", len(godowns))
	utils.SuccessJson(w, http.StatusOK, "godowns fetched successfully", godowns)
}

// UpdateGodownHandler - PATCH /api/godowns/{id}
func (h *Handler) UpdateGodownHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWNS] UpdateGodownHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWNS] UpdateGodownHandler called - ID: %d", id)

	type request struct {
		Name  *string `json:"name,omitempty"`
		Phone *string `json:"phone,omitempty"`
		Notes *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[GODOWNS] UpdateGodownHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	godown, err := h.app.Models.Godown.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWNS] UpdateGodownHandler ERROR: failed to fetch godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown")
		return
	}
	if godown == nil {
		log.Printf("[GODOWNS] UpdateGodownHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	if req.Name != nil {
		if *req.Name == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		if *req.Name != godown.Name {
			exists, err := h.app.Models.Godown.ExistsByName(r.Context(), *req.Name)
			if err != nil {
				log.Printf("[GODOWNS] UpdateGodownHandler ERROR: failed to check existing godown - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify godown")
				return
			}
			if exists {
				utils.ErrorJson(w, http.StatusConflict, "godown with this name already exists")
				return
			}
		}
		godown.Name = *req.Name
	}
	if req.Phone != nil {
		godown.Phone = req.Phone
	}
	if req.Notes != nil {
		godown.Notes = req.Notes
	}

	if err := h.app.Models.Godown.Update(r.Context(), godown); err != nil {
		log.Printf("[GODOWNS] UpdateGodownHandler ERROR: failed to update godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update godown")
		return
	}

	log.Printf("[GODOWNS] UpdateGodownHandler SUCCESS: updated godown ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated godown #" + strconv.FormatInt(id, 10),
		EntityType:  "godowns",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedGodown, _ := h.app.Models.Godown.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "godown updated successfully", updatedGodown)
}

// DeleteGodownHandler - DELETE /api/godowns/{id}
func (h *Handler) DeleteGodownHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[GODOWNS] DeleteGodownHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[GODOWNS] DeleteGodownHandler called - ID: %d", id)

	godown, err := h.app.Models.Godown.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[GODOWNS] DeleteGodownHandler ERROR: failed to fetch godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch godown")
		return
	}
	if godown == nil {
		log.Printf("[GODOWNS] DeleteGodownHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "godown not found")
		return
	}

	if err := h.app.Models.Godown.Delete(r.Context(), id); err != nil {
		log.Printf("[GODOWNS] DeleteGodownHandler ERROR: failed to delete godown - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete godown")
		return
	}

	log.Printf("[GODOWNS] DeleteGodownHandler SUCCESS: deleted godown ID=%d", id)

	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	h.logAudit(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted godown: " + godown.Name,
		EntityType:  "godowns",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "godown deleted successfully", nil)
}
