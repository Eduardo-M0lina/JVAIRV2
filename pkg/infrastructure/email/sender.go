package email

import "context"

// EmailSender define la interfaz para enviar emails
type EmailSender interface {
	// SendEmail envía un email con los parámetros especificados
	SendEmail(ctx context.Context, params SendEmailParams) error
}

// SendEmailParams contiene los parámetros para enviar un email
type SendEmailParams struct {
	To          []string          // Destinatarios
	Subject     string            // Asunto
	HTMLBody    string            // Cuerpo HTML
	TextBody    string            // Cuerpo texto plano (opcional)
	Attachments []Attachment      // Archivos adjuntos (opcional)
	ReplyTo     string            // Email de respuesta (opcional)
	CC          []string          // Copia (opcional)
	BCC         []string          // Copia oculta (opcional)
	Headers     map[string]string // Headers personalizados (opcional)
}

// Attachment representa un archivo adjunto
type Attachment struct {
	Filename    string // Nombre del archivo
	ContentType string // Tipo MIME
	Data        []byte // Contenido del archivo
}
