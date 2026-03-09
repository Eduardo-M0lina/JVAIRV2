package job_rate

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobRate, error) {
	return s.repo.GetByID(ctx, id)
}
