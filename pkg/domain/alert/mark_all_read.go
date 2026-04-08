package alert

import "context"

func (s *service) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	return s.repo.MarkAllRead(ctx, userID)
}
