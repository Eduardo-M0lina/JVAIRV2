package alert

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/alert"
)

type userExistsChecker struct {
	db *sql.DB
}

func NewUserExistsChecker(db *sql.DB) alert.UserExistsChecker {
	return &userExistsChecker{db: db}
}

func (c *userExistsChecker) UserExists(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	var exists int
	err := c.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
