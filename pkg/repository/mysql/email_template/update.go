package email_template

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/email_template"
)

func (r *Repository) Update(ctx context.Context, template *email_template.EmailTemplate) error {
	query := `
		UPDATE email_templates
		SET label = ?, subject = ?, body = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		template.Label,
		template.Subject,
		template.Body,
		template.IsActive,
		template.ID,
	)
	return err
}
