package job_activity_log

import "context"

func (s *service) List(ctx context.Context, jobID int64, limit, offset int) ([]*JobActivityLog, int64, error) {
	return s.repo.List(ctx, jobID, limit, offset)
}
