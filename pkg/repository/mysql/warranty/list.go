package warranty

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWarranty.Warranty, int, error) {
	where := []string{"w.deleted_at IS NULL"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(w.warranty_number LIKE ? OR w.agreement_number LIKE ? OR w.notes LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if jobID, ok := filters["job_id"].(int64); ok && jobID > 0 {
		where = append(where, "w.job_id = ?")
		args = append(args, jobID)
	}

	if typeID, ok := filters["warranty_type_id"].(int64); ok && typeID > 0 {
		where = append(where, "w.warranty_type_id = ?")
		args = append(args, typeID)
	}

	if statusID, ok := filters["warranty_status_id"].(int64); ok && statusID > 0 {
		where = append(where, "w.warranty_status_id = ?")
		args = append(args, statusID)
	}

	if weekNumber, ok := filters["week_number"].(string); ok && weekNumber != "" {
		where = append(where, "j.week_number = ?")
		args = append(args, weekNumber)
	}

	whereClause := strings.Join(where, " AND ")

	joinClause := `
		INNER JOIN jobs j ON w.job_id = j.id
		INNER JOIN properties p ON j.property_id = p.id
		INNER JOIN customers c ON p.customer_id = c.id`

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranties w %s WHERE %s", joinClause, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranties",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sort
	orderBy := "w.id DESC"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		direction := "ASC"
		if dir, ok := filters["direction"].(string); ok && (dir == "desc" || dir == "DESC") {
			direction = "DESC"
		}
		switch sort {
		case "warranty_number":
			orderBy = fmt.Sprintf("w.warranty_number %s", direction)
		case "date_submitted":
			orderBy = fmt.Sprintf("w.date_submitted %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("w.created_at %s", direction)
		case "week_number":
			orderBy = fmt.Sprintf("j.week_number %s", direction)
		default:
			orderBy = fmt.Sprintf("w.id %s", direction)
		}
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT w.id, w.warranty_number, w.job_id, w.warranty_type_id, w.warranty_status_id,
			w.date_submitted, w.agreement_number, w.audit_done, w.notes,
			w.created_at, w.updated_at, w.deleted_at,
			j.id, j.week_number, j.completion_date,
			p.id, CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip),
			c.name
		FROM warranties w
		%s
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, joinClause, whereClause, orderBy)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranties",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var warranties []*domainWarranty.Warranty
	for rows.Next() {
		w := &domainWarranty.Warranty{}
		var dateSubmitted sql.NullTime
		var agreementNumber sql.NullString
		var notes sql.NullString
		var jobID int64
		var weekNumber sql.NullInt32
		var completionDate sql.NullTime
		var propertyID int64
		var propertyAddress string
		var customerName string

		if err := rows.Scan(
			&w.ID,
			&w.WarrantyNumber,
			&w.JobID,
			&w.WarrantyTypeID,
			&w.WarrantyStatusID,
			&dateSubmitted,
			&agreementNumber,
			&w.AuditDone,
			&notes,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.DeletedAt,
			&jobID, &weekNumber, &completionDate,
			&propertyID, &propertyAddress,
			&customerName,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if dateSubmitted.Valid {
			w.DateSubmitted = &dateSubmitted.Time
		}
		if agreementNumber.Valid {
			w.AgreementNumber = &agreementNumber.String
		}
		if notes.Valid {
			w.Notes = &notes.String
		}

		w.Job = &domainWarranty.Job{
			ID: jobID,
			Property: domainWarranty.Property{
				ID:      propertyID,
				Address: propertyAddress,
				Customer: domainWarranty.Customer{
					Name: customerName,
				},
			},
		}

		if weekNumber.Valid {
			wn := int(weekNumber.Int32)
			w.Job.WeekNumber = &wn
		}
		if completionDate.Valid {
			w.Job.CompletionDate = &completionDate.Time
		}

		warranties = append(warranties, w)
	}

	return warranties, total, nil
}
