package job_rate

import "context"

func (s *service) Update(ctx context.Context, rate *JobRate) error {
	if err := rate.ValidateUpdate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, rate.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.IsDeleted() {
		return ErrNotFound
	}

	exists, err := s.userChecker.UserExists(ctx, rate.UserID)
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

	return s.repo.Update(ctx, rate)
}
