package job_activity_log

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_activity_log"
)

func (r *repository) List(ctx context.Context, jobID int64, limit, offset int) ([]*job_activity_log.JobActivityLog, int64, error) {
	query := `
		SELECT id, job_id, type, log, user_id, created_at, updated_at
		FROM job_activity_logs
		WHERE job_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var logs []*job_activity_log.JobActivityLog
	for rows.Next() {
		log := &job_activity_log.JobActivityLog{}
		err := rows.Scan(
			&log.ID,
			&log.JobID,
			&log.Type,
			&log.Log,
			&log.UserID,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM job_activity_logs WHERE job_id = ?"
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, jobID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
