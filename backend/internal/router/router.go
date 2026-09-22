package router

import (
	"backend/internal/app"
	"backend/internal/handlers"
	"backend/internal/middlewares"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func RegisterRouter(app *app.Application, handler *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	limiter := middlewares.NewRateLimiter(app.Config.RateLimit, app.Config.RateInterval)
	r.Use(limiter.LimitRate)

	protected := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.Authenticate(app.Config.JWTKey, h)
	}

	higherManagementOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return protected(middlewares.Authorize(app, "admin", "manager")(h))
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handler.HealthCheckHandler)

		// =====================================================
		// USERS
		// =====================================================
		r.Route("/users", func(r chi.Router) {
			// Public routes (no auth)
			r.Post("/login", handler.LoginHandler)
			r.Post("/refresh-session", handler.RefreshHandler)

			// Protected routes (auth required)
			r.Delete("/logout", protected(handler.LogoutHandler))
			r.Get("/current-user", protected(handler.GetCurrentUserHandler))
			r.Patch("/change-password", protected(handler.ChangePasswordHandler))

			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllUsersHandler))
			r.Post("/", higherManagementOnly(handler.CreateUserHandler))
			r.Patch("/reset-all-passwords", higherManagementOnly(handler.ResetAllPasswordsHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetUserHandler))
			r.Patch("/{id}/profile", higherManagementOnly(handler.UpdateUserProfileHandler))
			r.Patch("/{id}/change-password", higherManagementOnly(handler.ChangeUserPasswordHandler))
			r.Patch("/{id}/change-role", higherManagementOnly(handler.ChangeUserRoleHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleUserActiveHandler))
			r.Patch("/{id}/salary", higherManagementOnly(handler.UpdateUserSalaryHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteUserHandler))
		})

		// =====================================================
		// SALARIES
		// =====================================================
		r.Route("/salaries", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllSalariesHandler))
			r.Post("/", higherManagementOnly(handler.CreateSalaryHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetSalaryHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateSalaryHandler))
			r.Patch("/{id}/pay", higherManagementOnly(handler.MarkSalaryAsPaidHandler))
			r.Patch("/{id}/cancel", higherManagementOnly(handler.CancelSalaryHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteSalaryHandler))

			// Special routes
			r.Get("/employee/{id}", higherManagementOnly(handler.GetSalariesByEmployeeHandler))
			r.Get("/month/{month}", higherManagementOnly(handler.GetSalariesByMonthHandler))
		})

		// =====================================================
		// BROKERS
		// =====================================================
		r.Route("/brokers", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllBrokersHandler))
			r.Post("/", higherManagementOnly(handler.CreateBrokerHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetBrokerHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateBrokerHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteBrokerHandler))
		})

		// =====================================================
		// MAJHIS
		// =====================================================
		r.Route("/majhis", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllMajhisHandler))
			r.Post("/", higherManagementOnly(handler.CreateMajhiHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetMajhiHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateMajhiHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteMajhiHandler))
		})

		// =====================================================
		// GODOWNS
		// =====================================================
		r.Route("/godowns", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllGodownsHandler))
			r.Post("/", higherManagementOnly(handler.CreateGodownHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetGodownHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateGodownHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleGodownActiveHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteGodownHandler))
		})

		// =====================================================
		// RENTS
		// =====================================================
		r.Route("/rents", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllRentsHandler))
			r.Get("/current-month", higherManagementOnly(handler.GetCurrentMonthRentsHandler))
			r.Post("/", higherManagementOnly(handler.CreateRentHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetRentHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateRentHandler))
			r.Patch("/{id}/pay", higherManagementOnly(handler.MarkRentAsPaidHandler))
			r.Patch("/{id}/cancel", higherManagementOnly(handler.CancelRentHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteRentHandler))
		})

		// =====================================================
		// CUSTOMERS
		// =====================================================
		r.Route("/customers", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllCustomersHandler))
			r.Post("/", higherManagementOnly(handler.CreateCustomerHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetCustomerHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateCustomerHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteCustomerHandler))
		})

		// =====================================================
		// LOTS
		// =====================================================
		r.Route("/lots", func(r chi.Router) {
			// Admin/Manager only routes
			r.Get("/", higherManagementOnly(handler.GetAllLotsHandler))
			r.Post("/", higherManagementOnly(handler.CreateLotHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetLotHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateLotHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleLotActiveHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteLotHandler))

			r.Patch("/{id}/customer-payment", higherManagementOnly(handler.UpdateLotCustomerPaymentHandler))
			r.Patch("/{id}/customer-unload-payment", higherManagementOnly(handler.UpdateLotCustomerUnloadPaymentHandler)) // ADD THIS
			r.Patch("/{id}/majhi-payment", higherManagementOnly(handler.UpdateLotMajhiPaymentHandler))
		})

		// =====================================================
		// STORES
		// =====================================================
		r.Route("/stores", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllStoresHandler))
			r.Post("/", higherManagementOnly(handler.CreateStoreHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetStoreHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateStoreHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleStoreActiveHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteStoreHandler))
		})

		// =====================================================
		// DAMAGES
		// =====================================================
		r.Route("/damages", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllDamagesHandler))
			r.Post("/", higherManagementOnly(handler.CreateDamageHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetDamageHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDamageHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDamageHandler))
		})

		// =====================================================
		// DELIVERIES
		// =====================================================
		r.Route("/deliveries", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllDeliveriesHandler))
			r.Post("/", higherManagementOnly(handler.CreateDeliveryHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetDeliveryHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDeliveryHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDeliveryHandler))
		})

		// =====================================================
		// DELIVERY ITEMS
		// =====================================================
		r.Route("/delivery-items", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllDeliveryItemsHandler))
			r.Post("/", higherManagementOnly(handler.CreateDeliveryItemHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetDeliveryItemHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDeliveryItemHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDeliveryItemHandler))

			// Payment endpoints for delivery items
			r.Patch("/{id}/customer-unload-payment", higherManagementOnly(handler.UpdateDeliveryItemCustomerUnloadPaymentHandler))
			r.Patch("/{id}/majhi-payment", higherManagementOnly(handler.UpdateDeliveryItemMajhiPaymentHandler))
		})

		// =====================================================
		// TRANSPORTS
		// =====================================================
		r.Route("/transports", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllTransportsHandler))
			r.Post("/", higherManagementOnly(handler.CreateTransportHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetTransportHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateTransportHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteTransportHandler))
			r.Patch("/{id}/customer-payment", higherManagementOnly(handler.UpdateCustomerPaymentHandler)) // NEW
		})

		// =====================================================
		// VEHICLES
		// =====================================================
		r.Route("/vehicles", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllVehiclesHandler))
			r.Post("/", higherManagementOnly(handler.CreateVehicleHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetVehicleHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateVehicleHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteVehicleHandler))
			r.Patch("/{id}/broker-payment", higherManagementOnly(handler.UpdateVehicleBrokerPaymentHandler)) // NEW
		})

		// =====================================================
		// EXPENSES
		// =====================================================
		r.Route("/expenses", func(r chi.Router) {
			// Admin/Manager only routes (using higherManagementOnly)
			r.Get("/", higherManagementOnly(handler.GetAllExpensesHandler))
			r.Post("/", higherManagementOnly(handler.CreateExpenseHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetExpenseHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateExpenseHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteExpenseHandler))
		})

		// =====================================================
		// LOGS
		// =====================================================
		r.Route("/logs", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllLogsHandler))
		})

		// =====================================================
		// IMAGES
		// =====================================================
		r.Route("/images", func(r chi.Router) {
			r.Post("/{type}/{id}", protected(handler.UploadImageHandler))
			r.Delete("/{type}/{id}", protected(handler.DeleteImageHandler))
		})

		// =====================================================
		// BACKUP
		// =====================================================
		r.Get("/backup", protected(handler.BackupHandler))
	})

	// Serve static files from bucket directory
	r.Handle("/bucket/*", http.StripPrefix("/bucket/", http.FileServer(http.Dir("./bucket"))))

	// Serve SPA
	publicDir := "./public"
	staticFS := http.FileServer(http.Dir(publicDir))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ROUTER] NotFound handler called for: %s %s", r.Method, r.URL.Path)

		if strings.HasPrefix(r.URL.Path, "/api") {
			log.Printf("[ROUTER] API endpoint not found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		path := publicDir + r.URL.Path
		if _, err := os.Stat(path); err == nil {
			log.Printf("[ROUTER] Serving static file: %s", path)
			staticFS.ServeHTTP(w, r)
			return
		}

		log.Printf("[ROUTER] Serving index.html for SPA route: %s", r.URL.Path)
		http.ServeFile(w, r, publicDir+"/index.html")
	})

	log.Println("[ROUTER] Router initialized successfully")
	return r
}
