package job_task

import "context"

func (s *service) List(ctx context.Context, jobID int64, limit, offset int) ([]*JobTask, int64, error) {
	return s.repo.List(ctx, jobID, limit, offset)
}

func (s *service) ListAll(ctx context.Context, limit, offset int) ([]*JobTask, int64, error) {
	return s.repo.ListAll(ctx, limit, offset)
}
