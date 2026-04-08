package account

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/angumol/jvairv2/pkg/domain/account"
	domainUser "github.com/angumol/jvairv2/pkg/domain/user"
	"github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/angumol/jvairv2/pkg/rest/middleware"
)

// Handler maneja las solicitudes HTTP relacionadas con la gestión de cuenta
type Handler struct {
	accountUseCase *account.UseCase
}

// NewHandler crea una nueva instancia del handler de account
func NewHandler(accountUseCase *account.UseCase) *Handler {
	return &Handler{
		accountUseCase: accountUseCase,
	}
}

// GetProfile obtiene el perfil del usuario actual
// @Summary Obtener perfil del usuario actual
// @Description Obtiene la información del perfil del usuario autenticado
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} account.ProfileResponse
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/account [get]
// @Security BearerAuth
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		slog.Error("Usuario no encontrado en el contexto")
		handler.RespondWithError(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Obtener perfil
	profile, err := h.accountUseCase.GetProfile(r.Context(), userCtx.ID)
	if err != nil {
		slog.Error("Error al obtener perfil",
			"user_id", userCtx.ID,
			"error", err,
		)
		handler.RespondWithError(w, http.StatusInternalServerError, "Error al obtener perfil")
		return
	}

	slog.Info("Perfil obtenido exitosamente",
		"user_id", userCtx.ID,
	)

	handler.RespondWithJSON(w, http.StatusOK, profile)
}

// UpdateProfile actualiza el perfil del usuario actual
// @Summary Actualizar perfil del usuario actual
// @Description Actualiza el nombre y email del usuario autenticado
// @Tags Account
// @Accept json
// @Produce json
// @Param profile body account.UpdateProfileRequest true "Datos del perfil a actualizar"
// @Success 200 {object} map[string]string "Perfil actualizado exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 409 {object} map[string]string "Email ya está en uso"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/account [put]
// @Security BearerAuth
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		slog.Error("Usuario no encontrado en el contexto")
		handler.RespondWithError(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Decodificar request
	var req account.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Error al decodificar solicitud de actualización de perfil",
			"user_id", userCtx.ID,
			"error", err,
		)
		handler.RespondWithError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Validar request
	if req.Name == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "El nombre es requerido")
		return
	}
	if req.Email == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "El email es requerido")
		return
	}

	// Actualizar perfil
	err := h.accountUseCase.UpdateProfile(r.Context(), userCtx.ID, &req)
	if err != nil {
		if err == account.ErrEmailAlreadyInUse {
			slog.Warn("Intento de actualizar perfil con email duplicado",
				"user_id", userCtx.ID,
				"email", req.Email,
			)
			handler.RespondWithError(w, http.StatusConflict, "El email ya está en uso")
			return
		}
		slog.Error("Error al actualizar perfil",
			"user_id", userCtx.ID,
			"error", err,
		)
		handler.RespondWithError(w, http.StatusInternalServerError, "Error al actualizar perfil")
		return
	}

	slog.Info("Perfil actualizado exitosamente",
		"user_id", userCtx.ID,
		"name", req.Name,
		"email", req.Email,
	)

	handler.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Perfil actualizado exitosamente",
	})
}

// ChangePassword cambia la contraseña del usuario actual
// @Summary Cambiar contraseña del usuario actual
// @Description Cambia la contraseña del usuario autenticado validando la contraseña actual y los requisitos de seguridad
// @Tags Account
// @Accept json
// @Produce json
// @Param password body account.ChangePasswordRequest true "Datos para cambiar la contraseña"
// @Success 200 {object} map[string]string "Contraseña cambiada exitosamente"
// @Failure 400 {object} map[string]string "Datos inválidos"
// @Failure 401 {object} map[string]string "Usuario no autenticado o contraseña actual incorrecta"
// @Failure 422 {object} map[string]string "La contraseña no cumple con los requisitos de seguridad"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/account/password [put]
// @Security BearerAuth
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto
	userCtx, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userCtx == nil {
		slog.Error("Usuario no encontrado en el contexto")
		handler.RespondWithError(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	// Decodificar request
	var req account.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Error al decodificar solicitud de cambio de contraseña",
			"user_id", userCtx.ID,
			"error", err,
		)
		handler.RespondWithError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Validar request
	if req.CurrentPassword == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "La contraseña actual es requerida")
		return
	}
	if req.NewPassword == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "La nueva contraseña es requerida")
		return
	}
	if req.ConfirmPassword == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "La confirmación de contraseña es requerida")
		return
	}
	if req.NewPassword != req.ConfirmPassword {
		handler.RespondWithError(w, http.StatusBadRequest, "Las contraseñas no coinciden")
		return
	}

	// Cambiar contraseña
	err := h.accountUseCase.ChangePassword(r.Context(), userCtx.ID, &req)
	if err != nil {
		if err == account.ErrInvalidCurrentPassword {
			slog.Warn("Intento de cambio de contraseña con contraseña actual incorrecta",
				"user_id", userCtx.ID,
			)
			handler.RespondWithError(w, http.StatusUnauthorized, "La contraseña actual es incorrecta")
			return
		}
		if err == account.ErrPasswordInHistory {
			slog.Warn("Intento de reutilizar contraseña anterior",
				"user_id", userCtx.ID,
			)
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña ya fue utilizada anteriormente")
			return
		}
		if err == account.ErrPasswordTooWeak {
			slog.Warn("Contraseña no cumple con los requisitos",
				"user_id", userCtx.ID,
			)
			handler.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		slog.Error("Error al cambiar contraseña",
			"user_id", userCtx.ID,
			"error", err,
		)
		handler.RespondWithError(w, http.StatusInternalServerError, "Error al cambiar contraseña")
		return
	}

	slog.Info("Contraseña cambiada exitosamente",
		"user_id", userCtx.ID,
	)

	handler.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Contraseña cambiada exitosamente",
	})
}
