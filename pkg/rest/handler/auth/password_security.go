package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/angumol/jvairv2/pkg/domain/auth"
	domainUser "github.com/angumol/jvairv2/pkg/domain/user"
	handler "github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/angumol/jvairv2/pkg/rest/middleware"
)

// PasswordSecurityHandler maneja las solicitudes HTTP relacionadas con seguridad de contraseñas
type PasswordSecurityHandler struct {
	passwordSecurityUseCase *auth.PasswordSecurityUseCase
}

// NewPasswordSecurityHandler crea una nueva instancia del handler de seguridad de contraseñas
func NewPasswordSecurityHandler(passwordSecurityUseCase *auth.PasswordSecurityUseCase) *PasswordSecurityHandler {
	return &PasswordSecurityHandler{
		passwordSecurityUseCase: passwordSecurityUseCase,
	}
}

// ForgotPassword maneja la solicitud de recuperación de contraseña
// @Summary Solicitar recuperación de contraseña
// @Description Envía un email con el link para resetear la contraseña
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.PasswordResetRequest true "Email del usuario"
// @Success 200 {string} string "Se ha enviado un email de recuperación"
// @Failure 400 {string} string "Error al decodificar la solicitud"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /auth/forgot-password [post]
func (h *PasswordSecurityHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.PasswordResetRequest

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Error al decodificar solicitud de forgot-password",
			"error", err,
			"remote_addr", r.RemoteAddr,
		)
		handler.RespondWithError(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar email
	if req.Email == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Email es requerido")
		return
	}

	// Procesar forgot password
	if err := h.passwordSecurityUseCase.ForgotPassword(r.Context(), req.Email); err != nil {
		slog.Error("Error al procesar forgot-password",
			"email", req.Email,
			"error", err,
		)
		// Por seguridad, no revelar si el email existe o no
	}

	// Siempre retornar éxito para no revelar si el email existe
	handler.RespondWithSuccess(w, "Si el email existe en nuestro sistema, recibirás instrucciones para recuperar tu contraseña")
}

// ResetPassword maneja la solicitud de reseteo de contraseña con token
// @Summary Resetear contraseña
// @Description Resetea la contraseña usando un token válido
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.ResetPasswordRequest true "Token y nueva contraseña"
// @Success 200 {string} string "Contraseña actualizada exitosamente"
// @Failure 400 {string} string "Error en la solicitud"
// @Failure 401 {string} string "Token inválido o expirado"
// @Failure 422 {string} string "La contraseña no cumple con las políticas"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /auth/reset-password [post]
func (h *PasswordSecurityHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ResetPasswordRequest

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar campos requeridos
	if req.Token == "" || req.Email == "" || req.Password == "" || req.PasswordConfirm == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Todos los campos son requeridos")
		return
	}

	// Procesar reset de contraseña
	if err := h.passwordSecurityUseCase.ResetPassword(r.Context(), &req); err != nil {
		switch err {
		case auth.ErrPasswordMismatch:
			handler.RespondWithError(w, http.StatusBadRequest, "Las contraseñas no coinciden")
		case auth.ErrInvalidResetToken:
			handler.RespondWithError(w, http.StatusUnauthorized, "Token inválido")
		case auth.ErrTokenExpired:
			handler.RespondWithError(w, http.StatusUnauthorized, "Token expirado")
		case auth.ErrPasswordTooShort:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña no cumple con la longitud mínima")
		case auth.ErrPasswordNoNumbers:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña debe contener al menos un número")
		case auth.ErrPasswordNoSymbols:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña debe contener al menos un símbolo")
		case auth.ErrPasswordReuse:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "No puedes reutilizar una contraseña anterior")
		default:
			slog.Error("Error al resetear contraseña",
				"error", err,
			)
			handler.RespondWithError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}

	handler.RespondWithSuccess(w, "Contraseña actualizada exitosamente")
}

// ChangePassword maneja la solicitud de cambio de contraseña del usuario autenticado
// @Summary Cambiar contraseña
// @Description Cambia la contraseña del usuario autenticado
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.ChangePasswordRequest true "Contraseña actual y nueva"
// @Success 200 {string} string "Contraseña cambiada exitosamente"
// @Failure 400 {string} string "Error en la solicitud"
// @Failure 401 {string} string "No autorizado"
// @Failure 422 {string} string "La contraseña no cumple con las políticas"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /auth/change-password [post]
func (h *PasswordSecurityHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ChangePasswordRequest

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar campos requeridos
	if req.CurrentPassword == "" || req.NewPassword == "" || req.NewPasswordConfirm == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Todos los campos son requeridos")
		return
	}

	// Obtener usuario del contexto (seteado por el middleware de auth)
	u, ok := r.Context().Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || u == nil {
		slog.Error("Usuario no encontrado en el contexto",
			"path", r.URL.Path,
		)
		handler.RespondWithError(w, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	userID := u.ID

	// Procesar cambio de contraseña
	if err := h.passwordSecurityUseCase.ChangePassword(r.Context(), userID, &req); err != nil {
		switch err {
		case auth.ErrPasswordMismatch:
			handler.RespondWithError(w, http.StatusBadRequest, "Las contraseñas no coinciden")
		case auth.ErrInvalidCurrentPassword:
			handler.RespondWithError(w, http.StatusUnauthorized, "Contraseña actual incorrecta")
		case auth.ErrPasswordSameAsCurrent:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La nueva contraseña no puede ser igual a la actual")
		case auth.ErrPasswordTooShort:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña no cumple con la longitud mínima")
		case auth.ErrPasswordNoNumbers:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña debe contener al menos un número")
		case auth.ErrPasswordNoSymbols:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "La contraseña debe contener al menos un símbolo")
		case auth.ErrPasswordReuse:
			handler.RespondWithError(w, http.StatusUnprocessableEntity, "No puedes reutilizar una contraseña anterior")
		default:
			slog.Error("Error al cambiar contraseña",
				"userID", userID,
				"error", err,
			)
			handler.RespondWithError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}

	handler.RespondWithSuccess(w, "Contraseña cambiada exitosamente")
}
