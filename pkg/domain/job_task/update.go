package job_task

import "context"

func (s *service) Update(ctx context.Context, task *JobTask) error {
	if err := task.ValidateUpdate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.IsDeleted() {
		return ErrNotFound
	}

	exists, err := s.userChecker.UserExists(ctx, task.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	exists, err = s.taskStatusChecker.TaskStatusExists(ctx, task.TaskStatusID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTaskStatusNotFound
	}

	return s.repo.Update(ctx, task)
}
