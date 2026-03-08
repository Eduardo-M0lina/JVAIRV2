package alert

import "context"

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM alerts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
