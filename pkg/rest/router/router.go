package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/your-org/jvairv2/docs" // Importación de documentación Swagger
	"github.com/your-org/jvairv2/pkg/domain/user"
	"github.com/your-org/jvairv2/pkg/rest/handler"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	customerHandler "github.com/your-org/jvairv2/pkg/rest/handler/customer"
	jobCategoryHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_category"
	jobPriorityHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_priority"
	jobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_status"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	propertyHandler "github.com/your-org/jvairv2/pkg/rest/handler/property"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	settingsHandler "github.com/your-org/jvairv2/pkg/rest/handler/settings"
	taskStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/task_status"
	techJobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/technician_job_status"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
	workflowHandler "github.com/your-org/jvairv2/pkg/rest/handler/workflow"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

// New crea un nuevo router HTTP con las rutas configuradas
func New(
	healthHandler *handler.HealthHandler,
	authHandler *authHandler.Handler,
	userHandler *userHandler.Handler,
	roleHandler *roleHandler.Handler,
	abilityHandler *abilityHandler.Handler,
	assignedRoleHandler *assignedRoleHandler.Handler,
	permissionHandler *permissionHandler.Handler,
	settingsHandler *settingsHandler.Handler,
	workflowHandler *workflowHandler.Handler,
	customerHandler *customerHandler.Handler,
	propertyHandler *propertyHandler.Handler,
	jobCategoryHandler *jobCategoryHandler.Handler,
	jobStatusHandler *jobStatusHandler.Handler,
	jobPriorityHandler *jobPriorityHandler.Handler,
	techJobStatusHandler *techJobStatusHandler.Handler,
	taskStatusHandler *taskStatusHandler.Handler,
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
		RegisterAuthRoutes(r, authHandler)
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
			// Rutas de workflows
			SetupWorkflowRoutes(r, workflowHandler, authMiddleware)
			// Rutas de customers
			RegisterCustomerRoutes(r, customerHandler)
			// Rutas de properties
			RegisterPropertyRoutes(r, propertyHandler)
			// Rutas de catálogos de trabajos
			jobCategoryHandler.RegisterRoutes(r)
			jobStatusHandler.RegisterRoutes(r)
			jobPriorityHandler.RegisterRoutes(r)
			techJobStatusHandler.RegisterRoutes(r)
			taskStatusHandler.RegisterRoutes(r)
		})
	})
	return r
}
