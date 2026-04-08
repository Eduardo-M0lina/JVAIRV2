package job

import (
	"context"
	"log/slog"
)

// Update actualiza un job existente
func (uc *UseCase) Update(ctx context.Context, j *Job) error {
	if err := j.ValidateUpdate(); err != nil {
		return err
	}

	// Verificar que el job existe
	existing, err := uc.repo.GetByID(ctx, j.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Job not found for update",
			slog.Int64("id", j.ID),
			slog.String("error", err.Error()))
		return ErrJobNotFound
	}

	if existing.IsDeleted() {
		return ErrJobNotFound
	}

	// Verificar categoría si cambió
	if j.JobCategoryID > 0 && j.JobCategoryID != existing.JobCategoryID {
		if _, err := uc.jobCategoryRepo.GetByID(ctx, j.JobCategoryID); err != nil {
			return ErrInvalidJobCategory
		}
	}

	// Verificar prioridad si cambió
	if j.JobPriorityID > 0 && j.JobPriorityID != existing.JobPriorityID {
		if _, err := uc.jobPriorityRepo.GetByID(ctx, j.JobPriorityID); err != nil {
			return ErrInvalidJobPriority
		}
	}

	// Verificar status si cambió
	if j.JobStatusID > 0 && j.JobStatusID != existing.JobStatusID {
		if _, err := uc.jobStatusRepo.GetByID(ctx, j.JobStatusID); err != nil {
			return ErrInvalidJobStatus
		}
	}

	// Verificar workflow si cambió
	if j.WorkflowID > 0 && j.WorkflowID != existing.WorkflowID {
		if _, err := uc.workflowRepo.GetByID(ctx, j.WorkflowID); err != nil {
			return ErrInvalidWorkflow
		}
	}

	// Verificar usuario si cambió
	if j.UserID != nil && *j.UserID > 0 {
		if existing.UserID == nil || *j.UserID != *existing.UserID {
			if _, err := uc.userRepo.GetByID(ctx, *j.UserID); err != nil {
				return ErrInvalidUser
			}
		}
	}

	// Lógica de tech_status -> job_status automática (fiel al original)
	if j.TechnicianJobStatusID != nil && *j.TechnicianJobStatusID > 0 {
		techChanged := existing.TechnicianJobStatusID == nil || *j.TechnicianJobStatusID != *existing.TechnicianJobStatusID
		if techChanged {
			// Verificar que el tech status existe
			if _, err := uc.technicianJobStatusRepo.GetByID(ctx, *j.TechnicianJobStatusID); err != nil {
				return ErrInvalidTechnicianJobStatus
			}

			// Si el tech status tiene un job_status_id vinculado, actualizar automáticamente
			linkedJobStatusID, err := uc.technicianJobStatusRepo.GetLinkedJobStatusID(ctx, *j.TechnicianJobStatusID)
			if err == nil && linkedJobStatusID != nil {
				j.JobStatusID = *linkedJobStatusID
			}
		}
	}

	if err := uc.repo.Update(ctx, j); err != nil {
		slog.ErrorContext(ctx, "Failed to update job",
			slog.Int64("id", j.ID),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Job updated successfully",
		slog.Int64("id", j.ID))

	return nil
}
