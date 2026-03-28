package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/your-org/jvairv2/docs" // Importación de documentación Swagger
	"github.com/your-org/jvairv2/pkg/domain/user"
	"github.com/your-org/jvairv2/pkg/rest/handler"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
	alertHandler "github.com/your-org/jvairv2/pkg/rest/handler/alert"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	customerHandler "github.com/your-org/jvairv2/pkg/rest/handler/customer"
	dashboardHandler "github.com/your-org/jvairv2/pkg/rest/handler/dashboard"
	emailTemplateHandler "github.com/your-org/jvairv2/pkg/rest/handler/email_template"
	invoiceHandler "github.com/your-org/jvairv2/pkg/rest/handler/invoice"
	invoicePaymentHandler "github.com/your-org/jvairv2/pkg/rest/handler/invoice_payment"
	jobHandler "github.com/your-org/jvairv2/pkg/rest/handler/job"
	jobActivityLogHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_activity_log"
	jobCategoryHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_category"
	jobEmailHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_email"
	jobEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_equipment"
	jobPriorityHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_priority"
	jobRateHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_rate"
	jobRateStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_rate_status"
	jobResidentHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_resident"
	jobSMSHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_sms"
	jobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_status"
	jobTaskHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_task"
	jobVisitHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_visit"
	"github.com/your-org/jvairv2/pkg/rest/handler/new_dashboard"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	propertyHandler "github.com/your-org/jvairv2/pkg/rest/handler/property"
	propEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/property_equipment"
	quoteHandler "github.com/your-org/jvairv2/pkg/rest/handler/quote"
	quoteStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/quote_status"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	settingsHandler "github.com/your-org/jvairv2/pkg/rest/handler/settings"
	smsTemplateHandler "github.com/your-org/jvairv2/pkg/rest/handler/sms_template"
	supervisorHandler "github.com/your-org/jvairv2/pkg/rest/handler/supervisor"
	taskStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/task_status"
	techJobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/technician_job_status"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
	warrantyHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty"
	warrantyClaimHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim"
	warrantyClaimStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim_status"
	warrantyClaimTypeHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim_type"
	warrantyEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_equipment"
	warrantyStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_status"
	warrantyTypeHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_type"
	workflowHandler "github.com/your-org/jvairv2/pkg/rest/handler/workflow"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

