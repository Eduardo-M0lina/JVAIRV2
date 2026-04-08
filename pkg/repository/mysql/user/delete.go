package user

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// Delete elimina un usuario (soft delete)
func (r *Repository) Delete(ctx context.Context, id string) error {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	query := `
		UPDATE users
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()

	result, err := r.db.ExecContext(ctx, query, now, idInt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
