package router

import (
	"github.com/go-chi/chi/v5"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
)

// RegisterAssignedRoleRoutes registra las rutas de assigned-roles
func RegisterAssignedRoleRoutes(r chi.Router, handler *assignedRoleHandler.Handler) {
	r.Route("/assigned-roles", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Assign)
		r.Get("/{id}", handler.Get)
		r.Get("/entity/{entityType}/{entityId}", handler.GetByEntity)
		r.Get("/check/{roleId}/{entityType}/{entityId}", handler.HasRole)
		r.Delete("/revoke/{roleId}/{entityType}/{entityId}", handler.Revoke)
	})
}
