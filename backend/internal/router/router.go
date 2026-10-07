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

	adminOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return protected(middlewares.Authorize(app, "admin")(h))
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handler.HealthCheckHandler)

		// =====================================================
		// USERS
		// =====================================================
		r.Route("/users", func(r chi.Router) {
			r.Post("/login", handler.LoginHandler)
			r.Post("/refresh-session", handler.RefreshHandler)

			r.Delete("/logout", protected(handler.LogoutHandler))
			r.Get("/current-user", protected(handler.GetCurrentUserHandler))
			r.Patch("/change-password", protected(handler.ChangePasswordHandler))

			r.Get("/", higherManagementOnly(handler.GetAllUsersHandler))
			r.Post("/", higherManagementOnly(handler.CreateUserHandler))
			r.Patch("/reset-all-passwords", higherManagementOnly(handler.ResetAllPasswordsHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetUserHandler))
			r.Patch("/{id}/profile", protected(handler.UpdateUserProfileHandler))
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
			r.Get("/", higherManagementOnly(handler.GetAllSalariesHandler))
			r.Post("/", higherManagementOnly(handler.CreateSalaryHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetSalaryHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateSalaryHandler))
			r.Patch("/{id}/pay", higherManagementOnly(handler.MarkSalaryAsPaidHandler))
			r.Patch("/{id}/cancel", higherManagementOnly(handler.CancelSalaryHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteSalaryHandler))

			r.Get("/employee/{id}", higherManagementOnly(handler.GetSalariesByEmployeeHandler))
			r.Get("/month/{month}", higherManagementOnly(handler.GetSalariesByMonthHandler))
		})

		// =====================================================
		// BROKERS
		// =====================================================
		r.Route("/brokers", func(r chi.Router) {
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
			r.Delete("/{id}", higherManagementOnly(handler.DeleteGodownHandler))
		})

		// =====================================================
		// CUSTOMERS
		// =====================================================
		r.Route("/customers", func(r chi.Router) {
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
			r.Get("/", higherManagementOnly(handler.GetAllLotsHandler))
			r.Post("/", higherManagementOnly(handler.CreateLotHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetLotHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateLotHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteLotHandler))

			r.Get("/{id}/transfers", higherManagementOnly(handler.GetLotTransfersHandler))
			r.Post("/{id}/transfer", higherManagementOnly(handler.CreateLotTransferHandler))
		})

		// =====================================================
		// STORES
		// =====================================================
		r.Route("/stores", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllStoresHandler))
			r.Post("/", higherManagementOnly(handler.CreateStoreHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetStoreHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateStoreHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleStoreActiveHandler))
			r.Post("/{id}/refresh-bills", higherManagementOnly(handler.RefreshStoreBillsHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteStoreHandler))

			r.Get("/{id}/adjustments", higherManagementOnly(handler.GetStoreAdjustmentsHandler))
			r.Post("/{id}/adjustments", higherManagementOnly(handler.CreateStoreAdjustmentHandler))

			r.Get("/{id}/damages", higherManagementOnly(handler.GetStoreDamagesHandler))

			r.Get("/{id}/transfers", higherManagementOnly(handler.GetStoreTransfersHandler))
			r.Post("/{id}/transfer", higherManagementOnly(handler.CreateStoreTransferHandler))
		})

		// =====================================================
		// STORE ADJUSTMENTS
		// =====================================================
		r.Route("/store-adjustments", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllStoreAdjustmentsHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetStoreAdjustmentHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteStoreAdjustmentHandler))
		})

		// =====================================================
		// STORE TRANSFERS
		// =====================================================
		r.Route("/store-transfers", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllStoreTransfersHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetStoreTransferHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteStoreTransferHandler))
		})

		// =====================================================
		// LOT TRANSFERS
		// =====================================================
		r.Route("/lot-transfers", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllLotTransfersHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetLotTransferHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteLotTransferHandler))
		})

		// =====================================================
		// EXPENSES
		// =====================================================
		r.Route("/expenses", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllExpensesHandler))
			r.Post("/", higherManagementOnly(handler.CreateExpenseHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetExpenseHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateExpenseHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteExpenseHandler))
		})

		// =====================================================
		// INCOMES
		// =====================================================
		r.Route("/incomes", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllIncomesHandler))
			r.Post("/", higherManagementOnly(handler.CreateIncomeHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetIncomeHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateIncomeHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteIncomeHandler))
		})

		// =====================================================
		// DELIVERIES
		// =====================================================
		r.Route("/deliveries", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllDeliveriesHandler))
			r.Post("/", higherManagementOnly(handler.CreateDeliveryHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetDeliveryHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDeliveryHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDeliveryHandler))

			r.Get("/{id}/items", higherManagementOnly(handler.GetDeliveryItemsHandler))
			r.Post("/{id}/items", higherManagementOnly(handler.CreateDeliveryItemHandler))
		})

		// =====================================================
		// DELIVERY ITEMS (standalone operations)
		// =====================================================
		r.Route("/delivery-items", func(r chi.Router) {
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDeliveryItemHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDeliveryItemHandler))
		})

		// =====================================================
		// DAMAGES
		// =====================================================
		r.Route("/damages", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllDamagesHandler))
			r.Post("/", higherManagementOnly(handler.CreateDamageHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetDamageHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateDamageHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteDamageHandler))
		})

		// =====================================================
		// LOGS
		// =====================================================
		r.Route("/logs", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllLogsHandler))
		})

		// =====================================================
		// GODOWN BILLS (MONTHLY)
		// =====================================================
		r.Route("/godown-bills", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllGodownBillsHandler))
			r.Get("/current-month", higherManagementOnly(handler.GetCurrentMonthGodownBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateGodownBillHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetGodownBillHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateGodownBillHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteGodownBillHandler))

			r.Get("/{id}/payments", higherManagementOnly(handler.GetGodownBillPaymentsHandler))
			r.Post("/{id}/payments", higherManagementOnly(handler.CreateGodownBillPaymentHandler))
		})

		// =====================================================
		// MAJHI BILLS (ONE-TIME - store or delivery item)
		// =====================================================
		r.Route("/majhi-bills", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllMajhiBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateMajhiBillHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetMajhiBillHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateMajhiBillHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteMajhiBillHandler))

			r.Get("/{id}/payments", higherManagementOnly(handler.GetMajhiBillPaymentsHandler))
			r.Post("/{id}/payments", higherManagementOnly(handler.CreateMajhiBillPaymentHandler))
		})

		// =====================================================
		// CUSTOMER STORE BILLS (MONTHLY)
		// =====================================================
		r.Route("/customer-store-bills", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllCustomerStoreBillsHandler))
			r.Get("/current-month", higherManagementOnly(handler.GetCurrentMonthCustomerStoreBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateCustomerStoreBillHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetCustomerStoreBillHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateCustomerStoreBillHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteCustomerStoreBillHandler))

			r.Get("/{id}/payments", higherManagementOnly(handler.GetCustomerStoreBillPaymentsHandler))
			r.Post("/{id}/payments", higherManagementOnly(handler.CreateCustomerStoreBillPaymentHandler))
		})

		// =====================================================
		// CUSTOMER DELIVERY BILLS (ONE-TIME)
		// =====================================================
		r.Route("/customer-delivery-bills", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllCustomerDeliveryBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateCustomerDeliveryBillHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetCustomerDeliveryBillHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateCustomerDeliveryBillHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteCustomerDeliveryBillHandler))

			r.Get("/{id}/payments", higherManagementOnly(handler.GetCustomerDeliveryBillPaymentsHandler))
			r.Post("/{id}/payments", higherManagementOnly(handler.CreateCustomerDeliveryBillPaymentHandler))
		})

		// =====================================================
		// BILL PAYMENTS (standalone operations)
		// =====================================================
		r.Route("/bill-payments", func(r chi.Router) {
			r.Delete("/{id}", higherManagementOnly(handler.DeleteBillPaymentHandler))
		})

		// =====================================================
		// CUSTOMER ADDITIONAL BILLS (ONE-TIME)
		// =====================================================
		r.Route("/customer-additional-bills", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllCustomerAdditionalBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateCustomerAdditionalBillHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetCustomerAdditionalBillHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateCustomerAdditionalBillHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteCustomerAdditionalBillHandler))

			r.Get("/{id}/payments", higherManagementOnly(handler.GetCustomerAdditionalBillPaymentsHandler))
			r.Post("/{id}/payments", higherManagementOnly(handler.CreateCustomerAdditionalBillPaymentHandler))
		})

		// =====================================================
		// INVOICES
		// =====================================================
		r.Route("/invoices", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllInvoicesHandler))
			r.Get("/unbilled", higherManagementOnly(handler.GetUnbilledBillsHandler))
			r.Post("/", higherManagementOnly(handler.CreateInvoiceHandler))
			r.Get("/{id}", higherManagementOnly(handler.GetInvoiceHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateInvoiceHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteInvoiceHandler))
		})

		// =====================================================
		// IMAGES
		// =====================================================
		r.Route("/images", func(r chi.Router) {
			r.Post("/{type}/{id}", protected(handler.UploadImageHandler))
			r.Delete("/{type}/{id}", protected(handler.DeleteImageHandler))
		})

		// =====================================================
		// BACKUP - admin only
		// =====================================================
		r.Get("/backup", adminOnly(handler.BackupHandler))
	})

	r.Handle("/bucket/*", http.StripPrefix("/bucket/", http.FileServer(http.Dir("./bucket"))))

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
