package job_resident

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_resident"
)

func (r *repository) List(ctx context.Context, jobID int64, limit, offset int) ([]*job_resident.JobResident, int64, error) {
	query := `
		SELECT id, job_id, name, mobile_phone, home_phone, email, created_at, updated_at, deleted_at
		FROM job_residents
		WHERE job_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var residents []*job_resident.JobResident
	for rows.Next() {
		resident := &job_resident.JobResident{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, err
		}
		residents = append(residents, resident)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM job_residents WHERE job_id = ? AND deleted_at IS NULL"
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, jobID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return residents, total, nil
}
