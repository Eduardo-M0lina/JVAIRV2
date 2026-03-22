package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/jvairv2/pkg/domain/password_history"
	"github.com/your-org/jvairv2/pkg/domain/password_reset"
	"github.com/your-org/jvairv2/pkg/domain/settings"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

var (
	ErrPasswordExpired        = errors.New("la contraseña ha expirado, debe cambiarla")
	ErrPasswordReuse          = errors.New("la nueva contraseña no puede ser igual a una contraseña anterior")
	ErrPasswordSameAsCurrent  = errors.New("la nueva contraseña no puede ser igual a la actual")
	ErrInvalidCurrentPassword = errors.New("la contraseña actual es incorrecta")
	ErrInvalidResetToken      = errors.New("token de reseteo inválido o expirado")
	ErrPasswordTooShort       = errors.New("la contraseña no cumple con la longitud mínima")
	ErrPasswordNoNumbers      = errors.New("la contraseña debe contener al menos un número")
	ErrPasswordNoSymbols      = errors.New("la contraseña debe contener al menos un símbolo")
	ErrTokenExpired           = errors.New("el token de reseteo ha expirado")
	ErrPasswordMismatch       = errors.New("las contraseñas no coinciden")
)

// PasswordResetRequest representa la solicitud de reseteo de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest representa la solicitud para resetear contraseña con token
type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

// ChangePasswordRequest representa la solicitud de cambio de contraseña
type ChangePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" validate:"required"`
}

// PasswordSecurityUseCase maneja la lógica de seguridad de contraseñas
type PasswordSecurityUseCase struct {
	userRepo            user.Repository
	passwordResetRepo   password_reset.Repository
	passwordHistoryRepo password_history.Repository
	settingsRepo        settings.Repository
	emailService        EmailService
}

// EmailService define la interfaz para envío de emails
type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, toEmail, resetLink string) error
}

// NewPasswordSecurityUseCase crea una nueva instancia del caso de uso de seguridad de contraseñas
func NewPasswordSecurityUseCase(
	userRepo user.Repository,
	passwordResetRepo password_reset.Repository,
	passwordHistoryRepo password_history.Repository,
	settingsRepo settings.Repository,
	emailService EmailService,
) *PasswordSecurityUseCase {
	return &PasswordSecurityUseCase{
		userRepo:            userRepo,
		passwordResetRepo:   passwordResetRepo,
		passwordHistoryRepo: passwordHistoryRepo,
		settingsRepo:        settingsRepo,
		emailService:        emailService,
	}
}

// ForgotPassword genera un token y envía email de reseteo
func (uc *PasswordSecurityUseCase) ForgotPassword(ctx context.Context, email string) error {
	// Verificar que el usuario existe
	u, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Por seguridad, no revelar si el email existe o no
		return nil
	}

	// Generar token único
	token := generateToken()

	// Eliminar tokens previos del usuario
	_ = uc.passwordResetRepo.DeleteByEmail(ctx, email)

	// Crear nuevo token
	pr := &password_reset.PasswordReset{
		Email:     email,
		Token:     token,
		CreatedAt: &[]time.Time{time.Now()}[0],
	}

	if err := uc.passwordResetRepo.Create(ctx, pr); err != nil {
		return fmt.Errorf("error al crear token de reseteo: %w", err)
	}

	// Enviar email con el link de reseteo
	// TODO: Obtener la URL base desde configuración
	resetLink := fmt.Sprintf("https://app.jvair.com/reset-password?token=%s&email=%s", token, email)
	if err := uc.emailService.SendPasswordResetEmail(ctx, u.Email, resetLink); err != nil {
		// Log el error pero no fallar la operación
		// El token ya fue creado y el usuario puede intentar de nuevo
		return fmt.Errorf("error al enviar email de reseteo: %w", err)
	}

	return nil
}

// ResetPassword resetea la contraseña usando un token
func (uc *PasswordSecurityUseCase) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// Validar que las contraseñas coincidan
	if req.Password != req.PasswordConfirm {
		return ErrPasswordMismatch
	}

	// Obtener configuraciones
	cfg, err := uc.settingsRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("error al obtener configuraciones: %w", err)
	}

	// Validar políticas de contraseña
	if err := uc.validatePasswordPolicy(req.Password, cfg); err != nil {
		return err
	}

	// Verificar token
	pr, err := uc.passwordResetRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return ErrInvalidResetToken
	}
	if pr == nil || pr.Email != req.Email {
		return ErrInvalidResetToken
	}

	// Verificar que el token no haya expirado (24 horas por defecto)
	if time.Since(*pr.CreatedAt) > 24*time.Hour {
		return ErrTokenExpired
	}

	// Obtener usuario
	u, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return ErrInvalidResetToken
	}

	// Verificar historial de contraseñas
	if err := uc.checkPasswordHistory(ctx, u.ID, req.Password, cfg); err != nil {
		return err
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}

	// Actualizar contraseña
	u.Password = string(hashedPassword)
	u.IsChangePassword = false
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	// Guardar en historial
	if err := uc.savePasswordHistory(ctx, u.ID, string(hashedPassword), cfg); err != nil {
		// Log error pero no fallar
		_ = err
	}

	// Eliminar token usado
	_ = uc.passwordResetRepo.Delete(ctx, req.Token)

	return nil
}

