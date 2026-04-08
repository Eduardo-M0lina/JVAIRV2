package job_rate

import (
	"context"
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_rate"
)

type jobExistsChecker struct {
	db *sql.DB
}

func NewJobExistsChecker(db *sql.DB) job_rate.JobExistsChecker {
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

type userExistsChecker struct {
	db *sql.DB
}

func NewUserExistsChecker(db *sql.DB) job_rate.UserExistsChecker {
	return &userExistsChecker{db: db}
}

func (c *userExistsChecker) UserExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err := c.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type jobRateStatusExistsChecker struct {
	db *sql.DB
}

func NewJobRateStatusExistsChecker(db *sql.DB) job_rate.JobRateStatusExistsChecker {
	return &jobRateStatusExistsChecker{db: db}
}

func (c *jobRateStatusExistsChecker) JobRateStatusExists(ctx context.Context, statusID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM job_rate_statuses WHERE id = ? AND deleted_at IS NULL)"
	err := c.db.QueryRowContext(ctx, query, statusID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
