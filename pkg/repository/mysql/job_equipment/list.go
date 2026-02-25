package job_equipment

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

func (r *Repository) List(ctx context.Context, jobID int64, equipmentType string) ([]*job_equipment.JobEquipment, error) {
	query := `
		SELECT
			id, job_id, type, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at
		FROM job_equipment
		WHERE job_id = ?
	`

	var args []interface{}
	args = append(args, jobID)

	if equipmentType != "" {
		query += " AND type = ?"
		args = append(args, equipmentType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query job equipment",
			slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", closeErr.Error()))
		}
	}()

	var equipment []*job_equipment.JobEquipment
	for rows.Next() {
		e := &job_equipment.JobEquipment{}
		err := rows.Scan(
			&e.ID, &e.JobID, &e.Type, &e.Area,
			&e.OutdoorBrand, &e.OutdoorModel, &e.OutdoorSerial, &e.OutdoorInstalled,
			&e.FurnaceBrand, &e.FurnaceModel, &e.FurnaceSerial, &e.FurnaceInstalled,
			&e.EvaporatorBrand, &e.EvaporatorModel, &e.EvaporatorSerial, &e.EvaporatorInstalled,
			&e.AirHandlerBrand, &e.AirHandlerModel, &e.AirHandlerSerial, &e.AirHandlerInstalled,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan job equipment row",
				slog.String("error", err.Error()))
			return nil, err
		}
		equipment = append(equipment, e)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating job equipment rows",
			slog.String("error", err.Error()))
		return nil, err
	}

	return equipment, nil
}
