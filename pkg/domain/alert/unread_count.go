package alert

import "context"

func (s *service) UnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}
