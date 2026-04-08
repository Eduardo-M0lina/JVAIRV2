package sms_template

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/sms_template"
)

func (r *Repository) Create(ctx context.Context, template *sms_template.SMSTemplate) error {
	query := `
		INSERT INTO sms_templates (label, message, is_active, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		template.Label,
		template.Message,
		template.IsActive,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create sms template",
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
