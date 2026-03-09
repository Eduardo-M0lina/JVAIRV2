package job_rate

import "context"

func (s *service) Delete(ctx context.Context, id int64) error {
	rate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rate == nil || rate.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
