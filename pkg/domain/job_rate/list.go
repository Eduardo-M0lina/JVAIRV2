package job_rate

import "context"

func (s *service) List(ctx context.Context, jobID int64, limit, offset int) ([]*JobRate, int64, error) {
	return s.repo.List(ctx, jobID, limit, offset)
}
