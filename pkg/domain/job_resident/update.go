package job_resident

import "context"

func (s *service) Update(ctx context.Context, resident *JobResident) error {
	if err := resident.ValidateUpdate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, resident.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.IsDeleted() {
		return ErrNotFound
	}

	return s.repo.Update(ctx, resident)
}
