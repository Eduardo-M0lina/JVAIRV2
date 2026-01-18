package settings

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/settings"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con configuraciones
type Handler struct {
	settingsUseCase *settings.UseCase
	validate        *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de configuraciones
func NewHandler(settingsUseCase *settings.UseCase) *Handler {
	return &Handler{
		settingsUseCase: settingsUseCase,
		validate:        validator.New(),
	}
}

// UpdateSettingsRequest representa la solicitud para actualizar configuraciones
type UpdateSettingsRequest struct {
	IsTwilioEnabled               bool    `json:"isTwilioEnabled"`
	TwilioSID                     *string `json:"twilioSid,omitempty"`
	TwilioAuthToken               *string `json:"twilioAuthToken,omitempty"`
	TwilioFromNumber              *string `json:"twilioFromNumber,omitempty"`
	IsEnforceRoutinePasswordReset bool    `json:"isEnforceRoutinePasswordReset"`
	PasswordExpireDays            int     `json:"passwordExpireDays" validate:"required,min=1"`
	PasswordHistoryCount          int     `json:"passwordHistoryCount" validate:"required,min=0"`
	PasswordMinimumLength         int     `json:"passwordMinimumLength" validate:"required,min=4"`
	PasswordAge                   int     `json:"passwordAge" validate:"required,min=0"`
	PasswordIncludeNumbers        bool    `json:"passwordIncludeNumbers"`
	PasswordIncludeSymbols        bool    `json:"passwordIncludeSymbols"`
}

// SettingsResponse representa la respuesta de configuraciones
type SettingsResponse struct {
	ID                            int64   `json:"id"`
	IsTwilioEnabled               bool    `json:"isTwilioEnabled"`
	TwilioSID                     *string `json:"twilioSid,omitempty"`
	TwilioAuthToken               *string `json:"twilioAuthToken,omitempty"`
	TwilioFromNumber              *string `json:"twilioFromNumber,omitempty"`
	IsEnforceRoutinePasswordReset bool    `json:"isEnforceRoutinePasswordReset"`
	PasswordExpireDays            int     `json:"passwordExpireDays"`
	PasswordHistoryCount          int     `json:"passwordHistoryCount"`
	PasswordMinimumLength         int     `json:"passwordMinimumLength"`
	PasswordAge                   int     `json:"passwordAge"`
	PasswordIncludeNumbers        bool    `json:"passwordIncludeNumbers"`
	PasswordIncludeSymbols        bool    `json:"passwordIncludeSymbols"`
}

// Get maneja la solicitud de obtención de configuraciones
// @Summary Obtener configuraciones del sistema
// @Description Obtiene las configuraciones generales del sistema (políticas de contraseñas, Twilio, etc.)
// @Tags Settings
// @Accept json
// @Produce json
// @Success 200 {object} SettingsResponse "Configuraciones obtenidas exitosamente"
// @Failure 403 {string} string "No tiene permisos para ver configuraciones"
// @Failure 404 {string} string "Configuraciones no encontradas"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/settings [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_settings") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver configuraciones")
		return
	}

	// Obtener las configuraciones
	settingsData, err := h.settingsUseCase.Get(r.Context())
	if err != nil {
		if err == settings.ErrSettingsNotFound {
			response.Error(w, http.StatusNotFound, "Configuraciones no encontradas")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener las configuraciones")
		return
	}

	// Preparar la respuesta
	resp := SettingsResponse{
		ID:                            settingsData.ID,
		IsTwilioEnabled:               settingsData.IsTwilioEnabled,
		TwilioSID:                     settingsData.TwilioSID,
		TwilioAuthToken:               settingsData.TwilioAuthToken,
		TwilioFromNumber:              settingsData.TwilioFromNumber,
		IsEnforceRoutinePasswordReset: settingsData.IsEnforceRoutinePasswordReset,
		PasswordExpireDays:            settingsData.PasswordExpireDays,
		PasswordHistoryCount:          settingsData.PasswordHistoryCount,
		PasswordMinimumLength:         settingsData.PasswordMinimumLength,
		PasswordAge:                   settingsData.PasswordAge,
		PasswordIncludeNumbers:        settingsData.PasswordIncludeNumbers,
		PasswordIncludeSymbols:        settingsData.PasswordIncludeSymbols,
	}

	response.JSON(w, http.StatusOK, resp)
}

// Update maneja la solicitud de actualización de configuraciones
// @Summary Actualizar configuraciones del sistema
// @Description Actualiza las configuraciones generales del sistema. Ejemplo: {"is_twilio_enabled": false, "twilio_sid": null, "twilio_auth_token": null, "twilio_from_number": null, "is_enforce_routine_password_reset": true, "password_expire_days": 90, "password_history_count": 10, "password_minimum_length": 8, "password_age": 5, "password_include_numbers": true, "password_include_symbols": true}
// @Tags Settings
// @Accept json
// @Produce json
// @Param settings body UpdateSettingsRequest true "Datos de las configuraciones a actualizar"
// @Success 200 {object} SettingsResponse "Configuraciones actualizadas exitosamente"
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 403 {string} string "No tiene permisos para actualizar configuraciones"
// @Failure 404 {string} string "Configuraciones no encontradas"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/settings [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_settings") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar configuraciones")
		return
	}

	// Obtener las configuraciones actuales primero
	currentSettings, err := h.settingsUseCase.Get(r.Context())
	if err != nil {
		if err == settings.ErrSettingsNotFound {
			response.Error(w, http.StatusNotFound, "Configuraciones no encontradas")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener las configuraciones")
		return
	}

	// Decodificar la solicitud
	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Actualizar los campos
	currentSettings.IsTwilioEnabled = req.IsTwilioEnabled
	currentSettings.TwilioSID = req.TwilioSID
	currentSettings.TwilioAuthToken = req.TwilioAuthToken
	currentSettings.TwilioFromNumber = req.TwilioFromNumber
	currentSettings.IsEnforceRoutinePasswordReset = req.IsEnforceRoutinePasswordReset
	currentSettings.PasswordExpireDays = req.PasswordExpireDays
	currentSettings.PasswordHistoryCount = req.PasswordHistoryCount
	currentSettings.PasswordMinimumLength = req.PasswordMinimumLength
	currentSettings.PasswordAge = req.PasswordAge
	currentSettings.PasswordIncludeNumbers = req.PasswordIncludeNumbers
	currentSettings.PasswordIncludeSymbols = req.PasswordIncludeSymbols

	// Actualizar las configuraciones
	if err := h.settingsUseCase.Update(r.Context(), currentSettings); err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al actualizar las configuraciones: "+err.Error())
		return
	}

	// Preparar la respuesta
	resp := SettingsResponse{
		ID:                            currentSettings.ID,
		IsTwilioEnabled:               currentSettings.IsTwilioEnabled,
		TwilioSID:                     currentSettings.TwilioSID,
		TwilioAuthToken:               currentSettings.TwilioAuthToken,
		TwilioFromNumber:              currentSettings.TwilioFromNumber,
		IsEnforceRoutinePasswordReset: currentSettings.IsEnforceRoutinePasswordReset,
		PasswordExpireDays:            currentSettings.PasswordExpireDays,
		PasswordHistoryCount:          currentSettings.PasswordHistoryCount,
		PasswordMinimumLength:         currentSettings.PasswordMinimumLength,
		PasswordAge:                   currentSettings.PasswordAge,
		PasswordIncludeNumbers:        currentSettings.PasswordIncludeNumbers,
		PasswordIncludeSymbols:        currentSettings.PasswordIncludeSymbols,
	}

	response.JSON(w, http.StatusOK, resp)
}
