package settings

import "time"

// Settings representa la configuración general del sistema
type Settings struct {
	ID                            int64      `json:"id"`
	IsTwilioEnabled               bool       `json:"is_twilio_enabled"`
	TwilioSID                     *string    `json:"twilio_sid,omitempty"`
	TwilioAuthToken               *string    `json:"twilio_auth_token,omitempty"`
	TwilioFromNumber              *string    `json:"twilio_from_number,omitempty"`
	IsEnforceRoutinePasswordReset bool       `json:"is_enforce_routine_password_reset"`
	PasswordExpireDays            int        `json:"password_expire_days"`
	PasswordHistoryCount          int        `json:"password_history_count"`
	PasswordMinimumLength         int        `json:"password_minimum_length"`
	PasswordAge                   int        `json:"password_age"`
	PasswordIncludeNumbers        bool       `json:"password_include_numbers"`
	PasswordIncludeSymbols        bool       `json:"password_include_symbols"`
	CreatedAt                     *time.Time `json:"created_at,omitempty"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty"`
}
