package sms_template

import "context"

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sms_templates WHERE id = ?`, id)
	return err
}
