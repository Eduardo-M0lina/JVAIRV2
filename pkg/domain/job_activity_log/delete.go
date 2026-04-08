package job_activity_log

import "context"

func (s *service) Delete(ctx context.Context, id int64) error {
	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if log == nil {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
