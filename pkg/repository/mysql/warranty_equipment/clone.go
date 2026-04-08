package warranty_equipment

import (
	"context"
	"log/slog"
)

func (r *Repository) CloneFromJobEquipment(ctx context.Context, warrantyID int64, jobID int64) error {
	query := `
		INSERT INTO warranty_equipment (warranty_id, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at)
		SELECT ?, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			NOW(), NOW()
		FROM job_equipment
		WHERE job_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, warrantyID, jobID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to clone job equipment to warranty",
			slog.String("error", err.Error()),
			slog.Int64("warrantyId", warrantyID),
			slog.Int64("jobId", jobID))
		return err
	}

	return nil
}
