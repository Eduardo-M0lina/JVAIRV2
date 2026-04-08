package dashboard

import (
	"log/slog"
	"net/http"
	"strconv"

	domainDashboard "github.com/angumol/jvairv2/pkg/domain/dashboard"
	domainUser "github.com/angumol/jvairv2/pkg/domain/user"
	"github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/angumol/jvairv2/pkg/rest/middleware"
	"github.com/go-chi/chi/v5"
)

// adminRoles define los roles que ven el dashboard de administrador
var adminRoles = map[string]bool{
	"administrator": true,
	"executive":     true,
}

// Handler maneja las peticiones HTTP para el dashboard
type Handler struct {
	useCase domainDashboard.Service
}

// NewHandler crea una nueva instancia del handler de dashboard
func NewHandler(useCase domainDashboard.Service) *Handler {
	return &Handler{useCase: useCase}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/", h.GetDashboard)
	})
}

// GetDashboard retorna el dashboard según el rol del usuario autenticado
// @Summary Obtener dashboard
// @Description Retorna el dashboard con jobs según el rol del usuario. Admin ve jobs awaiting dispatch y urgent jobs. Técnicos ven jobs dispatched a ellos y sus urgent jobs. Soporta filtros, ordenamiento, búsqueda y paginación.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página (default: 1)"
// @Param pageSize query int false "Tamaño de página (default: 20)"
// @Param limit query int false "Alias de pageSize para backward compatibility"
// @Param filters[user_id] query string false "Filtrar por técnico (ID numérico o 'unassigned')"
// @Param filters[job_status_id] query int false "Filtrar por estado del job"
// @Param filters[job_priority_id] query int false "Filtrar por prioridad del job"
// @Param filters[is_closed] query string false "Filtrar por estado abierto/cerrado ('0', '1', 'all')"
// @Param filters[year] query int false "Filtrar por año"
// @Param filters[week] query int false "Filtrar por número de semana"
// @Param filters[last_days] query int false "Filtrar jobs creados en los últimos X días (ej: 7, 30, 90)"
// @Param search query string false "Búsqueda de texto libre en work_order, property, customer, etc."
// @Param sort query string false "Campo para ordenar (work_order, date_received, dispatch_date, due_date, completion_date, week_number, job_sales_price, property.city, property.zip, property.customer.name, user_id, status, priority.order)"
// @Param direction query string false "Dirección del ordenamiento ('asc' o 'desc', default: 'desc')"
// @Success 200 {object} dashboard.AdminDashboard "Dashboard de administrador"
// @Success 200 {object} dashboard.TechnicianDashboard "Dashboard de técnico"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/dashboard [get]
func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		handler.RespondWithError(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Parsear filtros desde query params
	filters := domainDashboard.DashboardFilters{
		Page:     1,
		PageSize: 20,
	}

	// Paginación
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			filters.Page = p
		}
	}
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			filters.PageSize = ps
		}
	}
	// Backward compatibility con "limit"
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			filters.PageSize = l
		}
	}

	// Filtros
	if userIDStr := r.URL.Query().Get("filters[user_id]"); userIDStr != "" {
		if userIDStr == "unassigned" {
			uid := int64(-1)
			filters.UserID = &uid
		} else if uid, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filters.UserID = &uid
		}
	}

	if statusIDStr := r.URL.Query().Get("filters[job_status_id]"); statusIDStr != "" {
		if sid, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			filters.JobStatusID = &sid
		}
	}

	if priorityIDStr := r.URL.Query().Get("filters[job_priority_id]"); priorityIDStr != "" {
		if pid, err := strconv.ParseInt(priorityIDStr, 10, 64); err == nil {
			filters.JobPriorityID = &pid
		}
	}

	if isClosedStr := r.URL.Query().Get("filters[is_closed]"); isClosedStr != "" {
		filters.IsClosed = &isClosedStr
	}

	if yearStr := r.URL.Query().Get("filters[year]"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			filters.Year = &y
		}
	}

	if weekStr := r.URL.Query().Get("filters[week]"); weekStr != "" {
		if w, err := strconv.Atoi(weekStr); err == nil {
			filters.Week = &w
		}
	}

	if lastDaysStr := r.URL.Query().Get("filters[last_days]"); lastDaysStr != "" {
		if ld, err := strconv.Atoi(lastDaysStr); err == nil && ld > 0 {
			filters.LastDays = &ld
		}
	}

	// Búsqueda
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = &search
	}

	// Ordenamiento
	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters.Sort = &sort
	}
	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters.Direction = &direction
	}

	// Determinar si el usuario es admin
	isAdmin := false
	if userCtx.RoleName != nil {
		isAdmin = adminRoles[*userCtx.RoleName]
	}

	if isAdmin {
		// Vista admin
		dashboard, err := h.useCase.GetAdminDashboard(r.Context(), filters)
		if err != nil {
			slog.ErrorContext(r.Context(), "Failed to get admin dashboard",
				slog.String("error", err.Error()),
				slog.Int64("userId", userCtx.ID))
			handler.RespondWithError(w, http.StatusInternalServerError, "Error al obtener dashboard")
			return
		}
		handler.RespondWithJSON(w, http.StatusOK, dashboard)
		return
	}

	// Vista técnico
	dashboard, err := h.useCase.GetTechnicianDashboard(r.Context(), userCtx.ID, filters)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get technician dashboard",
			slog.String("error", err.Error()),
			slog.Int64("userId", userCtx.ID))
		handler.RespondWithError(w, http.StatusInternalServerError, "Error al obtener dashboard")
		return
	}
	handler.RespondWithJSON(w, http.StatusOK, dashboard)
}
