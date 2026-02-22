package job

import (
	"context"
	"log/slog"
	"time"
)

// Create crea un nuevo job
func (uc *UseCase) Create(ctx context.Context, j *Job) error {
	if err := j.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que la categoría existe
	if _, err := uc.jobCategoryRepo.GetByID(ctx, j.JobCategoryID); err != nil {
		slog.ErrorContext(ctx, "Invalid job category",
			slog.Int64("jobCategoryId", j.JobCategoryID),
			slog.String("error", err.Error()))
		return ErrInvalidJobCategory
	}

	// Verificar que la prioridad existe
	if _, err := uc.jobPriorityRepo.GetByID(ctx, j.JobPriorityID); err != nil {
		slog.ErrorContext(ctx, "Invalid job priority",
			slog.Int64("jobPriorityId", j.JobPriorityID),
			slog.String("error", err.Error()))
		return ErrInvalidJobPriority
	}

	// Verificar que la propiedad existe
	if _, err := uc.propertyRepo.GetByID(ctx, j.PropertyID); err != nil {
		slog.ErrorContext(ctx, "Invalid property",
			slog.Int64("propertyId", j.PropertyID),
			slog.String("error", err.Error()))
		return ErrInvalidProperty
	}

	// Obtener workflow_id desde la propiedad (customer -> workflow)
	workflowID, err := uc.propertyRepo.GetWorkflowID(ctx, j.PropertyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workflow for property",
			slog.Int64("propertyId", j.PropertyID),
			slog.String("error", err.Error()))
		return ErrInvalidWorkflow
	}
	j.WorkflowID = workflowID

	// Obtener el status inicial del workflow
	initialStatusID, err := uc.workflowRepo.GetInitialStatusID(ctx, j.WorkflowID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get initial status for workflow",
			slog.Int64("workflowId", j.WorkflowID),
			slog.String("error", err.Error()))
		return err
	}
	j.JobStatusID = initialStatusID

	// Verificar usuario si se proporcionó
	if j.UserID != nil && *j.UserID > 0 {
		if _, err := uc.userRepo.GetByID(ctx, *j.UserID); err != nil {
			slog.ErrorContext(ctx, "Invalid user",
				slog.Int64("userId", *j.UserID),
				slog.String("error", err.Error()))
			return ErrInvalidUser
		}
	}

	// Establecer fecha de recepción si no se proporcionó
	if j.DateReceived.IsZero() {
		j.DateReceived = time.Now()
	}

	if err := uc.repo.Create(ctx, j); err != nil {
		slog.ErrorContext(ctx, "Failed to create job",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Job created successfully",
		slog.Int64("id", j.ID))

	return nil
}
