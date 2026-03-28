package new_dashboard

import "context"

// Repository define las operaciones de acceso a datos para el dashboard enriquecido
type Repository interface {
	// GetEnhancedStats obtiene estadísticas expandidas del dashboard
	// Si userID es nil, retorna stats globales (admin)
	// Si userID tiene valor, retorna stats filtradas por ese usuario (técnico)
	GetEnhancedStats(ctx context.Context, userID *int64, timeRange TimeRange) (*EnhancedStats, error)

	// GetAlertSummary obtiene el resumen de alertas
	// Si userID es nil, retorna todas las alertas (admin)
	// Si userID tiene valor, retorna solo alertas del usuario (técnico)
	GetAlertSummary(ctx context.Context, userID *int64, timeRange TimeRange) (*AlertSummary, error)

	// GetTaskSummary obtiene el resumen de tareas
	// Si userID es nil, retorna todas las tareas (admin)
	// Si userID tiene valor, retorna solo tareas del usuario (técnico)
	GetTaskSummary(ctx context.Context, userID *int64, timeRange TimeRange) (*TaskSummary, error)

	// GetRecentActivity obtiene actividad reciente
	// Si userID es nil, retorna toda la actividad (admin)
	// Si userID tiene valor, retorna solo actividad de jobs del usuario (técnico)
	GetRecentActivity(ctx context.Context, userID *int64, timeRange TimeRange) ([]*Activity, error)

	// GetInvoiceSummary obtiene el resumen de facturación (solo admin)
	GetInvoiceSummary(ctx context.Context, timeRange TimeRange) (*InvoiceSummary, error)

	// GetQuoteSummary obtiene el resumen de cotizaciones (solo admin)
	GetQuoteSummary(ctx context.Context, timeRange TimeRange) (*QuoteSummary, error)

	// GetWarrantySummary obtiene el resumen de garantías (solo admin)
	GetWarrantySummary(ctx context.Context, timeRange TimeRange) (*WarrantySummary, error)

	// GetJobsByCategory obtiene la distribución de jobs por categoría (solo admin)
	GetJobsByCategory(ctx context.Context, timeRange TimeRange) ([]*CategoryCount, error)

	// GetJobsByStatus obtiene la distribución de jobs por estado
	// Si userID es nil, retorna todos los jobs (admin)
	// Si userID tiene valor, retorna solo jobs del usuario (técnico)
	GetJobsByStatus(ctx context.Context, userID *int64, timeRange TimeRange) ([]*StatusCount, error)

	// GetJobsDueThisWeek obtiene jobs que vencen esta semana
	// Si userID es nil, retorna todos los jobs (admin)
	// Si userID tiene valor, retorna solo jobs del usuario (técnico)
	GetJobsDueThisWeek(ctx context.Context, userID *int64) ([]*DueJob, error)

	// GetTechnicianWorkload obtiene la carga de trabajo por técnico (solo admin)
	GetTechnicianWorkload(ctx context.Context, timeRange TimeRange) ([]*TechnicianLoad, error)
}
