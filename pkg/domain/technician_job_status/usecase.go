package technician_job_status

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_status"
)

type Service interface {
	Create(ctx context.Context, status *TechnicianJobStatus) error
	GetByID(ctx context.Context, id int64) (*TechnicianJobStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*TechnicianJobStatus, int, error)
	Update(ctx context.Context, status *TechnicianJobStatus) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo          Repository
	jobStatusRepo job_status.Repository
}

func NewUseCase(repo Repository, jobStatusRepo job_status.Repository) *UseCase {
	return &UseCase{
		repo:          repo,
		jobStatusRepo: jobStatusRepo,
	}
}
