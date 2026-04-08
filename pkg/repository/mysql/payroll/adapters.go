package payroll

import (
	"context"
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/payroll"
)

// UserExistsChecker implementa la interfaz para verificar existencia de usuarios
type userExistsChecker struct {
	db *sql.DB
}

// NewUserExistsChecker crea un nuevo checker de existencia de usuarios
func NewUserExistsChecker(db *sql.DB) payroll.UserExistsChecker {
	return &userExistsChecker{db: db}
}

func (c *userExistsChecker) UserExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND deleted_at IS NULL)"
	err := c.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	return exists, err
}
