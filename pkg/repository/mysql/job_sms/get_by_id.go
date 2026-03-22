package job_sms

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_sms"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_sms.JobSMS, error) {
	query := `SELECT id, job_id, recipients, type, message, created_at, updated_at FROM job_sms WHERE id = ?`

	item := &job_sms.JobSMS{}
	var recipientsRaw []byte
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.JobID,
		&recipientsRaw,
		&item.Type,
		&item.Message,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, job_sms.ErrNotFound
		}
		return nil, err
	}

	item.Recipients = decodeRecipients(recipientsRaw)
	return item, nil
}
