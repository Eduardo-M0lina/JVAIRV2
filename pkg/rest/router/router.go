package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/your-org/jvairv2/docs" // Importación de documentación Swagger
	"github.com/your-org/jvairv2/pkg/rest/handler"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

// New crea un nuevo router HTTP con las rutas configuradas
func New(healthHandler *handler.HealthHandler, authHandler *authHandler.Handler, authMiddleware *middleware.AuthMiddleware) *chi.Mux {
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

		// API v1
		r.Route("/api/v1", func(r chi.Router) {
			// Aquí se agregarán las rutas de la API
			// r.Mount("/customers", customerHandler.Routes())
			// r.Mount("/jobs", jobHandler.Routes())
			// r.Mount("/properties", propertyHandler.Routes())
		})
	})

	return r
}
