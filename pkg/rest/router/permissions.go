package router

import (
	"github.com/go-chi/chi/v5"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
)

// RegisterPermissionRoutes registra las rutas de permisos
func RegisterPermissionRoutes(r chi.Router, handler *permissionHandler.Handler) {
	r.Route("/permissions", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/check/{abilityId}/{entityType}/{entityId}", handler.List) // TODO: Implementar método específico para verificar permisos
	})
}
