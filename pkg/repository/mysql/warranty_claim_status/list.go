package warranty_claim_status

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*warranty_claim_status.WarrantyClaimStatus, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "label LIKE ?")
		args = append(args, "%"+search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranty_claim_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty claim statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, class, `+"`order`"+`, is_active, created_at, updated_at
		FROM warranty_claim_statuses
		WHERE %s
		ORDER BY `+"`order`"+` ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty claim statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*warranty_claim_status.WarrantyClaimStatus
	for rows.Next() {
		wcs := &warranty_claim_status.WarrantyClaimStatus{}
		var class sql.NullString

		if err := rows.Scan(
			&wcs.ID,
			&wcs.Label,
			&class,
			&wcs.Order,
			&wcs.IsActive,
			&wcs.CreatedAt,
			&wcs.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty claim status row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if class.Valid {
			wcs.Class = &class.String
		}

		statuses = append(statuses, wcs)
	}

	return statuses, total, nil
}
