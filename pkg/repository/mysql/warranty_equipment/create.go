package warranty_equipment

import (
	"context"
	"log/slog"

	domainWE "github.com/angumol/jvairv2/pkg/domain/warranty_equipment"
)

func (r *Repository) Create(ctx context.Context, we *domainWE.WarrantyEquipment) error {
	query := `
		INSERT INTO warranty_equipment (warranty_id, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		we.WarrantyID,
		we.Area,
		we.OutdoorBrand, we.OutdoorModel, we.OutdoorSerial, we.OutdoorInstalled,
		we.FurnaceBrand, we.FurnaceModel, we.FurnaceSerial, we.FurnaceInstalled,
		we.EvaporatorBrand, we.EvaporatorModel, we.EvaporatorSerial, we.EvaporatorInstalled,
		we.AirHandlerBrand, we.AirHandlerModel, we.AirHandlerSerial, we.AirHandlerInstalled,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty equipment query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	we.ID = id
	return nil
}
