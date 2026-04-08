package job_email

import "context"

func (s *service) Create(ctx context.Context, item *JobEmail) error {
	if err := item.ValidateCreate(); err != nil {
		return err
	}

	exists, err := s.jobChecker.JobExists(ctx, item.JobID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobNotFound
	}

	return s.repo.Create(ctx, item)
}
