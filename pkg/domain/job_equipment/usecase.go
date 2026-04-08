package job_equipment

// JobChecker verifica la existencia de un job
type JobChecker interface {
	JobExists(id int64) (bool, error)
}

// UseCase orquesta las operaciones de negocio para equipos de trabajo
type UseCase struct {
	repo       Repository
	jobChecker JobChecker
}

// NewUseCase crea una nueva instancia de UseCase
func NewUseCase(repo Repository, jobChecker JobChecker) *UseCase {
	return &UseCase{
		repo:       repo,
		jobChecker: jobChecker,
	}
}
