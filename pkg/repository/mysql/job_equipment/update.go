package job_equipment

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

func (r *Repository) Update(ctx context.Context, e *job_equipment.JobEquipment) error {
	query := `
		UPDATE job_equipment
		SET job_id = ?, type = ?, area = ?,
		    outdoor_brand = ?, outdoor_model = ?, outdoor_serial = ?, outdoor_installed = ?,
		    furnace_brand = ?, furnace_model = ?, furnace_serial = ?, furnace_installed = ?,
		    evaporator_brand = ?, evaporator_model = ?, evaporator_serial = ?, evaporator_installed = ?,
		    air_handler_brand = ?, air_handler_model = ?, air_handler_serial = ?, air_handler_installed = ?,
		    updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		e.JobID, e.Type, e.Area,
		e.OutdoorBrand, e.OutdoorModel, e.OutdoorSerial, e.OutdoorInstalled,
		e.FurnaceBrand, e.FurnaceModel, e.FurnaceSerial, e.FurnaceInstalled,
		e.EvaporatorBrand, e.EvaporatorModel, e.EvaporatorSerial, e.EvaporatorInstalled,
		e.AirHandlerBrand, e.AirHandlerModel, e.AirHandlerSerial, e.AirHandlerInstalled,
		e.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute update job_equipment query",
			slog.String("error", err.Error()))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get rows affected",
			slog.String("error", err.Error()))
		return err
	}

	if rowsAffected == 0 {
		slog.WarnContext(ctx, "No job equipment updated",
			slog.Int64("equipment_id", e.ID))
	}

	return nil
}
