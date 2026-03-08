package alert

import "context"

func (s *service) Create(ctx context.Context, alert *Alert) error {
	if err := alert.ValidateCreate(); err != nil {
		return err
	}

	if alert.UserID != nil {
		exists, err := s.userChecker.UserExists(ctx, *alert.UserID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrUserNotFound
		}
	}

	return s.repo.Create(ctx, alert)
}
