package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/handler/workflow"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

// SetupWorkflowRoutes configura las rutas para el módulo de workflows
func SetupWorkflowRoutes(r chi.Router, handler *workflow.Handler, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/workflows", func(r chi.Router) {
		// Aplicar middleware de autenticación
		r.Use(authMiddleware.Authenticate)

		// Rutas de workflows
		r.Get("/", handler.List)                     // GET /api/v1/workflows
		r.Post("/", handler.Create)                  // POST /api/v1/workflows
		r.Get("/{id}", handler.Get)                  // GET /api/v1/workflows/{id}
		r.Put("/{id}", handler.Update)               // PUT /api/v1/workflows/{id}
		r.Delete("/{id}", handler.Delete)            // DELETE /api/v1/workflows/{id}
		r.Post("/{id}/duplicate", handler.Duplicate) // POST /api/v1/workflows/{id}/duplicate
	})
}
