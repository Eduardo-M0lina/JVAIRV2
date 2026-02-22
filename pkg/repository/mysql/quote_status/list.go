package quote_status

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/quote_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*quote_status.QuoteStatus, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "label LIKE ?")
		args = append(args, "%"+search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quote_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count quote statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, class, `+"`order`"+`, created_at, updated_at
		FROM quote_statuses
		WHERE %s
		ORDER BY `+"`order`"+` ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list quote statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*quote_status.QuoteStatus
	for rows.Next() {
		qs := &quote_status.QuoteStatus{}
		var class sql.NullString

		if err := rows.Scan(
			&qs.ID,
			&qs.Label,
			&class,
			&qs.Order,
			&qs.CreatedAt,
			&qs.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan quote status row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if class.Valid {
			qs.Class = &class.String
		}

		statuses = append(statuses, qs)
	}

	return statuses, total, nil
}
