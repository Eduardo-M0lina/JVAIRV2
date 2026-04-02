package new_dashboard

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// UseCase implementa la lógica de negocio del dashboard enriquecido
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// GetAdminEnhancedDashboard retorna el dashboard enriquecido para administradores
// Ejecuta todas las queries en paralelo usando goroutines
func (uc *UseCase) GetAdminEnhancedDashboard(ctx context.Context, timeRange TimeRange) (*AdminEnhancedDashboard, error) {
	// Timeout de 10 segundos para todo el dashboard
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := []error{}

	dashboard := &AdminEnhancedDashboard{}

	// Ejecutar todas las queries en paralelo
	wg.Add(11)

	// 1. Enhanced Stats
	go func() {
		defer wg.Done()
		stats, err := uc.repo.GetEnhancedStats(ctx, nil, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get enhanced stats", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.Stats = *stats
		mu.Unlock()
	}()

	// 2. Alert Summary
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetAlertSummary(ctx, nil, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get alert summary", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.AlertSummary = summary
		mu.Unlock()
	}()

	// 3. Task Summary
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetTaskSummary(ctx, nil, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get task summary", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.TaskSummary = summary
		mu.Unlock()
	}()

	// 4. Recent Activity
	go func() {
		defer wg.Done()
		activities, err := uc.repo.GetRecentActivity(ctx, nil, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get recent activity", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.RecentActivity = activities
		mu.Unlock()
	}()

	// 5. Invoice Summary
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetInvoiceSummary(ctx, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get invoice summary", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.InvoiceSummary = summary
		mu.Unlock()
	}()

	// 6. Quote Summary
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetQuoteSummary(ctx, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get quote summary", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.QuoteSummary = summary
		mu.Unlock()
	}()

	// 7. Warranty Summary
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetWarrantySummary(ctx, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get warranty summary", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.WarrantySummary = summary
		mu.Unlock()
	}()

	// 8. Jobs by Category
	go func() {
		defer wg.Done()
		categories, err := uc.repo.GetJobsByCategory(ctx, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get jobs by category", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.JobsByCategory = categories
		mu.Unlock()
	}()

	// 9. Jobs by Status
	go func() {
		defer wg.Done()
		statuses, err := uc.repo.GetJobsByStatus(ctx, nil, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get jobs by status", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.JobsByStatus = statuses
		mu.Unlock()
	}()

	// 10. Jobs Due This Week
	go func() {
		defer wg.Done()
		jobs, err := uc.repo.GetJobsDueThisWeek(ctx, nil)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get jobs due this week", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.JobsDueThisWeek = jobs
		mu.Unlock()
	}()

	// 11. Technician Workload
	go func() {
		defer wg.Done()
		workloads, err := uc.repo.GetTechnicianWorkload(ctx, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get technician workload", slog.String("error", err.Error()))
			return
		}
		mu.Lock()
		dashboard.TechnicianWorkload = workloads
		mu.Unlock()
	}()

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Si hubo errores críticos (más del 50% de nodos fallaron), retornar error
	if len(errors) > 5 {
		slog.ErrorContext(ctx, "Too many errors in dashboard", slog.Int("errorCount", len(errors)))
		return nil, errors[0]
	}

	// Si algunos nodos fallaron pero la mayoría funcionó, retornar dashboard parcial
	if len(errors) > 0 {
		slog.WarnContext(ctx, "Some dashboard nodes failed", slog.Int("errorCount", len(errors)))
	}

	return dashboard, nil
}

// GetTechnicianEnhancedDashboard retorna el dashboard enriquecido para técnicos
// Todos los datos están filtrados por el user_id del técnico
func (uc *UseCase) GetTechnicianEnhancedDashboard(ctx context.Context, userID int64, timeRange TimeRange) (*TechnicianEnhancedDashboard, error) {
	// Timeout de 10 segundos para todo el dashboard
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := []error{}

	dashboard := &TechnicianEnhancedDashboard{}

	// Ejecutar queries en paralelo (solo 6 nodos para técnicos)
	wg.Add(6)

	// 1. Enhanced Stats (filtrado por user_id)
	go func() {
		defer wg.Done()
		stats, err := uc.repo.GetEnhancedStats(ctx, &userID, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get enhanced stats", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.Stats = *stats
		mu.Unlock()
	}()

	// 2. Alert Summary (filtrado por user_id)
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetAlertSummary(ctx, &userID, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get alert summary", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.AlertSummary = summary
		mu.Unlock()
	}()

	// 3. Task Summary (filtrado por user_id)
	go func() {
		defer wg.Done()
		summary, err := uc.repo.GetTaskSummary(ctx, &userID, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get task summary", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.TaskSummary = summary
		mu.Unlock()
	}()

	// 4. Recent Activity (filtrado por user_id - solo actividad de sus jobs)
	go func() {
		defer wg.Done()
		activities, err := uc.repo.GetRecentActivity(ctx, &userID, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get recent activity", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.RecentActivity = activities
		mu.Unlock()
	}()

	// 5. Jobs by Status (filtrado por user_id)
	go func() {
		defer wg.Done()
		statuses, err := uc.repo.GetJobsByStatus(ctx, &userID, timeRange)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get jobs by status", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.JobsByStatus = statuses
		mu.Unlock()
	}()

	// 6. Jobs Due This Week (filtrado por user_id)
	go func() {
		defer wg.Done()
		jobs, err := uc.repo.GetJobsDueThisWeek(ctx, &userID)
		if err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
			slog.ErrorContext(ctx, "Failed to get jobs due this week", slog.String("error", err.Error()), slog.Int64("userId", userID))
			return
		}
		mu.Lock()
		dashboard.JobsDueThisWeek = jobs
		mu.Unlock()
	}()

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Si hubo errores críticos (más del 50% de nodos fallaron), retornar error
	if len(errors) > 3 {
		slog.ErrorContext(ctx, "Too many errors in technician dashboard", slog.Int("errorCount", len(errors)), slog.Int64("userId", userID))
		return nil, errors[0]
	}

	// Si algunos nodos fallaron pero la mayoría funcionó, retornar dashboard parcial
	if len(errors) > 0 {
		slog.WarnContext(ctx, "Some dashboard nodes failed", slog.Int("errorCount", len(errors)), slog.Int64("userId", userID))
	}

	return dashboard, nil
}
