package job_sms

import "context"

func (s *service) GetByID(ctx context.Context, id int64) (*JobSMS, error) {
	return s.repo.GetByID(ctx, id)
}
