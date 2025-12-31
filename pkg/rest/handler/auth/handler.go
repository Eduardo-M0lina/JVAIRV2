package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/auth"
)

// Handler maneja las solicitudes HTTP relacionadas con autenticación
type Handler struct {
	authUseCase *auth.UseCase
}

// NewHandler crea una nueva instancia del handler de autenticación
func NewHandler(authUseCase *auth.UseCase) *Handler {
	return &Handler{
		authUseCase: authUseCase,
	}
}

// Login maneja la solicitud de inicio de sesión
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Validar la solicitud
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email y contraseña son requeridos", http.StatusBadRequest)
		return
	}

	// Autenticar usuario
	resp, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
			return
		}
		if err == auth.ErrUserInactive {
			http.Error(w, "Usuario inactivo", http.StatusForbidden)
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Responder con los tokens
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}

// Logout maneja la solicitud de cierre de sesión
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Obtener token de acceso del encabezado de autorización
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token de acceso no proporcionado", http.StatusBadRequest)
		return
	}

	// Extraer token del encabezado (Bearer token)
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Formato de token inválido", http.StatusBadRequest)
		return
	}
	accessToken := splitToken[1]

	// Cerrar sesión
	err := h.authUseCase.Logout(r.Context(), accessToken)
	if err != nil {
		http.Error(w, "Error al cerrar sesión", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusNoContent)
}

// RefreshToken maneja la solicitud de refresco de token
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Validar la solicitud
	if req.RefreshToken == "" {
		http.Error(w, "Token de refresco es requerido", http.StatusBadRequest)
		return
	}

	// Refrescar token
	td, err := h.authUseCase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if err == auth.ErrInvalidToken {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Responder con los nuevos tokens
	resp := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
		ExpiresAt:    td.AtExpires,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}
