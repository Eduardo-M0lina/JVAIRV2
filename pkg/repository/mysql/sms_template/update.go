package sms_template

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/sms_template"
)

func (r *Repository) Update(ctx context.Context, template *sms_template.SMSTemplate) error {
	query := `
		UPDATE sms_templates
		SET label = ?, message = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		template.Label,
		template.Message,
		template.IsActive,
		template.ID,
	)
	return err
}
