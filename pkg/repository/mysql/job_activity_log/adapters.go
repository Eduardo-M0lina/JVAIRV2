package job_activity_log

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_activity_log"
)

type jobExistsChecker struct {
	db *sql.DB
}

func NewJobExistsChecker(db *sql.DB) job_activity_log.JobExistsChecker {
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

func NewUserExistsChecker(db *sql.DB) job_activity_log.UserExistsChecker {
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