// ChangePassword cambia la contraseña del usuario autenticado
func (uc *PasswordSecurityUseCase) ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error {
	// Validar que las contraseñas coincidan
	if req.NewPassword != req.NewPasswordConfirm {
		return ErrPasswordMismatch
	}

	// Obtener usuario
	u, err := uc.userRepo.GetByID(ctx, fmt.Sprintf("%d", userID))
	if err != nil {
		return user.ErrUserNotFound
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidCurrentPassword
	}

	// Verificar que la nueva contraseña no sea igual a la actual
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.NewPassword)); err == nil {
		return ErrPasswordSameAsCurrent
	}

	// Obtener configuraciones
	cfg, err := uc.settingsRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("error al obtener configuraciones: %w", err)
	}

	// Validar políticas de contraseña
	if err := uc.validatePasswordPolicy(req.NewPassword, cfg); err != nil {
		return err
	}

	// Verificar historial de contraseñas
	if err := uc.checkPasswordHistory(ctx, userID, req.NewPassword, cfg); err != nil {
		return err
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}

	// Actualizar contraseña
	u.Password = string(hashedPassword)
	u.IsChangePassword = false
	if err := uc.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	// Guardar en historial
	if err := uc.savePasswordHistory(ctx, userID, string(hashedPassword), cfg); err != nil {
		// Log error pero no fallar
		_ = err
	}

	return nil
}

// CheckPasswordExpiry verifica si la contraseña del usuario ha expirado
func (uc *PasswordSecurityUseCase) CheckPasswordExpiry(ctx context.Context, userID int64) (bool, error) {
	// Obtener configuraciones
	cfg, err := uc.settingsRepo.Get(ctx)
	if err != nil {
		return false, err
	}

	// Si no se aplica reset obligatorio, no está expirada
	if !cfg.IsEnforceRoutinePasswordReset || cfg.PasswordExpireDays <= 0 {
		return false, nil
	}

	// Obtener historial más reciente
	histories, err := uc.passwordHistoryRepo.GetByUserID(ctx, userID, 1)
	if err != nil {
		return false, err
	}

	// Si no hay historial, usar la fecha de creación del usuario
	if len(histories) == 0 {
		u, err := uc.userRepo.GetByID(ctx, fmt.Sprintf("%d", userID))
		if err != nil {
			return false, err
		}
		if u.CreatedAt == nil {
			return false, nil
		}
		return time.Since(*u.CreatedAt) > time.Duration(cfg.PasswordExpireDays)*24*time.Hour, nil
	}

	// Verificar si la contraseña más reciente ha expirado
	return time.Since(*histories[0].CreatedAt) > time.Duration(cfg.PasswordExpireDays)*24*time.Hour, nil
}

// validatePasswordPolicy valida que la contraseña cumpla con las políticas
func (uc *PasswordSecurityUseCase) validatePasswordPolicy(password string, cfg *settings.Settings) error {
	// Validar longitud mínima
	if len(password) < cfg.PasswordMinimumLength {
		return ErrPasswordTooShort
	}

	// Validar números si es requerido
	if cfg.PasswordIncludeNumbers {
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		if !hasNumber {
			return ErrPasswordNoNumbers
		}
	}

	// Validar símbolos si es requerido
	if cfg.PasswordIncludeSymbols {
		hasSymbol := regexp.MustCompile(`[@$!%*#?&]`).MatchString(password)
		if !hasSymbol {
			return ErrPasswordNoSymbols
		}
	}

	return nil
}

// checkPasswordHistory verifica que la contraseña no esté en el historial
func (uc *PasswordSecurityUseCase) checkPasswordHistory(ctx context.Context, userID int64, password string, cfg *settings.Settings) error {
	// Obtener historial (limitado por password_age)
	limit := cfg.PasswordAge
	if limit <= 0 {
		limit = 5 // valor por defecto
	}

	histories, err := uc.passwordHistoryRepo.GetByUserID(ctx, userID, limit)
	if err != nil {
		return err
	}

	// Verificar contra cada contraseña en el historial
	for _, h := range histories {
		if err := bcrypt.CompareHashAndPassword([]byte(h.Password), []byte(password)); err == nil {
			return ErrPasswordReuse
		}
	}

	return nil
}

// savePasswordHistory guarda la contraseña en el historial
func (uc *PasswordSecurityUseCase) savePasswordHistory(ctx context.Context, userID int64, hashedPassword string, cfg *settings.Settings) error {
	// Crear nuevo registro
	ph := &password_history.PasswordHistory{
		UserID:   userID,
		Password: hashedPassword,
	}

	if err := uc.passwordHistoryRepo.Create(ctx, ph); err != nil {
		return err
	}

	// Mantener solo el número configurado de contraseñas en historial
	if cfg.PasswordHistoryCount > 0 {
		_ = uc.passwordHistoryRepo.DeleteOldest(ctx, userID, cfg.PasswordHistoryCount)
	}

	return nil
}

// generateToken genera un token único para reseteo
func generateToken() string {
	// Generar token aleatorio de 64 caracteres
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(65 + (i % 26)) // A-Z para simplificar
	}
	return fmt.Sprintf("%x", b)
}
