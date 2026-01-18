package settings

import "time"

// Settings representa la configuración general del sistema
type Settings struct {
	ID                            int64      `json:"id"`
	IsTwilioEnabled               bool       `json:"isTwilioEnabled"`
	TwilioSID                     *string    `json:"twilioSid,omitempty"`
	TwilioAuthToken               *string    `json:"twilioAuthToken,omitempty"`
	TwilioFromNumber              *string    `json:"twilioFromNumber,omitempty"`
	IsEnforceRoutinePasswordReset bool       `json:"isEnforceRoutinePasswordReset"`
	PasswordExpireDays            int        `json:"passwordExpireDays"`
	PasswordHistoryCount          int        `json:"passwordHistoryCount"`
	PasswordMinimumLength         int        `json:"passwordMinimumLength"`
	PasswordAge                   int        `json:"passwordAge"`
	PasswordIncludeNumbers        bool       `json:"passwordIncludeNumbers"`
	PasswordIncludeSymbols        bool       `json:"passwordIncludeSymbols"`
	CreatedAt                     *time.Time `json:"createdAt,omitempty"`
	UpdatedAt                     *time.Time `json:"updatedAt,omitempty"`
}
