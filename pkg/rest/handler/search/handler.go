package search

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainSearch "github.com/your-org/jvairv2/pkg/domain/search"
	domainUser "github.com/your-org/jvairv2/pkg/domain/user"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las peticiones HTTP para búsqueda global
type Handler struct {
	useCase domainSearch.Service
}

// NewHandler crea una nueva instancia del handler de búsqueda
func NewHandler(useCase domainSearch.Service) *Handler {
	return &Handler{useCase: useCase}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/search", h.GlobalSearch)
}

// GlobalSearch realiza una búsqueda global en todas las entidades
// @Summary Búsqueda global
// @Description Busca en múltiples entidades (jobs, customers, properties, invoices, quotes, warranties, warranty claims, users) simultáneamente. Respeta el permiso job_view_user_only para filtrar resultados.
// @Tags Search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Término de búsqueda"
// @Param limit query int false "Límite de resultados por entidad (default: 10, max: 50)"
// @Success 200 {object} search.GlobalSearchResponse "Resultados de búsqueda agrupados por entidad"
// @Failure 400 {object} map[string]string "Parámetro de búsqueda requerido"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/search [get]
func (h *Handler) GlobalSearch(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		response.Error(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener query de búsqueda
	query := r.URL.Query().Get("q")
	if query == "" {
		response.Error(w, http.StatusBadRequest, "El parámetro 'q' es requerido")
		return
	}

	// Parsear límite
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 50 {
				limit = 50
			}
		}
	}

	// Verificar si el usuario tiene restricción job_view_user_only
	// Nota: Si el usuario tiene "*" (superadmin), NO debe tener esta restricción
	userOnly := middleware.HasAbilityExact(r.Context(), "job_view_user_only")

	slog.Info("Global search request",
		slog.String("query", query),
		slog.Int("limit", limit),
		slog.Int64("userId", userCtx.ID),
		slog.Bool("userOnly", userOnly))

	filters := domainSearch.SearchFilters{
		Query:    query,
		Limit:    limit,
		UserID:   &userCtx.ID,
		UserOnly: userOnly,
	}

	results, err := h.useCase.GlobalSearch(r.Context(), filters)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to perform global search",
			slog.String("error", err.Error()),
			slog.String("query", query),
			slog.Int64("userId", userCtx.ID))
		response.Error(w, http.StatusInternalServerError, "Error al realizar la búsqueda")
		return
	}

	response.JSON(w, http.StatusOK, results)
}
