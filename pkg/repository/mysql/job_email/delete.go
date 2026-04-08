package job_email

import "context"

func (r *repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM job_emails WHERE id = ?`, id)
	return err
}
