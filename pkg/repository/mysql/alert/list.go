package alert

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/alert"
)

func (r *repository) List(ctx context.Context, filters alert.ListFilters, limit, offset int) ([]*alert.Alert, int64, error) {
	var conditions []string
	var args []interface{}

	if filters.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filters.UserID)
	}

	if filters.IsRead != nil {
		conditions = append(conditions, "is_read = ?")
		args = append(args, *filters.IsRead)
	}

	if filters.AlertType != nil {
		conditions = append(conditions, "alert_type = ?")
		args = append(args, *filters.AlertType)
	}

	if filters.EntityType != nil {
		conditions = append(conditions, "entity_type = ?")
		args = append(args, *filters.EntityType)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, alert_type, entity_id, entity_type, message_level, message, is_read, created_at, updated_at
		FROM alerts
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	queryArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var alerts []*alert.Alert
	for rows.Next() {
		a := &alert.Alert{}
		var userID sql.NullInt64

		err := rows.Scan(
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
			return nil, 0, err
		}

		if userID.Valid {
			a.UserID = &userID.Int64
		}

		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alerts %s", whereClause)
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}
