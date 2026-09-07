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
// CUSTOMER MANAGEMENT
// =============================================================================

// CreateCustomerHandler - POST /api/customers
func (h *Handler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMERS] CreateCustomerHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	type request struct {
		CompanyName   *string `json:"company_name,omitempty"`
		ContactPerson *string `json:"contact_person,omitempty"`
		Phone         string  `json:"phone"`
		Email         *string `json:"email,omitempty"`
		Address       *string `json:"address,omitempty"` // <-- New
		Notes         *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMERS] CreateCustomerHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields - at least company_name or contact_person
	if (req.CompanyName == nil || *req.CompanyName == "") && (req.ContactPerson == nil || *req.ContactPerson == "") {
		utils.ErrorJson(w, http.StatusBadRequest, "either company_name or contact_person is required")
		return
	}

	// Validate phone
	if req.Phone == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "phone is required")
		return
	}

	// Validate email if provided
	if req.Email != nil && *req.Email != "" {
		if !utils.IsValidEmail(*req.Email) {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}
	}

	// Check if phone already exists
	exists, err := h.app.Models.Customer.ExistsByPhone(r.Context(), req.Phone)
	if err != nil {
		log.Printf("[CUSTOMERS] CreateCustomerHandler ERROR: failed to check existing customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify customer")
		return
	}
	if exists {
		utils.ErrorJson(w, http.StatusConflict, "customer with this phone already exists")
		return
	}

	// Check if email already exists
	if req.Email != nil && *req.Email != "" {
		exists, err = h.app.Models.Customer.ExistsByEmail(r.Context(), *req.Email)
		if err != nil {
			log.Printf("[CUSTOMERS] CreateCustomerHandler ERROR: failed to check existing email - %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify email")
			return
		}
		if exists {
			utils.ErrorJson(w, http.StatusConflict, "customer with this email already exists")
			return
		}
	}

	customer := &models.Customer{
		CompanyName:   req.CompanyName,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address, // <-- New
		Notes:         req.Notes,
	}

	id, err := h.app.Models.Customer.Insert(r.Context(), customer)
	if err != nil {
		log.Printf("[CUSTOMERS] CreateCustomerHandler ERROR: failed to insert customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	customer.ID = id

	log.Printf("[CUSTOMERS] CreateCustomerHandler SUCCESS: created customer ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer: " + req.Phone,
		EntityType:  "customers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer created successfully", customer)
}

// GetCustomerHandler - GET /api/customers/{id}
func (h *Handler) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMERS] GetCustomerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMERS] GetCustomerHandler called - ID: %d", id)

	customer, err := h.app.Models.Customer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMERS] GetCustomerHandler ERROR: failed to fetch customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer")
		return
	}

	if customer == nil {
		log.Printf("[CUSTOMERS] GetCustomerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	log.Printf("[CUSTOMERS] GetCustomerHandler SUCCESS: ID=%d", id)
	utils.SuccessJson(w, http.StatusOK, "customer details fetched", customer)
}

// GetAllCustomersHandler - GET /api/customers
func (h *Handler) GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CUSTOMERS] GetAllCustomersHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	search := r.URL.Query().Get("search")

	var customers []models.Customer
	var err error

	if search != "" {
		customers, err = h.app.Models.Customer.Search(r.Context(), search)
	} else {
		customers, err = h.app.Models.Customer.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("[CUSTOMERS] GetAllCustomersHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customers")
		return
	}

	if customers == nil {
		customers = []models.Customer{}
	}

	log.Printf("[CUSTOMERS] GetAllCustomersHandler SUCCESS: fetched %d customers", len(customers))
	utils.SuccessJson(w, http.StatusOK, "customers fetched successfully", customers)
}

