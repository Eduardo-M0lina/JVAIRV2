package dashboard

import "context"

// DashboardFilters representa los filtros disponibles para el dashboard
type DashboardFilters struct {
	UserID        *int64  // Filtrar por técnico (nil = todos, -1 = unassigned)
	JobStatusID   *int64  // Filtrar por estado
	JobPriorityID *int64  // Filtrar por prioridad
	IsClosed      *string // "0" = abiertos, "1" = cerrados, "all" = todos
	Year          *int    // Filtrar por año
	Week          *int    // Filtrar por semana
	LastDays      *int    // Filtrar jobs creados en los últimos X días (ej: 7, 30, 90)
	Search        *string // Búsqueda de texto libre
	Sort          *string // Campo para ordenar (ej: "work_order", "property.city")
	Direction     *string // Dirección: "asc" o "desc"
	Page          int     // Número de página (1-indexed)
	PageSize      int     // Tamaño de página
}

// DashboardJobsResult representa el resultado paginado de jobs
type DashboardJobsResult struct {
	Jobs       []DashboardJob `json:"jobs"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// Repository define las operaciones de acceso a datos para el dashboard
type Repository interface {
	// GetStats obtiene estadísticas generales del dashboard aplicando los filtros
	// Si userID es nil, retorna stats globales (admin)
	// Si userID tiene valor, retorna stats filtradas por ese usuario (técnico)
	// Los filtros adicionales se aplican para que las stats sean consistentes con los jobs mostrados
	GetStats(ctx context.Context, userID *int64, filters DashboardFilters) (*DashboardStats, error)

	// GetJobs obtiene jobs con filtros, ordenamiento, búsqueda y paginación
	GetJobs(ctx context.Context, baseFilters map[string]interface{}, filters DashboardFilters) (*DashboardJobsResult, error)
}
