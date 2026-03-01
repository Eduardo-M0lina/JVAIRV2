package job_visit

import "context"

// Service define las operaciones de negocio para visitas de trabajo
type Service interface {
	Create(ctx context.Context, visit *JobVisit) (int64, error)
	GetByID(ctx context.Context, id int64) (*JobVisit, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int, sort, direction string) ([]*JobVisit, int64, error)
	Update(ctx context.Context, visit *JobVisit) error
	Delete(ctx context.Context, id int64) error
}

// UseCase implementa la lógica de negocio para visitas de trabajo
type UseCase struct {
	repo        Repository
	jobChecker  JobChecker
	userChecker UserChecker
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository, jobChecker JobChecker, userChecker UserChecker) *UseCase {
	return &UseCase{
		repo:        repo,
		jobChecker:  jobChecker,
		userChecker: userChecker,
	}
}
