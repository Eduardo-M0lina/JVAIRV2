package job_activity_log

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobActivityLog, error) {
	return s.repo.GetByID(ctx, id)
}
