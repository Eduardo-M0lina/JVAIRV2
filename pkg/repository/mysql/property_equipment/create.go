package property_equipment

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/property_equipment"
)

func (r *Repository) Create(ctx context.Context, e *property_equipment.PropertyEquipment) error {
	query := `
		INSERT INTO property_equipment (
			property_id, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		e.PropertyID,
		e.Area,
		e.OutdoorBrand, e.OutdoorModel, e.OutdoorSerial, e.OutdoorInstalled,
		e.FurnaceBrand, e.FurnaceModel, e.FurnaceSerial, e.FurnaceInstalled,
		e.EvaporatorBrand, e.EvaporatorModel, e.EvaporatorSerial, e.EvaporatorInstalled,
		e.AirHandlerBrand, e.AirHandlerModel, e.AirHandlerSerial, e.AirHandlerInstalled,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert property_equipment query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	e.ID = id
	return nil
}
