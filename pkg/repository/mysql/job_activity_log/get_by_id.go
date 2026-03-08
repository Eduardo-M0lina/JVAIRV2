package job_activity_log

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_activity_log"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_activity_log.JobActivityLog, error) {
	query := `
		SELECT id, job_id, type, log, user_id, created_at, updated_at
		FROM job_activity_logs
		WHERE id = ?
	`

	log := &job_activity_log.JobActivityLog{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.JobID,
		&log.Type,
		&log.Log,
		&log.UserID,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, job_activity_log.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return log, nil
}
