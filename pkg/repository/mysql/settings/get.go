package settings

import (
	"context"
	"database/sql"

	domainSettings "github.com/your-org/jvairv2/pkg/domain/settings"
)

// Get obtiene las configuraciones del sistema (siempre retorna el primer registro)
func (r *Repository) Get(ctx context.Context) (*domainSettings.Settings, error) {
	query := `
		SELECT id, is_twilio_enabled, twilio_sid, twilio_auth_token, twilio_from_number,
		       is_enforce_routine_password_reset, password_expire_days, password_history_count,
		       password_minimum_length, password_age, password_include_numbers,
		       password_include_symbols, created_at, updated_at
		FROM settings
		LIMIT 1
	`

	var settings domainSettings.Settings
	var twilioSID, twilioAuthToken, twilioFromNumber sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query).Scan(
		&settings.ID,
		&settings.IsTwilioEnabled,
		&twilioSID,
		&twilioAuthToken,
		&twilioFromNumber,
		&settings.IsEnforceRoutinePasswordReset,
		&settings.PasswordExpireDays,
		&settings.PasswordHistoryCount,
		&settings.PasswordMinimumLength,
		&settings.PasswordAge,
		&settings.PasswordIncludeNumbers,
		&settings.PasswordIncludeSymbols,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainSettings.ErrSettingsNotFound
		}
		return nil, err
	}

	// Asignar valores opcionales
	if twilioSID.Valid {
		settings.TwilioSID = &twilioSID.String
	}
	if twilioAuthToken.Valid {
		settings.TwilioAuthToken = &twilioAuthToken.String
	}
	if twilioFromNumber.Valid {
		settings.TwilioFromNumber = &twilioFromNumber.String
	}
	if createdAt.Valid {
		settings.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		settings.UpdatedAt = &updatedAt.Time
	}

	return &settings, nil
}
