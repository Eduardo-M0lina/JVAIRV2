package job_sms

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_sms"
)

type jobExistsChecker struct {
	db *sql.DB
}

func NewJobExistsChecker(db *sql.DB) job_sms.JobExistsChecker {
	return &jobExistsChecker{db: db}
}

func (c *jobExistsChecker) JobExists(ctx context.Context, jobID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ? AND deleted_at IS NULL)"
	err := c.db.QueryRowContext(ctx, query, jobID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
