package alert

import "context"

func (r *repository) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE alerts SET is_read = 1, updated_at = NOW() WHERE user_id = ? AND is_read = 0`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
