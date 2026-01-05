package permission

import (
	"context"
	"database/sql"
)

// Delete elimina un permiso por su ID
func (r *Repository) Delete(ctx context.Context, id int64) error {
	// Verificar si el permiso existe
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Preparar la consulta
	query := `DELETE FROM permissions WHERE id = ?`

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
