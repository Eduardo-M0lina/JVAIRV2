package job_email

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_email"
)

func (r *repository) Create(ctx context.Context, item *job_email.JobEmail) error {
	recipients, err := encodeRecipients(item.Recipients)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO job_emails (job_id, recipients, type, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, item.JobID, recipients, item.Type)
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
