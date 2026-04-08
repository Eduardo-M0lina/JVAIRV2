package workflow

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrWorkflowNotFound se devuelve cuando no se encuentra un workflow
	ErrWorkflowNotFound = errors.New("workflow no encontrado")
	// ErrWorkflowNameRequired se devuelve cuando el nombre del workflow está vacío
	ErrWorkflowNameRequired = errors.New("el nombre del workflow es requerido")
)

// UseCase maneja la lógica de negocio de workflows
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso de workflows
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// List obtiene una lista paginada de workflows
func (uc *UseCase) List(ctx context.Context, filters Filters, page, pageSize int) ([]Workflow, int64, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}

// GetByID obtiene un workflow por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Workflow, error) {
	workflow, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Obtener los statuses asociados al workflow
	statuses, err := uc.repo.GetWorkflowStatuses(ctx, id)
	if err != nil {
		return nil, err
	}

	workflow.Statuses = statuses
	return workflow, nil
}

// Create crea un nuevo workflow
func (uc *UseCase) Create(ctx context.Context, workflow *Workflow, statusIDs []int64) error {
	// Validar que el nombre no esté vacío
	if workflow.Name == "" {
		return ErrWorkflowNameRequired
	}

	// Crear el workflow
	if err := uc.repo.Create(ctx, workflow); err != nil {
		return err
	}

	// Si se proporcionaron statuses, asociarlos al workflow
	if len(statusIDs) > 0 {
		statuses := make([]WorkflowStatus, len(statusIDs))
		for i, statusID := range statusIDs {
			statuses[i] = WorkflowStatus{
				JobStatusID: statusID,
				WorkflowID:  workflow.ID,
				Order:       i,
			}
		}

		if err := uc.repo.SetWorkflowStatuses(ctx, workflow.ID, statuses); err != nil {
			return err
		}
	}

	return nil
}

// Update actualiza un workflow existente
func (uc *UseCase) Update(ctx context.Context, workflow *Workflow, statusIDs []int64) error {
	// Validar que el nombre no esté vacío
	if workflow.Name == "" {
		return ErrWorkflowNameRequired
	}

	// Verificar que el workflow existe
	existing, err := uc.repo.GetByID(ctx, workflow.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrWorkflowNotFound
	}

	// Actualizar el workflow
	if err := uc.repo.Update(ctx, workflow); err != nil {
		return err
	}

	// Actualizar los statuses asociados
	if statusIDs != nil {
		statuses := make([]WorkflowStatus, len(statusIDs))
		for i, statusID := range statusIDs {
			statuses[i] = WorkflowStatus{
				JobStatusID: statusID,
				WorkflowID:  workflow.ID,
				Order:       i,
			}
		}

		if err := uc.repo.SetWorkflowStatuses(ctx, workflow.ID, statuses); err != nil {
			return err
		}
	}

	return nil
}

// Delete elimina un workflow
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Verificar que el workflow existe
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrWorkflowNotFound
	}

	return uc.repo.Delete(ctx, id)
}

// Duplicate duplica un workflow existente
func (uc *UseCase) Duplicate(ctx context.Context, id int64) (*Workflow, error) {
	// Verificar que el workflow existe
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrWorkflowNotFound
	}

	// Obtener los statuses del workflow original
	statuses, err := uc.repo.GetWorkflowStatuses(ctx, id)
	if err != nil {
		return nil, err
	}

	// Duplicar el workflow
	duplicated, err := uc.repo.Duplicate(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actualizar el nombre del workflow duplicado
	duplicated.Name = fmt.Sprintf("Copy of %s (%d)", existing.Name, duplicated.ID)
	if err := uc.repo.Update(ctx, duplicated); err != nil {
		return nil, err
	}

	// Copiar los statuses al workflow duplicado
	if len(statuses) > 0 {
		newStatuses := make([]WorkflowStatus, len(statuses))
		for i, status := range statuses {
			newStatuses[i] = WorkflowStatus{
				JobStatusID: status.JobStatusID,
				WorkflowID:  duplicated.ID,
				Order:       status.Order,
			}
		}

		if err := uc.repo.SetWorkflowStatuses(ctx, duplicated.ID, newStatuses); err != nil {
			return nil, err
		}
	}

	// Obtener el workflow completo con sus statuses
	return uc.GetByID(ctx, duplicated.ID)
}
