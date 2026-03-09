package job_rate

import "context"

func (s *service) Create(ctx context.Context, rate *JobRate) error {
	if err := rate.ValidateCreate(); err != nil {
		return err
	}

	exists, err := s.jobChecker.JobExists(ctx, rate.JobID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobNotFound
	}

	exists, err = s.userChecker.UserExists(ctx, rate.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	exists, err = s.statusChecker.JobRateStatusExists(ctx, rate.JobRateStatusID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJobRateStatusNotFound
	}

	rate.Payment = CalculatePayment(
		rate.SalePrice,
		rate.RatePercent,
		rate.RateFlat,
		rate.TechParts,
		rate.CompanyParts,
		rate.Deduction,
	)

	return s.repo.Create(ctx, rate)
}
