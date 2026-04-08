package sms_template

import (
	"context"
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/sms_template"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*sms_template.SMSTemplate, error) {
	query := `SELECT id, label, message, is_active, created_at, updated_at FROM sms_templates WHERE id = ?`

	item := &sms_template.SMSTemplate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Label,
		&item.Message,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sms_template.ErrSMSTemplateNotFound
		}
		return nil, err
	}

	return item, nil
}
