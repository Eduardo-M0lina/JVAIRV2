package job_resident

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_resident"
)

func (r *repository) Update(ctx context.Context, resident *job_resident.JobResident) error {
	query := `
		UPDATE job_residents
		SET name = ?, mobile_phone = ?, home_phone = ?, email = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		resident.Name,
		resident.MobilePhone,
		resident.HomePhone,
		resident.Email,
		resident.ID,
	)

	return err
}
