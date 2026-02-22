package job

import "context"

// Service define la interfaz del servicio de jobs
type Service interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id int64) (*Job, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Job, int, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id int64) error
	Close(ctx context.Context, id int64, jobStatusID int64) error
}

// UseCase implementa la lógica de negocio de jobs
type UseCase struct {
	repo                    Repository
	jobCategoryRepo         JobCategoryChecker
	jobPriorityRepo         JobPriorityChecker
	jobStatusRepo           JobStatusChecker
	workflowRepo            WorkflowChecker
	propertyRepo            PropertyChecker
	userRepo                UserChecker
	technicianJobStatusRepo TechnicianJobStatusChecker
}

// JobCategoryChecker verifica existencia de categorías de trabajo
type JobCategoryChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// JobPriorityChecker verifica existencia de prioridades de trabajo
type JobPriorityChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// JobStatusChecker verifica existencia de estados de trabajo
type JobStatusChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// WorkflowChecker verifica existencia de workflows y obtiene status inicial
type WorkflowChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	GetInitialStatusID(ctx context.Context, workflowID int64) (int64, error)
}

// PropertyChecker verifica existencia de propiedades
type PropertyChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	GetWorkflowID(ctx context.Context, propertyID int64) (int64, error)
}

// UserChecker verifica existencia de usuarios
type UserChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// TechnicianJobStatusChecker verifica existencia de estados de técnico y obtiene job_status_id vinculado
type TechnicianJobStatusChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	GetLinkedJobStatusID(ctx context.Context, id int64) (*int64, error)
}

// NewUseCase crea una nueva instancia del caso de uso de jobs
func NewUseCase(
	repo Repository,
	jobCategoryRepo JobCategoryChecker,
	jobPriorityRepo JobPriorityChecker,
	jobStatusRepo JobStatusChecker,
	workflowRepo WorkflowChecker,
	propertyRepo PropertyChecker,
	userRepo UserChecker,
	technicianJobStatusRepo TechnicianJobStatusChecker,
) *UseCase {
	return &UseCase{
		repo:                    repo,
		jobCategoryRepo:         jobCategoryRepo,
		jobPriorityRepo:         jobPriorityRepo,
		jobStatusRepo:           jobStatusRepo,
		workflowRepo:            workflowRepo,
		propertyRepo:            propertyRepo,
		userRepo:                userRepo,
		technicianJobStatusRepo: technicianJobStatusRepo,
	}
}
