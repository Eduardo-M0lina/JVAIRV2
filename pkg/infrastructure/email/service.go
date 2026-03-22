package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
)

// Service proporciona funcionalidades de alto nivel para envío de emails
type Service struct {
	sender   EmailSender
	renderer *TemplateRenderer
}

// TemplateRenderer renderiza templates estáticos
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer crea un nuevo renderer de templates a partir de un directorio
func NewTemplateRenderer(templatesDir string) (*TemplateRenderer, error) {
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error al parsear templates: %w", err)
	}
	return &TemplateRenderer{templates: templates}, nil
}

// Renderiza un template con los datos proporcionados
func (r *TemplateRenderer) Render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar template: %w", err)
	}

	return buf.String(), nil
}

// NewService crea una nueva instancia del servicio de email
func NewService(sender EmailSender, templatesDir string) (*Service, error) {
	renderer, err := NewTemplateRenderer(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar renderer: %w", err)
	}

	return &Service{
		sender:   sender,
		renderer: renderer,
	}, nil
}

// SendTemplatedEmail envía un email usando un template HTML estático
func (s *Service) SendTemplatedEmail(ctx context.Context, templateName string, data interface{}, params SendEmailParams) error {
	// Renderizar el template
	htmlBody, err := s.renderer.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("error al renderizar template '%s': %w", templateName, err)
	}

	// Actualizar params con el contenido renderizado
	params.HTMLBody = htmlBody

	// Enviar el email
	return s.sender.SendEmail(ctx, params)
}

// SendEmail envía un email sin template (directo)
func (s *Service) SendEmail(ctx context.Context, params SendEmailParams) error {
	return s.sender.SendEmail(ctx, params)
}
