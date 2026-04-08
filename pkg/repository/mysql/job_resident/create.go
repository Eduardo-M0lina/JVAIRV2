package job_resident

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_resident"
)

func (r *repository) Create(ctx context.Context, resident *job_resident.JobResident) error {
	query := `
		INSERT INTO job_residents (job_id, name, mobile_phone, home_phone, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		resident.JobID,
		resident.Name,
		resident.MobilePhone,
		resident.HomePhone,
		resident.Email,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	resident.ID = id
	return nil
}
