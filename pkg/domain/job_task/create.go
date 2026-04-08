package job_task

import "context"

func (s *service) Create(ctx context.Context, task *JobTask) error {
	if err := task.ValidateCreate(); err != nil {
		return err
	}

	exists, err := s.jobChecker.JobExists(ctx, task.JobID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobNotFound
	}

	exists, err = s.userChecker.UserExists(ctx, task.UserID)
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

	return s.repo.Create(ctx, task)
}
