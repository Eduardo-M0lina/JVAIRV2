package email_template

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/angumol/jvairv2/pkg/domain/email_template"
)

func (r *Repository) GetByLabel(ctx context.Context, label string) (*email_template.EmailTemplate, error) {
	query := `
		SELECT id, label, subject, body, is_active, created_at, updated_at
		FROM email_templates
		WHERE label = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var template email_template.EmailTemplate
	err := r.db.QueryRowContext(ctx, query, label).Scan(
		&template.ID,
		&template.Label,
		&template.Subject,
		&template.Body,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("email template with label '%s' not found", label)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying email template: %w", err)
	}

	return &template, nil
}
