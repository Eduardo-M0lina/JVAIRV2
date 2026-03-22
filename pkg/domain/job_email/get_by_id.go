package job_email

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobEmail, error) {
	return s.repo.GetByID(ctx, id)
}
