package job_task

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobTask, error) {
	return s.repo.GetByID(ctx, id)
}
