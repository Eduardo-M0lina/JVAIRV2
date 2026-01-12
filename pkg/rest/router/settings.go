package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/handler/settings"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

// SetupSettingsRoutes configura las rutas para el módulo de configuraciones
func SetupSettingsRoutes(r chi.Router, handler *settings.Handler, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/settings", func(r chi.Router) {
		// Aplicar middleware de autenticación
		r.Use(authMiddleware.Authenticate)

		// Rutas de configuraciones
		r.Get("/", handler.Get)      // GET /api/v1/settings
		r.Put("/", handler.Update)   // PUT /api/v1/settings
		r.Patch("/", handler.Update) // PATCH /api/v1/settings (alias de PUT)
	})
}
