package job_task

import "context"

func (s *service) Delete(ctx context.Context, id int64) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil || task.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
