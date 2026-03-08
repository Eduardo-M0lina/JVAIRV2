package alert

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/alert"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*alert.Alert, error) {
	query := `
		SELECT id, user_id, alert_type, entity_id, entity_type, message_level, message, is_read, created_at, updated_at
		FROM alerts
		WHERE id = ?
	`

	a := &alert.Alert{}
	var userID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&userID,
		&a.AlertType,
		&a.EntityID,
		&a.EntityType,
		&a.MessageLevel,
		&a.Message,
		&a.IsRead,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, alert.ErrNotFound
		}
		return nil, err
	}

	if userID.Valid {
		a.UserID = &userID.Int64
	}

	return a, nil
}
