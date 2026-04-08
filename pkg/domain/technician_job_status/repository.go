package technician_job_status

import "context"

type Repository interface {
	Create(ctx context.Context, status *TechnicianJobStatus) error
	GetByID(ctx context.Context, id int64) (*TechnicianJobStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*TechnicianJobStatus, int, error)
	Update(ctx context.Context, status *TechnicianJobStatus) error
	Delete(ctx context.Context, id int64) error
}
