package job_resident

import "context"

func (s *service) Delete(ctx context.Context, id int64) error {
	resident, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if resident == nil || resident.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
