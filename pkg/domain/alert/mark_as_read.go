package alert

import "context"

func (s *service) MarkAsRead(ctx context.Context, id int64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.MarkAsRead(ctx, id)
}
