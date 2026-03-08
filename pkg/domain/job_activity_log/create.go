package job_activity_log

import "context"

func (s *service) Create(ctx context.Context, log *JobActivityLog) error {
	if err := log.ValidateCreate(); err != nil {
		return err
	}

	exists, err := s.jobChecker.JobExists(ctx, log.JobID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobNotFound
	}

	exists, err = s.userChecker.UserExists(ctx, log.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	return s.repo.Create(ctx, log)
}
