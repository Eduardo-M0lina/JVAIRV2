package account

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/jvairv2/pkg/domain/password_history"
	"github.com/your-org/jvairv2/pkg/domain/settings"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

var (
	ErrInvalidCurrentPassword = errors.New("contraseña actual incorrecta")
	ErrPasswordInHistory      = errors.New("la contraseña ya fue utilizada anteriormente")
	ErrPasswordTooWeak        = errors.New("la contraseña no cumple con los requisitos de seguridad")
	ErrEmailAlreadyInUse      = errors.New("el email ya está en uso por otro usuario")
)

// UseCase define los casos de uso para la gestión de cuenta
type UseCase struct {
	userRepo            user.Repository
	passwordHistoryRepo password_history.Repository
	settingsRepo        settings.Repository
}

// NewUseCase crea una nueva instancia del caso de uso de account
func NewUseCase(
	userRepo user.Repository,
	passwordHistoryRepo password_history.Repository,
	settingsRepo settings.Repository,
) *UseCase {
	return &UseCase{
		userRepo:            userRepo,
		passwordHistoryRepo: passwordHistoryRepo,
		settingsRepo:        settingsRepo,
	}
}

// GetProfile obtiene el perfil del usuario actual
func (uc *UseCase) GetProfile(ctx context.Context, userID int64) (*ProfileResponse, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	u, err := uc.userRepo.GetByID(ctx, userIDStr)
	if err != nil {
		slog.Error("Error al obtener usuario",
			"user_id", userID,
			"error", err,
		)
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return &ProfileResponse{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		RoleID:           u.RoleID,
		RoleName:         u.RoleName,
		RoleTitle:        u.RoleTitle,
		IsActive:         u.IsActive,
		IsChangePassword: u.IsChangePassword,
	}, nil
}

// UpdateProfile actualiza el perfil del usuario actual
func (uc *UseCase) UpdateProfile(ctx context.Context, userID int64, req *UpdateProfileRequest) error {
	userIDStr := strconv.FormatInt(userID, 10)
	u, err := uc.userRepo.GetByID(ctx, userIDStr)
	if err != nil {
		slog.Error("Error al obtener usuario",
			"user_id", userID,
			"error", err,
		)
		return fmt.Errorf("error al obtener usuario: %w", err)
	}

	// Verificar si el email cambió y si ya está en uso por otro usuario
	if req.Email != u.Email {
		existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			slog.Warn("Intento de actualizar perfil con email duplicado",
				"user_id", userID,
				"email", req.Email,
			)
			return ErrEmailAlreadyInUse
		}
		// Si el error no es "usuario no encontrado", propagarlo
		if err != nil && !errors.Is(err, user.ErrUserNotFound) {
			slog.Error("Error al verificar email existente",
				"email", req.Email,
				"error", err,
			)
			return fmt.Errorf("error al verificar email: %w", err)
		}
	}

	// Actualizar campos
	u.Name = req.Name
	u.Email = req.Email
	now := time.Now()
	u.UpdatedAt = &now

	// Guardar cambios
	err = uc.userRepo.Update(ctx, u)
	if err != nil {
		slog.Error("Error al actualizar usuario",
			"user_id", userID,
			"error", err,
		)
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	slog.Info("Perfil actualizado exitosamente",
		"user_id", userID,
		"name", req.Name,
		"email", req.Email,
	)

	return nil
}

// ChangePassword cambia la contraseña del usuario actual
func (uc *UseCase) ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error {
	userIDStr := strconv.FormatInt(userID, 10)
	u, err := uc.userRepo.GetByID(ctx, userIDStr)
	if err != nil {
		slog.Error("Error al obtener usuario",
			"user_id", userID,
			"error", err,
		)
		return fmt.Errorf("error al obtener usuario: %w", err)
	}

	// Verificar contraseña actual
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.CurrentPassword))
	if err != nil {
		slog.Warn("Intento de cambio de contraseña con contraseña actual incorrecta",
			"user_id", userID,
		)
		return ErrInvalidCurrentPassword
	}

	// Obtener configuración de seguridad de contraseñas
	settingsData, err := uc.settingsRepo.Get(ctx)
	if err != nil {
		slog.Error("Error al obtener configuración de seguridad",
			"error", err,
		)
		return fmt.Errorf("error al obtener configuración: %w", err)
	}

	// Validar requisitos de contraseña
	if err := uc.validatePassword(req.NewPassword, settingsData); err != nil {
		slog.Warn("Contraseña no cumple con los requisitos",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	// Verificar historial de contraseñas
	if settingsData.PasswordHistoryCount > 0 {
		history, err := uc.passwordHistoryRepo.GetByUserID(ctx, userID, settingsData.PasswordHistoryCount)
		if err != nil {
			slog.Error("Error al obtener historial de contraseñas",
				"user_id", userID,
				"error", err,
			)
			return fmt.Errorf("error al obtener historial de contraseñas: %w", err)
		}

		// Verificar si la nueva contraseña ya fue usada
		for _, ph := range history {
			err := bcrypt.CompareHashAndPassword([]byte(ph.Password), []byte(req.NewPassword))
			if err == nil {
				slog.Warn("Intento de reutilizar contraseña anterior",
					"user_id", userID,
				)
				return ErrPasswordInHistory
			}
		}
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error al generar hash de contraseña",
			"error", err,
		)
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}

	// Guardar en historial de contraseñas
	now := time.Now()
	passwordHistory := &password_history.PasswordHistory{
		UserID:    userID,
		Password:  string(hashedPassword),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	err = uc.passwordHistoryRepo.Create(ctx, passwordHistory)
	if err != nil {
		slog.Error("Error al guardar historial de contraseña",
			"user_id", userID,
			"error", err,
		)
		// No fallar si no se puede guardar el historial
	}

	// Limpiar historial antiguo si excede el límite
	if settingsData.PasswordHistoryCount > 0 {
		err = uc.passwordHistoryRepo.DeleteOldest(ctx, userID, settingsData.PasswordHistoryCount)
		if err != nil {
			slog.Warn("Error al limpiar historial antiguo de contraseñas",
				"user_id", userID,
				"error", err,
			)
			// No fallar si no se puede limpiar el historial
		}
	}

	// Actualizar contraseña del usuario
	u.Password = string(hashedPassword)
	u.IsChangePassword = false // Resetear flag de cambio forzado
	u.UpdatedAt = &now

	err = uc.userRepo.Update(ctx, u)
	if err != nil {
		slog.Error("Error al actualizar contraseña del usuario",
			"user_id", userID,
			"error", err,
		)
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	slog.Info("Contraseña cambiada exitosamente",
		"user_id", userID,
	)

	return nil
}

// validatePassword valida que la contraseña cumpla con los requisitos de seguridad
func (uc *UseCase) validatePassword(password string, settings *settings.Settings) error {
	// Validar longitud mínima
	if len(password) < settings.PasswordMinimumLength {
		return fmt.Errorf("%w: debe tener al menos %d caracteres", ErrPasswordTooWeak, settings.PasswordMinimumLength)
	}

	// Validar que incluya números si es requerido
	if settings.PasswordIncludeNumbers {
		hasNumber := false
		for _, char := range password {
			if char >= '0' && char <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return fmt.Errorf("%w: debe incluir al menos un número", ErrPasswordTooWeak)
		}
	}

	// Validar que incluya símbolos si es requerido
	if settings.PasswordIncludeSymbols {
		hasSymbol := false
		symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
		for _, char := range password {
			for _, symbol := range symbols {
				if char == symbol {
					hasSymbol = true
					break
				}
			}
			if hasSymbol {
				break
			}
		}
		if !hasSymbol {
			return fmt.Errorf("%w: debe incluir al menos un símbolo especial", ErrPasswordTooWeak)
		}
	}

	return nil
}
