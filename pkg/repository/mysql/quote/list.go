package quote

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainQuote.Quote, int64, error) {
	where := []string{"q.deleted_at IS NULL"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "q.quote_number LIKE ?")
		args = append(args, "%"+search+"%")
	}

	if jobID, ok := filters["job_id"].(int64); ok {
		where = append(where, "q.job_id = ?")
		args = append(args, jobID)
	}

	if quoteStatusID, ok := filters["quote_status_id"].(int64); ok {
		where = append(where, "q.quote_status_id = ?")
		args = append(args, quoteStatusID)
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quotes q WHERE %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count quotes",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sorting
	sortColumn := "q.created_at"
	sortDirection := "DESC"
	if sort, ok := filters["sort"].(string); ok {
		switch sort {
		case "quote_number":
			sortColumn = "q.quote_number"
		case "amount":
			sortColumn = "q.amount"
		case "created_at":
			sortColumn = "q.created_at"
		}
	}
	if direction, ok := filters["direction"].(string); ok {
		switch strings.ToUpper(direction) {
		case "ASC":
			sortDirection = "ASC"
		case "DESC":
			sortDirection = "DESC"
		}
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT q.id, q.job_id, q.quote_number, q.quote_status_id, q.amount,
		       q.description, q.notes, q.created_at, q.updated_at
		FROM quotes q
		WHERE %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, sortColumn, sortDirection)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list quotes",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var quotes []*domainQuote.Quote
	for rows.Next() {
		q := &domainQuote.Quote{}
		var description, notes sql.NullString

		if err := rows.Scan(
			&q.ID,
			&q.JobID,
			&q.QuoteNumber,
			&q.QuoteStatusID,
			&q.Amount,
			&description,
			&notes,
			&q.CreatedAt,
			&q.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan quote row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if description.Valid {
			q.Description = &description.String
		}
		if notes.Valid {
			q.Notes = &notes.String
		}

		quotes = append(quotes, q)
	}

	return quotes, total, nil
}
