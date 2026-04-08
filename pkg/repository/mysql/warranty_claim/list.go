package warranty_claim

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWC.WarrantyClaim, int, error) {
	where := []string{"wc.deleted_at IS NULL"}
	args := []interface{}{}
	needsJoin := false

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(wc.internal_claim_number LIKE ? OR wc.claim_number LIKE ? OR wc.notes LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if jobID, ok := filters["job_id"].(int64); ok && jobID > 0 {
		where = append(where, "wc.job_id = ?")
		args = append(args, jobID)
	}

	if typeID, ok := filters["warranty_claim_type_id"].(int64); ok && typeID > 0 {
		where = append(where, "wc.warranty_claim_type_id = ?")
		args = append(args, typeID)
	}

	if statusID, ok := filters["warranty_claim_status_id"].(int64); ok && statusID > 0 {
		where = append(where, "wc.warranty_claim_status_id = ?")
		args = append(args, statusID)
	}

	if weekNumber, ok := filters["week_number"].(string); ok && weekNumber != "" {
		needsJoin = true
		where = append(where, "j.week_number = ?")
		args = append(args, weekNumber)
	}

	whereClause := strings.Join(where, " AND ")

	joinClause := ""
	if needsJoin {
		joinClause = "INNER JOIN jobs j ON wc.job_id = j.id"
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranty_claims wc %s WHERE %s", joinClause, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty claims",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sort
	orderBy := "wc.id DESC"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		direction := "ASC"
		if dir, ok := filters["direction"].(string); ok && (dir == "desc" || dir == "DESC") {
			direction = "DESC"
		}
		switch sort {
		case "internal_claim_number":
			orderBy = fmt.Sprintf("wc.internal_claim_number %s", direction)
		case "claim_number":
			orderBy = fmt.Sprintf("wc.claim_number %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("wc.created_at %s", direction)
		case "week_number":
			if !needsJoin {
				joinClause = "INNER JOIN jobs j ON wc.job_id = j.id"
			}
			orderBy = fmt.Sprintf("j.week_number %s", direction)
		default:
			orderBy = fmt.Sprintf("wc.id %s", direction)
		}
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT wc.id, wc.internal_claim_number, wc.warranty_claim_type_id, wc.warranty_claim_status_id, wc.job_id,
			wc.invoice_number, wc.work_done, wc.warranty_part, wc.manufacturer, wc.model_number,
			wc.part_number, wc.replacement_part_number, wc.part_distributor, wc.part_invoice_number,
			wc.old_part_serial_number, wc.new_part_serial_number, wc.esa_number, wc.serial,
			wc.claim_number, wc.approved, wc.parts_credit_received, wc.labor_payment_received,
			wc.notes, wc.created_at, wc.updated_at, wc.deleted_at
		FROM warranty_claims wc
		%s
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, joinClause, whereClause, orderBy)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty claims",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var claims []*domainWC.WarrantyClaim
	for rows.Next() {
		wc := &domainWC.WarrantyClaim{}
		var invoiceNumber, warrantyPart, manufacturer, modelNumber sql.NullString
		var partNumber, replacementPartNumber, partDistributor, partInvoiceNumber sql.NullString
		var oldPartSerialNumber, newPartSerialNumber, esaNumber, serial sql.NullString
		var claimNumber, notes sql.NullString

		if err := rows.Scan(
			&wc.ID,
			&wc.InternalClaimNumber,
			&wc.WarrantyClaimTypeID,
			&wc.WarrantyClaimStatusID,
			&wc.JobID,
			&invoiceNumber, &wc.WorkDone, &warrantyPart, &manufacturer, &modelNumber,
			&partNumber, &replacementPartNumber, &partDistributor, &partInvoiceNumber,
			&oldPartSerialNumber, &newPartSerialNumber, &esaNumber, &serial,
			&claimNumber, &wc.Approved, &wc.PartsCreditReceived, &wc.LaborPaymentReceived,
			&notes, &wc.CreatedAt, &wc.UpdatedAt, &wc.DeletedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty claim row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		wc.InvoiceNumber = fromNullString(invoiceNumber)
		wc.WarrantyPart = fromNullString(warrantyPart)
		wc.Manufacturer = fromNullString(manufacturer)
		wc.ModelNumber = fromNullString(modelNumber)
		wc.PartNumber = fromNullString(partNumber)
		wc.ReplacementPartNumber = fromNullString(replacementPartNumber)
		wc.PartDistributor = fromNullString(partDistributor)
		wc.PartInvoiceNumber = fromNullString(partInvoiceNumber)
		wc.OldPartSerialNumber = fromNullString(oldPartSerialNumber)
		wc.NewPartSerialNumber = fromNullString(newPartSerialNumber)
		wc.EsaNumber = fromNullString(esaNumber)
		wc.Serial = fromNullString(serial)
		wc.ClaimNumber = fromNullString(claimNumber)
		wc.Notes = fromNullString(notes)

		claims = append(claims, wc)
	}

	return claims, total, nil
}
