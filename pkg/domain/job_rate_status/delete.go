package job_rate_status

import "context"

func (s *service) Delete(ctx context.Context, id int64) error {
	status, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if status == nil || status.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
