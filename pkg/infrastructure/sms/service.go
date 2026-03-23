package sms

import (
	"context"
	"log/slog"

	domainSettings "github.com/your-org/jvairv2/pkg/domain/settings"
)

// SettingsProvider define la interfaz para obtener configuraciones del sistema
type SettingsProvider interface {
	Get(ctx context.Context) (*domainSettings.Settings, error)
}

// SelectiveSender enruta SMS a AWS SNS o Twilio según la configuración en BD
type SelectiveSender struct {
	snsSender        SMSSender
	settingsProvider SettingsProvider
}

// NewSelectiveSender crea un sender que elige entre SNS y Twilio en tiempo de ejecución
func NewSelectiveSender(snsSender SMSSender, settingsProvider SettingsProvider) *SelectiveSender {
	return &SelectiveSender{
		snsSender:        snsSender,
		settingsProvider: settingsProvider,
	}
}

// SendSMS envía un SMS usando Twilio si está habilitado en settings, o AWS SNS como fallback
func (s *SelectiveSender) SendSMS(ctx context.Context, to string, body string) error {
	cfg, err := s.settingsProvider.Get(ctx)
	if err != nil {
		slog.WarnContext(ctx, "[SelectiveSender] no se pudo leer settings, usando AWS SNS",
			slog.String("error", err.Error()),
		)
	} else if cfg.IsTwilioEnabled &&
		cfg.TwilioSID != nil &&
		cfg.TwilioAuthToken != nil &&
		cfg.TwilioFromNumber != nil {
		slog.InfoContext(ctx, "[SelectiveSender] usando Twilio", slog.String("to", to))
		twilio := NewTwilioSender(*cfg.TwilioSID, *cfg.TwilioAuthToken, *cfg.TwilioFromNumber)
		return twilio.SendSMS(ctx, to, body)
	} else {
		slog.InfoContext(ctx, "[SelectiveSender] Twilio deshabilitado o sin credenciales, usando AWS SNS",
			slog.Bool("isTwilioEnabled", cfg.IsTwilioEnabled),
			slog.Bool("tieneSID", cfg.TwilioSID != nil),
			slog.Bool("tieneToken", cfg.TwilioAuthToken != nil),
			slog.Bool("tieneFrom", cfg.TwilioFromNumber != nil),
		)
	}
	return s.snsSender.SendSMS(ctx, to, body)
}
