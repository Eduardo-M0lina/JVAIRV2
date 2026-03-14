package email_template

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/email_template"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*email_template.EmailTemplate, error) {
	query := `SELECT id, label, subject, body, is_active, created_at, updated_at FROM email_templates WHERE id = ?`

	item := &email_template.EmailTemplate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Label,
		&item.Subject,
		&item.Body,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, email_template.ErrEmailTemplateNotFound
		}
		return nil, err
	}

	return item, nil
}
