package alert

import "context"

func (r *repository) UnreadCount(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM alerts WHERE user_id = ? AND is_read = 0`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
