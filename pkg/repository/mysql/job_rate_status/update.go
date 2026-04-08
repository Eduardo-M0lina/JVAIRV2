package job_rate_status

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_rate_status"
)

func (r *repository) Update(ctx context.Context, status *job_rate_status.JobRateStatus) error {
	query := `
		UPDATE job_rate_statuses
		SET label = ?, class = ?, ` + "`order`" + ` = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, status.Label, status.Class, status.Order, status.ID)
	return err
}
