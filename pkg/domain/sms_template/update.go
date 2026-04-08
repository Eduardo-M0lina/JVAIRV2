package sms_template

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, template *SMSTemplate) error {
	if err := template.Validate(); err != nil {
		return err
	}

	if _, err := uc.repo.GetByID(ctx, template.ID); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, template); err != nil {
		slog.ErrorContext(ctx, "Failed to update sms template",
			slog.String("error", err.Error()),
			slog.Int64("id", template.ID))
		return err
	}

	return nil
}
