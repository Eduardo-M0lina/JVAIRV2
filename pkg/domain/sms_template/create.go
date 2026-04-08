package sms_template

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, template *SMSTemplate) error {
	if err := template.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, template); err != nil {
		slog.ErrorContext(ctx, "Failed to create sms template",
			slog.String("error", err.Error()),
			slog.String("label", template.Label))
		return err
	}

	return nil
}
