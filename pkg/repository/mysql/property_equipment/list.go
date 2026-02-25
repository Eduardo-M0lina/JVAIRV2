package property_equipment

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/property_equipment"
)

func (r *Repository) List(ctx context.Context, propertyID int64) ([]*property_equipment.PropertyEquipment, error) {
	query := `
		SELECT
			id, property_id, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at
		FROM property_equipment
		WHERE property_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, propertyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query property equipment",
			slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", closeErr.Error()))
		}
	}()

	var equipment []*property_equipment.PropertyEquipment
	for rows.Next() {
		e := &property_equipment.PropertyEquipment{}
		err := rows.Scan(
			&e.ID, &e.PropertyID, &e.Area,
			&e.OutdoorBrand, &e.OutdoorModel, &e.OutdoorSerial, &e.OutdoorInstalled,
			&e.FurnaceBrand, &e.FurnaceModel, &e.FurnaceSerial, &e.FurnaceInstalled,
			&e.EvaporatorBrand, &e.EvaporatorModel, &e.EvaporatorSerial, &e.EvaporatorInstalled,
			&e.AirHandlerBrand, &e.AirHandlerModel, &e.AirHandlerSerial, &e.AirHandlerInstalled,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan property equipment row",
				slog.String("error", err.Error()))
			return nil, err
		}
		equipment = append(equipment, e)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating property equipment rows",
			slog.String("error", err.Error()))
		return nil, err
	}

	return equipment, nil
}
