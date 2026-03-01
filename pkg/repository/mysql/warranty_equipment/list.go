package warranty_equipment

import (
	"context"
	"database/sql"
	"log/slog"

	domainWE "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
)

func (r *Repository) ListByWarrantyID(ctx context.Context, warrantyID int64) ([]*domainWE.WarrantyEquipment, error) {
	query := `
		SELECT id, warranty_id, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at
		FROM warranty_equipment
		WHERE warranty_id = ?
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, warrantyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty equipment",
			slog.String("error", err.Error()),
			slog.Int64("warrantyId", warrantyID))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var equipment []*domainWE.WarrantyEquipment
	for rows.Next() {
		we := &domainWE.WarrantyEquipment{}
		var outdoorBrand, outdoorModel, outdoorSerial sql.NullString
		var furnaceBrand, furnaceModel, furnaceSerial sql.NullString
		var evaporatorBrand, evaporatorModel, evaporatorSerial sql.NullString
		var airHandlerBrand, airHandlerModel, airHandlerSerial sql.NullString
		var outdoorInstalled, furnaceInstalled, evaporatorInstalled, airHandlerInstalled sql.NullTime

		if err := rows.Scan(
			&we.ID, &we.WarrantyID, &we.Area,
			&outdoorBrand, &outdoorModel, &outdoorSerial, &outdoorInstalled,
			&furnaceBrand, &furnaceModel, &furnaceSerial, &furnaceInstalled,
			&evaporatorBrand, &evaporatorModel, &evaporatorSerial, &evaporatorInstalled,
			&airHandlerBrand, &airHandlerModel, &airHandlerSerial, &airHandlerInstalled,
			&we.CreatedAt, &we.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty equipment row",
				slog.String("error", err.Error()))
			return nil, err
		}

		assignNullString(&we.OutdoorBrand, outdoorBrand)
		assignNullString(&we.OutdoorModel, outdoorModel)
		assignNullString(&we.OutdoorSerial, outdoorSerial)
		assignNullTime(&we.OutdoorInstalled, outdoorInstalled)
		assignNullString(&we.FurnaceBrand, furnaceBrand)
		assignNullString(&we.FurnaceModel, furnaceModel)
		assignNullString(&we.FurnaceSerial, furnaceSerial)
		assignNullTime(&we.FurnaceInstalled, furnaceInstalled)
		assignNullString(&we.EvaporatorBrand, evaporatorBrand)
		assignNullString(&we.EvaporatorModel, evaporatorModel)
		assignNullString(&we.EvaporatorSerial, evaporatorSerial)
		assignNullTime(&we.EvaporatorInstalled, evaporatorInstalled)
		assignNullString(&we.AirHandlerBrand, airHandlerBrand)
		assignNullString(&we.AirHandlerModel, airHandlerModel)
		assignNullString(&we.AirHandlerSerial, airHandlerSerial)
		assignNullTime(&we.AirHandlerInstalled, airHandlerInstalled)

		equipment = append(equipment, we)
	}

	return equipment, nil
}
