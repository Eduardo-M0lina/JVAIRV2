package password_history

import (
	"context"
	"database/sql"

	domainPasswordHistory "github.com/angumol/jvairv2/pkg/domain/password_history"
)

// Create crea un nuevo registro de historial
func (r *Repository) Create(ctx context.Context, ph *domainPasswordHistory.PasswordHistory) error {
	query := `
		INSERT INTO password_history (user_id, password, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, ph.UserID, ph.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ph.ID = id
	return nil
}

// GetByUserID obtiene el historial de contraseñas de un usuario
func (r *Repository) GetByUserID(ctx context.Context, userID int64, limit int) ([]*domainPasswordHistory.PasswordHistory, error) {
	query := `
		SELECT id, user_id, password, created_at, updated_at
		FROM password_history
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += ` LIMIT ?`
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = r.db.QueryContext(ctx, query, userID, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Log error but don't return it since we're in a defer
			_ = err
		}
	}()

	var histories []*domainPasswordHistory.PasswordHistory
	for rows.Next() {
		var ph domainPasswordHistory.PasswordHistory
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ph.ID,
			&ph.UserID,
			&ph.Password,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			ph.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			ph.UpdatedAt = &updatedAt.Time
		}

		histories = append(histories, &ph)
	}

	return histories, rows.Err()
}

// Delete elimina un registro de historial por ID
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM password_history WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteOldest elimina los registros más antiguos de un usuario, manteniendo solo los N más recientes
func (r *Repository) DeleteOldest(ctx context.Context, userID int64, keep int) error {
	query := `
		DELETE FROM password_history
		WHERE user_id = ?
		AND id NOT IN (
			SELECT id FROM (
				SELECT id
				FROM password_history
				WHERE user_id = ?
				ORDER BY created_at DESC
				LIMIT ?
			) AS subquery
		)
	`
	_, err := r.db.ExecContext(ctx, query, userID, userID, keep)
	return err
}
