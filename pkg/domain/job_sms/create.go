package job_sms

import "context"

func (s *service) Create(ctx context.Context, item *JobSMS) error {
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
