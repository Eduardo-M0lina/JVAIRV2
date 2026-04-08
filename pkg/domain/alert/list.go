package alert

import "context"

func (s *service) List(ctx context.Context, filters ListFilters, limit, offset int) ([]*Alert, int64, error) {
	return s.repo.List(ctx, filters, limit, offset)
}
