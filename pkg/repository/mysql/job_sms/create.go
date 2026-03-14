package job_sms

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_sms"
)

func (r *repository) Create(ctx context.Context, item *job_sms.JobSMS) error {
	recipients, err := encodeRecipients(item.Recipients)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO job_sms (job_id, recipients, type, message, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, item.JobID, recipients, item.Type, item.Message)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = id
	return nil
}
