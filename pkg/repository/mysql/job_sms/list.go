package job_sms

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_sms"
)

func (r *repository) List(ctx context.Context, jobID int64, limit, offset int) ([]*job_sms.JobSMS, int64, error) {
	query := `
		SELECT id, job_id, recipients, type, message, created_at, updated_at
		FROM job_sms
		WHERE job_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]*job_sms.JobSMS, 0)
	for rows.Next() {
		item := &job_sms.JobSMS{}
		var recipientsRaw []byte
		if err := rows.Scan(
			&item.ID,
			&item.JobID,
			&recipientsRaw,
			&item.Type,
			&item.Message,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		item.Recipients = decodeRecipients(recipientsRaw)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(*) FROM job_sms WHERE job_id = ?`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, jobID).Scan(&total); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
