package job_equipment

import (
	"database/sql"
	"log/slog"
)

// JobCheckerAdapter implementa job_equipment.JobChecker usando MySQL
type JobCheckerAdapter struct {
	db *sql.DB
}

// NewJobCheckerAdapter crea una nueva instancia del adapter
func NewJobCheckerAdapter(db *sql.DB) *JobCheckerAdapter {
	return &JobCheckerAdapter{db: db}
}

// JobExists verifica si un job existe y no está eliminado
func (a *JobCheckerAdapter) JobExists(id int64) (bool, error) {
	query := `SELECT COUNT(*) FROM jobs WHERE id = ? AND deleted_at IS NULL`

	var count int
	err := a.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		slog.Error("Failed to check job existence",
			slog.String("error", err.Error()),
			slog.Int64("job_id", id))
		return false, err
	}

	return count > 0, nil
}
