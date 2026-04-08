package email

import (
	"context"
	"fmt"
	"time"

	"github.com/angumol/jvairv2/configs"
	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSender implementa EmailSender usando Mailgun
type MailgunSender struct {
	mg        *mailgun.MailgunImpl
	fromName  string
	fromEmail string
	domain    string
	timeout   time.Duration
}

// NewMailgunSender crea una nueva instancia de MailgunSender
func NewMailgunSender(config configs.MailConfig) *MailgunSender {
	mg := mailgun.NewMailgun(config.MailgunDomain, config.MailgunSecret)

	return &MailgunSender{
		mg:        mg,
		fromName:  config.FromName,
		fromEmail: config.FromAddress,
		domain:    config.MailgunDomain,
		timeout:   30 * time.Second,
	}
}

// SendEmail envía un email usando Mailgun
func (s *MailgunSender) SendEmail(ctx context.Context, params SendEmailParams) error {
	// Validar parámetros
	if len(params.To) == 0 {
		return fmt.Errorf("al menos un destinatario es requerido")
	}
	if params.Subject == "" {
		return fmt.Errorf("el asunto es requerido")
	}
	if params.HTMLBody == "" && params.TextBody == "" {
		return fmt.Errorf("el cuerpo del email (HTML o texto) es requerido")
	}

	// Crear mensaje
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	message := mailgun.NewMessage(from, params.Subject, params.TextBody)

	// Agregar destinatarios
	for _, to := range params.To {
		if err := message.AddRecipient(to); err != nil {
			return fmt.Errorf("error al agregar destinatario %s: %w", to, err)
		}
	}

	// Agregar CC
	for _, cc := range params.CC {
		message.AddCC(cc)
	}

	// Agregar BCC
	for _, bcc := range params.BCC {
		message.AddBCC(bcc)
	}

	// Agregar cuerpo HTML si existe
	if params.HTMLBody != "" {
		message.SetHTML(params.HTMLBody)
	}

	// Agregar ReplyTo si existe
	if params.ReplyTo != "" {
		message.SetReplyTo(params.ReplyTo)
	}

	// Agregar headers personalizados
	for key, value := range params.Headers {
		message.AddHeader(key, value)
	}

	// Agregar archivos adjuntos
	for _, attachment := range params.Attachments {
		message.AddBufferAttachment(attachment.Filename, attachment.Data)
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Enviar mensaje
	_, _, err := s.mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("error al enviar email via Mailgun: %w", err)
	}

	return nil
}
