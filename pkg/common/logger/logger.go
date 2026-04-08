package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var defaultLogger *slog.Logger
var projectRoot string

// init se ejecuta automáticamente al importar el paquete
func init() {
	// Obtener la ruta del proyecto (JVAIRV2)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
}

// replaceSource reemplaza la ruta absoluta con una ruta relativa al proyecto
func replaceSource(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if source, ok := a.Value.Any().(*slog.Source); ok {
			// Reemplazar la ruta absoluta con la ruta relativa
			if strings.HasPrefix(source.File, projectRoot) {
				source.File = strings.TrimPrefix(source.File, projectRoot+string(filepath.Separator))
			}
		}
	}
	return a
}

// Init inicializa el logger global con configuración personalizada
func Init(env string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   true,
		ReplaceAttr: replaceSource,
	}

	// En desarrollo, usar formato de texto más legible
	if env == "development" || env == "dev" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// En producción, usar formato JSON
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// GetLogger retorna el logger global
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		Init("development")
	}
	return defaultLogger
}

// Info registra un mensaje informativo
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// InfoContext registra un mensaje informativo con contexto
func InfoContext(ctx context.Context, msg string, args ...any) {
	GetLogger().InfoContext(ctx, msg, args...)
}

// Warn registra un mensaje de advertencia
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// WarnContext registra un mensaje de advertencia con contexto
func WarnContext(ctx context.Context, msg string, args ...any) {
	GetLogger().WarnContext(ctx, msg, args...)
}

// Error registra un mensaje de error
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// ErrorContext registra un mensaje de error con contexto
func ErrorContext(ctx context.Context, msg string, args ...any) {
	GetLogger().ErrorContext(ctx, msg, args...)
}

// Debug registra un mensaje de depuración
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// DebugContext registra un mensaje de depuración con contexto
func DebugContext(ctx context.Context, msg string, args ...any) {
	GetLogger().DebugContext(ctx, msg, args...)
}

// With crea un nuevo logger con atributos adicionales
func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}

// WithGroup crea un nuevo logger con un grupo de atributos
func WithGroup(name string) *slog.Logger {
	return GetLogger().WithGroup(name)
}