// UpdateCustomerHandler - PATCH /api/customers/{id}
func (h *Handler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMERS] UpdateCustomerHandler called - ID: %d", id)

	type request struct {
		CompanyName   *string `json:"company_name,omitempty"`
		ContactPerson *string `json:"contact_person,omitempty"`
		Phone         *string `json:"phone,omitempty"`
		Email         *string `json:"email,omitempty"`
		Address       *string `json:"address,omitempty"` // <-- New
		Notes         *string `json:"notes,omitempty"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: invalid request body - %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing customer
	customer, err := h.app.Models.Customer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: failed to fetch customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer")
		return
	}
	if customer == nil {
		log.Printf("[CUSTOMERS] UpdateCustomerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	// Apply updates
	if req.CompanyName != nil {
		customer.CompanyName = req.CompanyName
	}
	if req.ContactPerson != nil {
		customer.ContactPerson = req.ContactPerson
	}
	if req.Phone != nil {
		if *req.Phone == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "phone cannot be empty")
			return
		}
		// Check if new phone conflicts with another customer
		if *req.Phone != customer.Phone {
			exists, err := h.app.Models.Customer.ExistsByPhone(r.Context(), *req.Phone)
			if err != nil {
				log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: failed to check existing phone - %v", err)
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify phone")
				return
			}
			if exists {
				utils.ErrorJson(w, http.StatusConflict, "customer with this phone already exists")
				return
			}
		}
		customer.Phone = *req.Phone
	}
	if req.Email != nil {
		if *req.Email != "" {
			if !utils.IsValidEmail(*req.Email) {
				utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
				return
			}
			// Check if new email conflicts with another customer
			if *req.Email != "" && (customer.Email == nil || *req.Email != *customer.Email) {
				exists, err := h.app.Models.Customer.ExistsByEmail(r.Context(), *req.Email)
				if err != nil {
					log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: failed to check existing email - %v", err)
					utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify email")
					return
				}
				if exists {
					utils.ErrorJson(w, http.StatusConflict, "customer with this email already exists")
					return
				}
			}
		}
		customer.Email = req.Email
	}
	if req.Address != nil { // <-- New
		customer.Address = req.Address
	}
	if req.Notes != nil {
		customer.Notes = req.Notes
	}

	// Validate that either company_name or contact_person is set
	if (customer.CompanyName == nil || *customer.CompanyName == "") && (customer.ContactPerson == nil || *customer.ContactPerson == "") {
		utils.ErrorJson(w, http.StatusBadRequest, "either company_name or contact_person is required")
		return
	}

	if err := h.app.Models.Customer.Update(r.Context(), customer); err != nil {
		log.Printf("[CUSTOMERS] UpdateCustomerHandler ERROR: failed to update customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	log.Printf("[CUSTOMERS] UpdateCustomerHandler SUCCESS: updated customer ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer #" + strconv.FormatInt(id, 10),
		EntityType:  "customers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	// Fetch updated customer
	updatedCustomer, _ := h.app.Models.Customer.GetByID(r.Context(), id)
	utils.SuccessJson(w, http.StatusOK, "customer updated successfully", updatedCustomer)
}

// DeleteCustomerHandler - DELETE /api/customers/{id}
func (h *Handler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParamID(w, r)
	if !ok {
		log.Printf("[CUSTOMERS] DeleteCustomerHandler ERROR: invalid parameter")
		return
	}

	log.Printf("[CUSTOMERS] DeleteCustomerHandler called - ID: %d", id)

	// Check if customer exists
	customer, err := h.app.Models.Customer.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("[CUSTOMERS] DeleteCustomerHandler ERROR: failed to fetch customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customer")
		return
	}
	if customer == nil {
		log.Printf("[CUSTOMERS] DeleteCustomerHandler NOT FOUND: ID=%d", id)
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	// Check if customer is referenced by any items, deliveries, transports, invoices
	// You can add these checks if needed - depends on your business rules
	// For now, we'll let the foreign key constraint handle it

	if err := h.app.Models.Customer.Delete(r.Context(), id); err != nil {
		log.Printf("[CUSTOMERS] DeleteCustomerHandler ERROR: failed to delete customer - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	log.Printf("[CUSTOMERS] DeleteCustomerHandler SUCCESS: deleted customer ID=%d", id)

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(r.Context(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer: " + customer.Phone,
		EntityType:  "customers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer deleted successfully", nil)
}
