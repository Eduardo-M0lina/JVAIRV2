package job_rate_status

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_rate_status"
)

func (r *repository) Create(ctx context.Context, status *job_rate_status.JobRateStatus) error {
	query := `
		INSERT INTO job_rate_statuses (label, class, ` + "`order`" + `, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, status.Label, status.Class, status.Order)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	status.ID = id
	return nil
}
