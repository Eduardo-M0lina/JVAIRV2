package alert

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/alert"
)

func (r *repository) Create(ctx context.Context, a *alert.Alert) error {
	query := `
		INSERT INTO alerts (user_id, alert_type, entity_id, entity_type, message_level, message, is_read, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	var userID sql.NullInt64
	if a.UserID != nil {
		userID = sql.NullInt64{Int64: *a.UserID, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query, userID, a.AlertType, a.EntityID, a.EntityType, a.MessageLevel, a.Message, a.IsRead)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = id
	return nil
}
