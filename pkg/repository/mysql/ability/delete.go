package ability

import (
	"context"
	"database/sql"
)

// Delete elimina una ability por su ID
func (r *Repository) Delete(ctx context.Context, id int64) error {
	// Verificar si la ability existe
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Preparar la consulta
	query := `DELETE FROM abilities WHERE id = ?`

	// Ejecutar la consulta
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Verificar que se haya eliminado al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
