package job_rate_status

import "context"

func (s *service) Update(ctx context.Context, status *JobRateStatus) error {
	if err := status.ValidateUpdate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, status.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Update(ctx, status)
}
