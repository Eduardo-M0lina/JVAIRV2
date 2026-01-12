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
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
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
		})
	})
	return r
}
