package job_activity_log

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_activity_log"
)

func (r *repository) Create(ctx context.Context, log *job_activity_log.JobActivityLog) error {
	query := `
		INSERT INTO job_activity_logs (job_id, type, log, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, log.JobID, log.Type, log.Log, log.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.ID = id
	return nil
}
