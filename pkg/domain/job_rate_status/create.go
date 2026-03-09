package job_rate_status

import "context"

func (s *service) Create(ctx context.Context, status *JobRateStatus) error {
	if err := status.ValidateCreate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, status)
}
