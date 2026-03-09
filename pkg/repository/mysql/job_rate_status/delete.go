package job_rate_status

import "context"

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := "UPDATE job_rate_statuses SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
