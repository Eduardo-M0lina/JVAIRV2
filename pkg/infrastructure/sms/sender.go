package sms

import "context"

// SMSSender define la interfaz para enviar SMS
type SMSSender interface {
	SendSMS(ctx context.Context, to string, body string) error
}
