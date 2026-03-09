package job_rate_status

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_rate_status"
)

func (r *repository) List(ctx context.Context) ([]*job_rate_status.JobRateStatus, error) {
	query := `
		SELECT id, label, class, ` + "`order`" + `, created_at, updated_at, deleted_at
		FROM job_rate_statuses
		WHERE deleted_at IS NULL
		ORDER BY ` + "`order`" + ` ASC, label ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*job_rate_status.JobRateStatus
	for rows.Next() {
		status := &job_rate_status.JobRateStatus{}
		err := rows.Scan(
			&status.ID,
			&status.Label,
			&status.Class,
			&status.Order,
			&status.CreatedAt,
			&status.UpdatedAt,
			&status.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}
