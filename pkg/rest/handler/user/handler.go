package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/user"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con usuarios
type Handler struct {
	userUseCase *user.UseCase
	validate    *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de usuarios
func NewHandler(userUseCase *user.UseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
		validate:    validator.New(),
	}
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	RoleID   *string `json:"roleId,omitempty"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password,omitempty" validate:"omitempty,min=6"`
	RoleID   *string `json:"roleId,omitempty"`
	IsActive bool    `json:"isActive"`
}

// UserResponse representa la respuesta de un usuario
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	RoleID    *string   `json:"roleId,omitempty"`
	Role      *RoleInfo `json:"role,omitempty"`
	IsActive  bool      `json:"isActive"`
	CreatedAt string    `json:"createdAt,omitempty"`
	UpdatedAt string    `json:"updatedAt,omitempty"`
}

// RoleInfo representa la información básica del rol en la respuesta
type RoleInfo struct {
	Name  string  `json:"name"`
	Title *string `json:"title,omitempty"`
}

// Create maneja la solicitud de creación de un usuario
// @Summary Crear usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags Users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} user.UserResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 409 {string} string "Email ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_user") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para crear usuarios")
		return
	}

	// Decodificar la solicitud
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Crear el usuario
	u := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	}

	if err := h.userUseCase.Create(r.Context(), u); err != nil {
		if err == user.ErrDuplicateEmail {
			response.Error(w, http.StatusConflict, "Email ya está en uso")
			return
		}
		// Registrar el error detallado para depuración
		slog.Error("Error al crear usuario",
			"name", req.Name,
			"email", req.Email,
			"error", err,
		)
		response.Error(w, http.StatusInternalServerError, "Error al crear el usuario: "+err.Error())
		return
	}

	// Preparar la respuesta
	resp := UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}

	// Agregar información del rol si existe
	if u.RoleName != nil {
		resp.Role = &RoleInfo{
			Name:  *u.RoleName,
			Title: u.RoleTitle,
		}
	}

	if u.CreatedAt != nil {
		resp.CreatedAt = u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if u.UpdatedAt != nil {
		resp.UpdatedAt = u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Get maneja la solicitud de obtención de un usuario por ID
// @Summary Obtener usuario
// @Description Obtiene un usuario por su ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} UserResponse
// @Failure 400 {string} string "ID de usuario inválido"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_user") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver usuarios")
		return
	}

	// Obtener el ID del usuario de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Obtener el usuario
	u, err := h.userUseCase.GetByID(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el usuario")
		return
	}

	// Preparar la respuesta
	resp := UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}

	// Agregar información del rol si existe
	if u.RoleName != nil {
		resp.Role = &RoleInfo{
			Name:  *u.RoleName,
			Title: u.RoleTitle,
		}
	}

	if u.CreatedAt != nil {
		resp.CreatedAt = u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if u.UpdatedAt != nil {
		resp.UpdatedAt = u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Update maneja la solicitud de actualización de un usuario
// @Summary Actualizar usuario
// @Description Actualiza un usuario existente
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param user body UpdateUserRequest true "Datos del usuario"
// @Success 200 {object} UserResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 409 {string} string "Email ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_user") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar usuarios")
		return
	}

	// Obtener el ID del usuario de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Decodificar la solicitud
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Actualizar el usuario
	u := &user.User{
		ID:       idInt,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
		IsActive: req.IsActive,
	}

	if err := h.userUseCase.Update(r.Context(), u); err != nil {
		if err == user.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		if err == user.ErrDuplicateEmail {
			response.Error(w, http.StatusConflict, "Email ya está en uso")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el usuario")
		return
	}

	// Preparar la respuesta
	resp := UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}

	// Agregar información del rol si existe
	if u.RoleName != nil {
		resp.Role = &RoleInfo{
			Name:  *u.RoleName,
			Title: u.RoleTitle,
		}
	}

	if u.CreatedAt != nil {
		resp.CreatedAt = u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if u.UpdatedAt != nil {
		resp.UpdatedAt = u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete maneja la solicitud de eliminación de un usuario
// @Summary Eliminar usuario
// @Description Elimina un usuario (soft delete)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 204 "No Content"
// @Failure 400 {string} string "ID de usuario inválido"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "delete_user") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para eliminar usuarios")
		return
	}

	// Obtener el ID del usuario de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Eliminar el usuario
	if err := h.userUseCase.Delete(r.Context(), id); err != nil {
		if err == user.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el usuario")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List maneja la solicitud de listado de usuarios
// @Summary Listar usuarios
// @Description Obtiene una lista paginada de usuarios
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param page_size query int false "Tamaño de página (por defecto: 10)"
// @Param name query string false "Filtrar por nombre"
// @Param email query string false "Filtrar por email"
// @Param role_id query string false "Filtrar por rol"
// @Param is_active query bool false "Filtrar por estado (activo/inactivo)"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {string} string "Parámetros de consulta inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "list_users") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar usuarios")
		return
	}

	// Obtener parámetros de consulta
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Construir filtros
	filters := make(map[string]interface{})

	if name := r.URL.Query().Get("name"); name != "" {
		filters["name"] = name
	}

	if email := r.URL.Query().Get("email"); email != "" {
		filters["email"] = email
	}

	if roleID := r.URL.Query().Get("role_id"); roleID != "" {
		filters["role_id"] = roleID
	}

	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		switch isActive {
		case "true":
			filters["is_active"] = true
		case "false":
			filters["is_active"] = false
		}
	}

	// Obtener la lista de usuarios
	users, total, err := h.userUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar usuarios")
		return
	}

	slog.Info("Usuarios obtenidos del repositorio",
		"count", len(users),
		"total", total,
	)

	// Preparar la respuesta
	var items []UserResponse
	for _, u := range users {
		slog.Debug("Procesando usuario para respuesta",
			"user_id", u.ID,
			"email", u.Email,
			"role_name", u.RoleName,
			"role_title", u.RoleTitle,
		)
		item := UserResponse{
			ID:       u.ID,
			Name:     u.Name,
			Email:    u.Email,
			RoleID:   u.RoleID,
			IsActive: u.IsActive,
		}

		// Agregar información del rol si existe
		if u.RoleName != nil {
			item.Role = &RoleInfo{
				Name:  *u.RoleName,
				Title: u.RoleTitle,
			}
			slog.Debug("Rol agregado a la respuesta",
				"user_id", u.ID,
				"role_name", *u.RoleName,
			)
		} else {
			slog.Debug("Usuario sin rol en la respuesta",
				"user_id", u.ID,
				"email", u.Email,
			)
		}

		if u.CreatedAt != nil {
			item.CreatedAt = u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if u.UpdatedAt != nil {
			item.UpdatedAt = u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	// Enviar respuesta paginada
	response.Paginated(w, items, page, pageSize, total)
}

// GetRoles maneja la solicitud de obtención de roles de un usuario
// @Summary Obtener roles de usuario
// @Description Obtiene los roles asignados a un usuario
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {array} role.Role
// @Failure 400 {string} string "ID de usuario inválido"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users/{id}/roles [get]
// @Security BearerAuth
func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_user_roles") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver roles de usuarios")
		return
	}

	// Obtener el ID del usuario de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Obtener los roles del usuario
	roles, err := h.userUseCase.GetUserRoles(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener roles del usuario")
		return
	}

	response.JSON(w, http.StatusOK, roles)
}

// GetAbilities maneja la solicitud de obtención de habilidades de un usuario
// @Summary Obtener habilidades de usuario
// @Description Obtiene las habilidades asignadas a un usuario
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {array} ability.Ability
// @Failure 400 {string} string "ID de usuario inválido"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/users/{id}/abilities [get]
// @Security BearerAuth
func (h *Handler) GetAbilities(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_user_abilities") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver habilidades de usuarios")
		return
	}

	// Obtener el ID del usuario de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Obtener las habilidades del usuario
	abilities, err := h.userUseCase.GetUserAbilities(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener habilidades del usuario")
		return
	}

	response.JSON(w, http.StatusOK, abilities)
}
