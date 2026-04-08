package email_template

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/email_template"
)

func (r *Repository) Create(ctx context.Context, template *email_template.EmailTemplate) error {
	query := `
		INSERT INTO email_templates (label, subject, body, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		template.Label,
		template.Subject,
		template.Body,
		template.IsActive,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create email template",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	template.ID = id
	return nil
}
