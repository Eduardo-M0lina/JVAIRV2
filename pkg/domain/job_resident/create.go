package job_resident

import "context"

func (s *service) Create(ctx context.Context, resident *JobResident) error {
	if err := resident.ValidateCreate(); err != nil {
		return err
	}

	exists, err := s.jobChecker.JobExists(ctx, resident.JobID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobNotFound
	}

	return s.repo.Create(ctx, resident)
}
