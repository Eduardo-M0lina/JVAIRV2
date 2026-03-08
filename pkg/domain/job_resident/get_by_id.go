package job_resident

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobResident, error) {
	return s.repo.GetByID(ctx, id)
}
