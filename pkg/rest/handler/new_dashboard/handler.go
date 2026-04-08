package new_dashboard

import (
	"log/slog"
	"net/http"

	domainNewDashboard "github.com/angumol/jvairv2/pkg/domain/new_dashboard"
	domainUser "github.com/angumol/jvairv2/pkg/domain/user"
	"github.com/angumol/jvairv2/pkg/rest/middleware"
	"github.com/angumol/jvairv2/pkg/rest/response"

	"github.com/go-chi/chi/v5"
)

// Handler maneja las peticiones HTTP del dashboard enriquecido
type Handler struct {
	useCase *domainNewDashboard.UseCase
}

// NewHandler crea una nueva instancia del handler
func NewHandler(useCase *domainNewDashboard.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Roles que tienen acceso completo (admin)
var adminRoles = map[string]bool{
	"administrator": true,
	"executive":     true,
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/new-dashboard", func(r chi.Router) {
		r.Get("/", h.GetEnhancedDashboard)
	})
}

// GetEnhancedDashboard retorna el dashboard enriquecido según el rol del usuario
// @Summary Obtener dashboard enriquecido
// @Description Retorna un dashboard enriquecido con datos de múltiples módulos (invoices, quotes, tasks, warranties, alerts, activity). Admin ve todos los nodos, técnicos ven solo sus datos filtrados por user_id. Todas las queries se ejecutan en paralelo.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param range query string false "Rango de tiempo: 7days, 30days (default), 90days, thisMonth, lastMonth, thisYear"
// @Success 200 {object} new_dashboard.AdminEnhancedDashboard "Dashboard enriquecido de administrador"
// @Success 200 {object} new_dashboard.TechnicianEnhancedDashboard "Dashboard enriquecido de técnico"
// @Failure 400 {object} map[string]string "Rango de tiempo inválido"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/new-dashboard [get]
func (h *Handler) GetEnhancedDashboard(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		response.Error(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Parsear rango de tiempo (default: 30days)
	rangeParam := r.URL.Query().Get("range")
	if rangeParam == "" {
		rangeParam = "30days"
	}

	// Validar rango de tiempo
	timeRange := domainNewDashboard.TimeRange(rangeParam)
	validRanges := map[domainNewDashboard.TimeRange]bool{
		domainNewDashboard.TimeRange7Days:     true,
		domainNewDashboard.TimeRange30Days:    true,
		domainNewDashboard.TimeRange90Days:    true,
		domainNewDashboard.TimeRangeThisMonth: true,
		domainNewDashboard.TimeRangeLastMonth: true,
		domainNewDashboard.TimeRangeThisYear:  true,
	}

	if !validRanges[timeRange] {
		response.Error(w, http.StatusBadRequest, "Rango de tiempo inválido. Valores permitidos: 7days, 30days, 90days, thisMonth, lastMonth, thisYear")
		return
	}

	// Determinar si el usuario es admin
	isAdmin := false
	if userCtx.RoleName != nil {
		isAdmin = adminRoles[*userCtx.RoleName]
	}

	if isAdmin {
		// Vista admin - todos los nodos con datos globales
		dashboard, err := h.useCase.GetAdminEnhancedDashboard(r.Context(), timeRange)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to get admin enhanced dashboard",
				slog.String("error", err.Error()),
				slog.Int64("userId", userCtx.ID),
				slog.String("range", string(timeRange)))
			response.Error(w, http.StatusInternalServerError, "Error al obtener dashboard")
			return
		}
		response.JSON(w, http.StatusOK, dashboard)
		return
	}

	// Vista técnico - solo 6 nodos filtrados por user_id
	dashboard, err := h.useCase.GetTechnicianEnhancedDashboard(r.Context(), userCtx.ID, timeRange)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get technician enhanced dashboard",
			slog.String("error", err.Error()),
			slog.Int64("userId", userCtx.ID),
			slog.String("range", string(timeRange)))
		response.Error(w, http.StatusInternalServerError, "Error al obtener dashboard")
		return
	}
	response.JSON(w, http.StatusOK, dashboard)
}
