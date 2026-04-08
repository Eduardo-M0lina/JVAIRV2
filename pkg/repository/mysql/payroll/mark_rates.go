package payroll

import (
	"context"
	"fmt"
	"strings"
)

func (r *repository) MarkRatesAsPaid(ctx context.Context, rateIDs []int64) error {
	if len(rateIDs) == 0 {
		return nil
	}

	// Obtener el ID del status "Paid"
	statusID, err := r.GetStatusIDByLabel(ctx, "Paid")
	if err != nil {
		return fmt.Errorf("error getting Paid status ID: %w", err)
	}

	// Construir placeholders para IN clause
	placeholders := make([]string, len(rateIDs))
	args := make([]interface{}, len(rateIDs)+1)
	args[0] = statusID
	for i, id := range rateIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		UPDATE job_rates
		SET job_rate_status_id = ?, updated_at = NOW()
		WHERE id IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ","))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error marking rates as paid: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rates were updated - verify that rate IDs %v exist and are not deleted", rateIDs)
	}

	return nil
}

func (r *repository) MarkRatesAsHolding(ctx context.Context, rateIDs []int64) error {
	if len(rateIDs) == 0 {
		return nil
	}

	// Obtener el ID del status "Holding"
	statusID, err := r.GetStatusIDByLabel(ctx, "Holding")
	if err != nil {
		return fmt.Errorf("error getting Holding status ID: %w", err)
	}

	// Construir placeholders para IN clause
	placeholders := make([]string, len(rateIDs))
	args := make([]interface{}, len(rateIDs)+1)
	args[0] = statusID
	for i, id := range rateIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		UPDATE job_rates
		SET job_rate_status_id = ?, updated_at = NOW()
		WHERE id IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ","))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error marking rates as holding: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rates were updated - verify that rate IDs %v exist and are not deleted", rateIDs)
	}

	return nil
}

func (r *repository) GetStatusIDByLabel(ctx context.Context, label string) (int64, error) {
	var statusID int64
	query := "SELECT id FROM job_rate_statuses WHERE label = ? AND deleted_at IS NULL LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, label).Scan(&statusID)
	if err != nil {
		return 0, fmt.Errorf("status '%s' not found: %w", label, err)
	}
	return statusID, nil
}
