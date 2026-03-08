package job_resident

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_resident"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_resident.JobResident, error) {
	query := `
		SELECT id, job_id, name, mobile_phone, home_phone, email, created_at, updated_at, deleted_at
		FROM job_residents
		WHERE id = ?
	`

	resident := &job_resident.JobResident{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resident.ID,
		&resident.JobID,
		&resident.Name,
		&resident.MobilePhone,
		&resident.HomePhone,
		&resident.Email,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, job_resident.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return resident, nil
}
