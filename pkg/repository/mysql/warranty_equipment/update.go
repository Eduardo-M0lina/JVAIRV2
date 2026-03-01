package warranty_equipment

import (
	"context"
	"log/slog"

	domainWE "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
)

func (r *Repository) Update(ctx context.Context, we *domainWE.WarrantyEquipment) error {
	query := `
		UPDATE warranty_equipment
		SET area = ?,
			outdoor_brand = ?, outdoor_model = ?, outdoor_serial = ?, outdoor_installed = ?,
			furnace_brand = ?, furnace_model = ?, furnace_serial = ?, furnace_installed = ?,
			evaporator_brand = ?, evaporator_model = ?, evaporator_serial = ?, evaporator_installed = ?,
			air_handler_brand = ?, air_handler_model = ?, air_handler_serial = ?, air_handler_installed = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		we.Area,
		we.OutdoorBrand, we.OutdoorModel, we.OutdoorSerial, we.OutdoorInstalled,
		we.FurnaceBrand, we.FurnaceModel, we.FurnaceSerial, we.FurnaceInstalled,
		we.EvaporatorBrand, we.EvaporatorModel, we.EvaporatorSerial, we.EvaporatorInstalled,
		we.AirHandlerBrand, we.AirHandlerModel, we.AirHandlerSerial, we.AirHandlerInstalled,
		we.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty equipment",
			slog.String("error", err.Error()),
			slog.Int64("id", we.ID))
		return err
	}

	return nil
}