// New crea un nuevo router HTTP con las rutas configuradas
func New(
	healthHandler *handler.HealthHandler,
	authHandler *authHandler.Handler,
	passwordSecurityHandler *authHandler.PasswordSecurityHandler,
	userHandler *userHandler.Handler,
	roleHandler *roleHandler.Handler,
	abilityHandler *abilityHandler.Handler,
	assignedRoleHandler *assignedRoleHandler.Handler,
	permissionHandler *permissionHandler.Handler,
	settingsHandler *settingsHandler.Handler,
	emailTemplateHandler *emailTemplateHandler.Handler,
	workflowHandler *workflowHandler.Handler,
	customerHandler *customerHandler.Handler,
	propertyHandler *propertyHandler.Handler,
	jobHandler *jobHandler.Handler,
	jobCategoryHandler *jobCategoryHandler.Handler,
	jobStatusHandler *jobStatusHandler.Handler,
	jobPriorityHandler *jobPriorityHandler.Handler,
	techJobStatusHandler *techJobStatusHandler.Handler,
	taskStatusHandler *taskStatusHandler.Handler,
	quoteHandler *quoteHandler.Handler,
	quoteStatusHandler *quoteStatusHandler.Handler,
	supervisorHandler *supervisorHandler.Handler,
	propEquipHandler *propEquipHandler.Handler,
	jobEquipHandler *jobEquipHandler.Handler,
	invoiceHandler *invoiceHandler.Handler,
	invoicePaymentHandler *invoicePaymentHandler.Handler,
	warrantyTypeHandler *warrantyTypeHandler.Handler,
	warrantyStatusHandler *warrantyStatusHandler.Handler,
	warrantyClaimTypeHandler *warrantyClaimTypeHandler.Handler,
	warrantyClaimStatusHandler *warrantyClaimStatusHandler.Handler,
	warrantyHandler *warrantyHandler.Handler,
	warrantyEquipHandler *warrantyEquipHandler.Handler,
	warrantyClaimHandler *warrantyClaimHandler.Handler,
	jobVisitHandler *jobVisitHandler.Handler,
	jobActivityLogHandler *jobActivityLogHandler.Handler,
	jobEmailHandler *jobEmailHandler.Handler,
	jobResidentHandler *jobResidentHandler.Handler,
	jobRateStatusHandler *jobRateStatusHandler.Handler,
	jobSMSHandler *jobSMSHandler.Handler,
	smsTemplateHandler *smsTemplateHandler.Handler,
	jobTaskHandler *jobTaskHandler.Handler,
	jobRateHandler *jobRateHandler.Handler,
	alertHandler *alertHandler.Handler,
	dashboardHandler *dashboardHandler.Handler,
	newDashboardHdlr *new_dashboard.Handler,
	authMiddleware *middleware.AuthMiddleware,
	userUseCase *user.UseCase, // Añadir esta dependencia
) *chi.Mux {
	r := chi.NewRouter()
	// Middlewares globales
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	// Rutas públicas
	r.Group(func(r chi.Router) {
		// Health check
		r.Get("/health", healthHandler.Check)
		// Rutas de autenticación
		RegisterAuthRoutes(r, authHandler, passwordSecurityHandler)
		// Swagger UI
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // URL para acceder a la documentación JSON
		))
	})
	// Rutas protegidas que requieren autenticación
	r.Group(func(r chi.Router) {
		// Middleware de autenticación
		r.Use(authMiddleware.Authenticate)

		// Middleware de habilidades - añadir esto
		r.Use(middleware.WithAbilities(userUseCase))

		// Ruta protegida de cambio de contraseña
		r.Post("/auth/change-password", passwordSecurityHandler.ChangePassword)

		// API v1
		r.Route("/api/v1", func(r chi.Router) {
			// Rutas de usuarios
			RegisterUserRoutes(r, userHandler)
			// Rutas de roles
			RegisterRoleRoutes(r, roleHandler)
			// Rutas de abilities
			RegisterAbilityRoutes(r, abilityHandler)
			// Rutas de assigned-roles
			RegisterAssignedRoleRoutes(r, assignedRoleHandler)
			// Rutas de permisos
			RegisterPermissionRoutes(r, permissionHandler)
			// Rutas de configuraciones
			SetupSettingsRoutes(r, settingsHandler, authMiddleware)
			// Rutas de plantillas de email
			emailTemplateHandler.RegisterRoutes(r)
			// Rutas de workflows
			SetupWorkflowRoutes(r, workflowHandler, authMiddleware)
			// Rutas de customers
			RegisterCustomerRoutes(r, customerHandler)
			// Rutas de properties
			RegisterPropertyRoutes(r, propertyHandler)
			// Rutas de trabajos
			jobHandler.RegisterRoutes(r)
			// Rutas de catálogos de trabajos
			jobCategoryHandler.RegisterRoutes(r)
			jobStatusHandler.RegisterRoutes(r)
			jobPriorityHandler.RegisterRoutes(r)
			techJobStatusHandler.RegisterRoutes(r)
			taskStatusHandler.RegisterRoutes(r)
			// Rutas de cotizaciones
			quoteHandler.RegisterRoutes(r)
			quoteStatusHandler.RegisterRoutes(r)
			// Rutas de supervisores
			supervisorHandler.RegisterRoutes(r)
			// Rutas de equipos de propiedad
			propEquipHandler.RegisterRoutes(r)
			// Rutas de equipos de trabajo
			jobEquipHandler.RegisterRoutes(r)
			// Rutas de facturas y pagos
			invoiceHandler.RegisterRoutes(r)
			invoicePaymentHandler.RegisterRoutes(r)
			// Rutas de garantías y catálogos
			warrantyTypeHandler.RegisterRoutes(r)
			warrantyStatusHandler.RegisterRoutes(r)
			warrantyClaimTypeHandler.RegisterRoutes(r)
			warrantyClaimStatusHandler.RegisterRoutes(r)
			warrantyHandler.RegisterRoutes(r)
			warrantyEquipHandler.RegisterRoutes(r)
			warrantyClaimHandler.RegisterRoutes(r)
			// Rutas de visitas de trabajo y archivos
			jobVisitHandler.RegisterRoutes(r)

			// Module 13: Activities and Communications
			// Job Activity Logs
			r.Route("/jobs/{jobId}/activities", func(r chi.Router) {
				r.Post("/", jobActivityLogHandler.Create)
				r.Get("/", jobActivityLogHandler.List)
				r.Delete("/{id}", jobActivityLogHandler.Delete)
			})

			// Job Emails
			r.Route("/jobs/{jobId}/emails", func(r chi.Router) {
				r.Post("/", jobEmailHandler.Create)
				r.Get("/", jobEmailHandler.List)
				r.Delete("/{id}", jobEmailHandler.Delete)
			})

			// Job Residents
			r.Route("/jobs/{jobId}/residents", func(r chi.Router) {
				r.Post("/", jobResidentHandler.Create)
				r.Get("/", jobResidentHandler.List)
				r.Put("/{id}", jobResidentHandler.Update)
				r.Delete("/{id}", jobResidentHandler.Delete)
			})

			// Job Rate Statuses (catalog)
			r.Route("/job-rate-statuses", func(r chi.Router) {
				r.Post("/", jobRateStatusHandler.Create)
				r.Get("/", jobRateStatusHandler.List)
				r.Get("/{id}", jobRateStatusHandler.GetByID)
				r.Put("/{id}", jobRateStatusHandler.Update)
				r.Delete("/{id}", jobRateStatusHandler.Delete)
			})

			// SMS Templates
			smsTemplateHandler.RegisterRoutes(r)

			// Job Tasks
			r.Route("/jobs/{jobId}/tasks", func(r chi.Router) {
				r.Post("/", jobTaskHandler.Create)
				r.Get("/", jobTaskHandler.List)
				r.Put("/{id}", jobTaskHandler.Update)
				r.Delete("/{id}", jobTaskHandler.Delete)
			})

			// Global tasks view
			r.Get("/tasks", jobTaskHandler.ListAll)

			// Task notification email
			r.Post("/tasks/{id}/notification", jobTaskHandler.SendNotification)

			// Job Rates
			r.Route("/jobs/{jobId}/rates", func(r chi.Router) {
				r.Post("/", jobRateHandler.Create)
				r.Get("/", jobRateHandler.List)
				r.Put("/{id}", jobRateHandler.Update)
				r.Delete("/{id}", jobRateHandler.Delete)
			})

			// Job SMS
			r.Route("/jobs/{jobId}/sms", func(r chi.Router) {
				r.Post("/", jobSMSHandler.Create)
				r.Get("/", jobSMSHandler.List)
				r.Delete("/{id}", jobSMSHandler.Delete)
			})

			// Calculate rate payment (standalone endpoint)
			r.Post("/calculate-rate-payment", jobRateHandler.CalculatePayment)

			// Module 15: Alerts
			alertHandler.RegisterRoutes(r)

			// Dashboard
			dashboardHandler.RegisterRoutes(r)

			// New Enhanced Dashboard
			newDashboardHdlr.RegisterRoutes(r)
		})
	})
	return r
}
