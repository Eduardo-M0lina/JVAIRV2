package dashboard

import (
	"context"
	"log/slog"
)

// Service define la interfaz del servicio de dashboard
type Service interface {
	GetAdminDashboard(ctx context.Context, filters DashboardFilters) (*AdminDashboard, error)
	GetTechnicianDashboard(ctx context.Context, userID int64, filters DashboardFilters) (*TechnicianDashboard, error)
}

// UseCase implementa la lógica de negocio del dashboard
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso de dashboard
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// GetAdminDashboard retorna el dashboard para administradores con filtros
func (uc *UseCase) GetAdminDashboard(ctx context.Context, filters DashboardFilters) (*AdminDashboard, error) {
	// Obtener estadísticas globales aplicando los mismos filtros
	stats, err := uc.repo.GetStats(ctx, nil, filters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get admin stats", slog.String("error", err.Error()))
		return nil, err
	}

	// Obtener jobs awaiting dispatch (job_status_id = 1)
	baseFilters := map[string]interface{}{
		"job_status_id": int64(1),
		"closed":        false,
	}
	jobsAwaitingDispatch, err := uc.repo.GetJobs(ctx, baseFilters, filters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get jobs awaiting dispatch", slog.String("error", err.Error()))
		return nil, err
	}

	// Obtener jobs urgentes (admin view: job_status_id = 9, job_priority_id = 4)
	urgentFilters := filters
	urgentFilters.JobStatusID = nil // Reset para usar el base filter
	urgentFilters.JobPriorityID = nil
	urgentBaseFilters := map[string]interface{}{
		"job_status_id":   int64(9),
		"job_priority_id": int64(4),
	}
	jobsUrgent, err := uc.repo.GetJobs(ctx, urgentBaseFilters, urgentFilters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get urgent jobs", slog.String("error", err.Error()))
		return nil, err
	}

	return &AdminDashboard{
		Stats:                *stats,
		JobsAwaitingDispatch: jobsAwaitingDispatch.Jobs,
		JobsUrgent:           jobsUrgent.Jobs,
	}, nil
}

// GetTechnicianDashboard retorna el dashboard para técnicos con filtros
func (uc *UseCase) GetTechnicianDashboard(ctx context.Context, userID int64, filters DashboardFilters) (*TechnicianDashboard, error) {
	// Obtener estadísticas filtradas por usuario aplicando los mismos filtros
	stats, err := uc.repo.GetStats(ctx, &userID, filters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get technician stats",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		return nil, err
	}

	// Obtener jobs dispatched al técnico (closed = 0, job_status_id = 2, user_id = userID)
	baseFilters := map[string]interface{}{
		"closed":        false,
		"job_status_id": int64(2),
		"user_id":       userID,
	}
	jobsDispatched, err := uc.repo.GetJobs(ctx, baseFilters, filters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get dispatched jobs",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		return nil, err
	}

	// Obtener jobs urgentes del técnico (closed = 0, job_status_id = 2, job_priority_id = 4, user_id = userID)
	urgentFilters := filters
	urgentFilters.JobStatusID = nil
	urgentFilters.JobPriorityID = nil
	urgentBaseFilters := map[string]interface{}{
		"closed":          false,
		"job_status_id":   int64(2),
		"job_priority_id": int64(4),
		"user_id":         userID,
	}
	jobsUrgent, err := uc.repo.GetJobs(ctx, urgentBaseFilters, urgentFilters)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get urgent jobs",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		return nil, err
	}

	return &TechnicianDashboard{
		Stats:          *stats,
		JobsDispatched: jobsDispatched.Jobs,
		JobsUrgent:     jobsUrgent.Jobs,
	}, nil
}
