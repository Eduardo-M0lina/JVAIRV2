package settings

import (
	"context"

	domainSettings "github.com/your-org/jvairv2/pkg/domain/settings"
)

// Update actualiza las configuraciones del sistema
func (r *Repository) Update(ctx context.Context, settings *domainSettings.Settings) error {
	query := `
		UPDATE settings
		SET is_twilio_enabled = ?,
		    twilio_sid = ?,
		    twilio_auth_token = ?,
		    twilio_from_number = ?,
		    is_enforce_routine_password_reset = ?,
		    password_expire_days = ?,
		    password_history_count = ?,
		    password_minimum_length = ?,
		    password_age = ?,
		    password_include_numbers = ?,
		    password_include_symbols = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.IsTwilioEnabled,
		settings.TwilioSID,
		settings.TwilioAuthToken,
		settings.TwilioFromNumber,
		settings.IsEnforceRoutinePasswordReset,
		settings.PasswordExpireDays,
		settings.PasswordHistoryCount,
		settings.PasswordMinimumLength,
		settings.PasswordAge,
		settings.PasswordIncludeNumbers,
		settings.PasswordIncludeSymbols,
		settings.ID,
	)

	return err
}
