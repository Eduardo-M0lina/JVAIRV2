package job_visit

import (
	"context"
	"database/sql"
	"log/slog"
)

// JobExistsChecker verifica la existencia de un job en la base de datos
type JobExistsChecker struct {
	db *sql.DB
}

// NewJobExistsChecker crea una nueva instancia del checker
func NewJobExistsChecker(db *sql.DB) *JobExistsChecker {
	return &JobExistsChecker{db: db}
}

// JobExists verifica si un job existe
func (c *JobExistsChecker) JobExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ?)`
	err := c.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check job existence",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return false, err
	}
	return exists, nil
}

// UserExistsChecker verifica la existencia de un usuario en la base de datos
type UserExistsChecker struct {
	db *sql.DB
}

// NewUserExistsChecker crea una nueva instancia del checker
func NewUserExistsChecker(db *sql.DB) *UserExistsChecker {
	return &UserExistsChecker{db: db}
}

// UserExists verifica si un usuario existe
func (c *UserExistsChecker) UserExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err := c.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check user existence",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return false, err
	}
	return exists, nil
}
