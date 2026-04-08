package job_rate_status

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobRateStatus, error) {
	return s.repo.GetByID(ctx, id)
}
