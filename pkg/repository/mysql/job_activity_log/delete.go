package job_activity_log

import "context"

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM job_activity_logs WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
