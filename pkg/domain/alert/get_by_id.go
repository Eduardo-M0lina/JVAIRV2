package alert

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*Alert, error) {
	return s.repo.GetByID(ctx, id)
}
