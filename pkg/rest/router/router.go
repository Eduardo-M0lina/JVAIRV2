package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/your-org/jvairv2/pkg/rest/handler"
)

// New crea un nuevo router HTTP con las rutas configuradas
func New(healthHandler *handler.HealthHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Rutas públicas
	r.Group(func(r chi.Router) {
		// Health check
		r.Get("/health", healthHandler.Check)
	})

	// Rutas protegidas (para implementar más adelante)
	r.Group(func(r chi.Router) {
		// Aquí se agregarán middlewares de autenticación
		// r.Use(authMiddleware)

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
