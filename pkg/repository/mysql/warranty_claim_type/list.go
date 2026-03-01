package warranty_claim_type

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/warranty_claim_type"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*warranty_claim_type.WarrantyClaimType, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(label LIKE ? OR label_plural LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranty_claim_types WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty claim types",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, label_plural, created_at, updated_at
		FROM warranty_claim_types
		WHERE %s
		ORDER BY label ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty claim types",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var types []*warranty_claim_type.WarrantyClaimType
	for rows.Next() {
		wct := &warranty_claim_type.WarrantyClaimType{}

		if err := rows.Scan(
			&wct.ID,
			&wct.Label,
			&wct.LabelPlural,
			&wct.CreatedAt,
			&wct.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty claim type row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		types = append(types, wct)
	}

	return types, total, nil
}
