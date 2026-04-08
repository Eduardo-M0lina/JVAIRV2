package job_sms

import "context"

func (s *service) List(ctx context.Context, jobID int64, limit, offset int) ([]*JobSMS, int64, error) {
	if limit < 1 {
		limit = 15
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, jobID, limit, offset)
}
