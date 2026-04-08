package job_rate_status

import (
	"context"
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_rate_status"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_rate_status.JobRateStatus, error) {
	query := `
		SELECT id, label, class, ` + "`order`" + `, created_at, updated_at, deleted_at
		FROM job_rate_statuses
		WHERE id = ?
	`

	status := &job_rate_status.JobRateStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&status.ID,
		&status.Label,
		&status.Class,
		&status.Order,
		&status.CreatedAt,
		&status.UpdatedAt,
		&status.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, job_rate_status.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return status, nil
}
