package job_equipment

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*job_equipment.JobEquipment, error) {
	query := `
		SELECT
			id, job_id, type, area,
			outdoor_brand, outdoor_model, outdoor_serial, outdoor_installed,
			furnace_brand, furnace_model, furnace_serial, furnace_installed,
			evaporator_brand, evaporator_model, evaporator_serial, evaporator_installed,
			air_handler_brand, air_handler_model, air_handler_serial, air_handler_installed,
			created_at, updated_at
		FROM job_equipment
		WHERE id = ?
	`

	e := &job_equipment.JobEquipment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.JobID, &e.Type, &e.Area,
		&e.OutdoorBrand, &e.OutdoorModel, &e.OutdoorSerial, &e.OutdoorInstalled,
		&e.FurnaceBrand, &e.FurnaceModel, &e.FurnaceSerial, &e.FurnaceInstalled,
		&e.EvaporatorBrand, &e.EvaporatorModel, &e.EvaporatorSerial, &e.EvaporatorInstalled,
		&e.AirHandlerBrand, &e.AirHandlerModel, &e.AirHandlerSerial, &e.AirHandlerInstalled,
		&e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Job equipment not found",
				slog.Int64("equipment_id", id))
			return nil, errors.New("job equipment not found")
		}
		slog.ErrorContext(ctx, "Failed to query job equipment by ID",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return nil, err
	}

	return e, nil
}
