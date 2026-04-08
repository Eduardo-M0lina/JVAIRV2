package warranty_equipment

import "context"

// Service define la interfaz del servicio de equipos de garantía
type Service interface {
	Create(ctx context.Context, equipment *WarrantyEquipment) error
	GetByID(ctx context.Context, id int64) (*WarrantyEquipment, error)
	ListByWarrantyID(ctx context.Context, warrantyID int64) ([]*WarrantyEquipment, error)
	Update(ctx context.Context, equipment *WarrantyEquipment) error
	Delete(ctx context.Context, id int64) error
	CloneFromJobEquipment(ctx context.Context, warrantyID int64, jobID int64) error
}

// WarrantyChecker verifica existencia de warranties
type WarrantyChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// UseCase implementa la lógica de negocio de equipos de garantía
type UseCase struct {
	repo          Repository
	warrantyCheck WarrantyChecker
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository, warrantyCheck WarrantyChecker) *UseCase {
	return &UseCase{
		repo:          repo,
		warrantyCheck: warrantyCheck,
	}
}
