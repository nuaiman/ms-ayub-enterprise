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
// ITEM MANAGEMENT
// =============================================================================

// CreateItemHandler - POST /api/items
func (h *Handler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ITEMS] CreateItemHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CustomerID  *int64  `json:"customer_id,omitempty"`
		ProductName *string `json:"product_name,omitempty"`
		Category    *string `json:"category,omitempty"`
		Notes       *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[ITEMS] CreateItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate - at least product_name or category is required
	if (req.ProductName == nil || *req.ProductName == "") && (req.Category == nil || *req.Category == "") {
		utils.ErrorJson(w, http.StatusBadRequest, "either product_name or category is required")
		return
	}

	// Check if customer exists if provided
	if req.CustomerID != nil {
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[ITEMS] CreateItemHandler ERROR: failed to verify customer - %v", err)
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
		log.Printf("[ITEMS] CreateItemHandler ERROR: invalid user context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	item := &models.Item{
		UserID:      userID,
		CustomerID:  req.CustomerID,
		ProductName: req.ProductName,
		Category:    req.Category,
		IsActive:    true, // Default to active
		Notes:       req.Notes,
	}

	id, err := h.app.Models.Item.Insert(r.Context(), item)
	if err != nil {
		log.Printf("[ITEMS] CreateItemHandler ERROR: failed to insert item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	item.ID = id

	log.Printf("[ITEMS] CreateItemHandler SUCCESS: created item ID=%d", id)

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created item #" + strconv.FormatInt(id, 10),
		EntityType:  "items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "item created successfully", item)
}

// GetItemHandler - GET /api/items/{id}
func (h *Handler) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ITEMS] GetItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ITEMS] GetItemHandler called - ID: %d", id)

	item, err := h.app.Models.Item.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ITEMS] GetItemHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}

	if item == nil {
		log.Printf("[ITEMS] GetItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	log.Printf("[ITEMS] GetItemHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "item details fetched", item)
}

// GetAllItemsHandler - GET /api/items
func (h *Handler) GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ITEMS] GetAllItemsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")
	userID := r.URL.Query().Get("user_id")
	customerID := r.URL.Query().Get("customer_id")
	activeOnly := r.URL.Query().Get("active") == "true"

	var items []models.Item
	var err error

	switch {
	case search != "":
		items, err = h.app.Models.Item.Search(r.Context(), search)
	case userID != "":
		id, _ := strconv.ParseInt(userID, 10, 64)
		items, err = h.app.Models.Item.GetByUserID(r.Context(), id)
	case customerID != "":
		id, _ := strconv.ParseInt(customerID, 10, 64)
		items, err = h.app.Models.Item.GetByCustomerID(r.Context(), id)
	case activeOnly:
		items, err = h.app.Models.Item.GetActive(r.Context())
	default:
		items, err = h.app.Models.Item.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[ITEMS] GetAllItemsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch items")
		return
	}

	if items == nil {
		items = []models.Item{}
	}

	log.Printf("[ITEMS] GetAllItemsHandler SUCCESS: fetched %d items", len(items))
	utils.SuccessJson(w, http.StatusOK, "items fetched successfully", items)
}

// UpdateItemHandler - PATCH /api/items/{id}
func (h *Handler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ITEMS] UpdateItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ITEMS] UpdateItemHandler called - ID: %d", id)

	type request struct {
		CustomerID  *int64  `json:"customer_id,omitempty"`
		ProductName *string `json:"product_name,omitempty"`
		Category    *string `json:"category,omitempty"`
		Notes       *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[ITEMS] UpdateItemHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing item
	item, err := h.app.Models.Item.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ITEMS] UpdateItemHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}
	if item == nil {
		log.Printf("[ITEMS] UpdateItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	// Check if customer exists if provided
	if req.CustomerID != nil {
		exists, err := h.app.Models.Customer.Exists(r.Context(), *req.CustomerID)
		if err != nil {
			log.Printf("[ITEMS] UpdateItemHandler ERROR: failed to verify customer - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
			return
		}
		if !exists {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
		item.CustomerID = req.CustomerID
	}

	if req.ProductName != nil {
		item.ProductName = req.ProductName
	}
	if req.Category != nil {
		item.Category = req.Category
	}
	if req.Notes != nil {
		item.Notes = req.Notes
	}

	if err := h.app.Models.Item.Update(r.Context(), item); err != nil {
		log.Printf("[ITEMS] UpdateItemHandler ERROR: failed to update item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	log.Printf("[ITEMS] UpdateItemHandler SUCCESS: updated item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated item #" + strconv.FormatInt(id, 10),
		EntityType:  "items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated item
	updatedItem, _ := h.app.Models.Item.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "item updated successfully", updatedItem)
}

// ToggleItemActiveHandler - PATCH /api/items/{id}/toggle-active
func (h *Handler) ToggleItemActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ITEMS] ToggleItemActiveHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ITEMS] ToggleItemActiveHandler called - ID: %d", id)

	item, err := h.app.Models.Item.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ITEMS] ToggleItemActiveHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}
	if item == nil {
		log.Printf("[ITEMS] ToggleItemActiveHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	newStatus := !item.IsActive

	if err := h.app.Models.Item.ToggleActive(r.Context(), id, newStatus); err != nil {
		log.Printf("[ITEMS] ToggleItemActiveHandler ERROR: failed to toggle status - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item status")
		return
	}

	log.Printf("[ITEMS] ToggleItemActiveHandler SUCCESS: item ID=%d now active=%v", id, newStatus)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	action := "deactivated"
	if newStatus {
		action = "activated"
	}
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Item " + action + ": #" + strconv.FormatInt(id, 10),
		EntityType:  "items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	item.IsActive = newStatus
	utils.SuccessJson(w, http.StatusOK, "item status updated successfully", item)
}

// DeleteItemHandler - DELETE /api/items/{id}
func (h *Handler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[ITEMS] DeleteItemHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[ITEMS] DeleteItemHandler called - ID: %d", id)

	item, err := h.app.Models.Item.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[ITEMS] DeleteItemHandler ERROR: failed to fetch item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}
	if item == nil {
		log.Printf("[ITEMS] DeleteItemHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	if err := h.app.Models.Item.Delete(r.Context(), id); err != nil {
		log.Printf("[ITEMS] DeleteItemHandler ERROR: failed to delete item - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	log.Printf("[ITEMS] DeleteItemHandler SUCCESS: deleted item ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted item: #" + strconv.FormatInt(id, 10),
		EntityType:  "items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "item deleted successfully", nil)
}
