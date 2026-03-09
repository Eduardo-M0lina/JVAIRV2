package job_rate_status

import "context"

func (s *service) List(ctx context.Context) ([]*JobRateStatus, error) {
	return s.repo.List(ctx)
}
