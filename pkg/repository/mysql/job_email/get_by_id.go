package job_email

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_email"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_email.JobEmail, error) {
	query := `SELECT id, job_id, recipients, type, created_at, updated_at FROM job_emails WHERE id = ?`

	item := &job_email.JobEmail{}
	var recipientsRaw []byte
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.JobID,
		&recipientsRaw,
		&item.Type,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, job_email.ErrNotFound
		}
		return nil, err
	}

	item.Recipients = decodeRecipients(recipientsRaw)
	return item, nil
}
