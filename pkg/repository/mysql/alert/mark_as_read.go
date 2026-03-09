package alert

import "context"

func (r *repository) MarkAsRead(ctx context.Context, id int64) error {
	query := `UPDATE alerts SET is_read = 1, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
