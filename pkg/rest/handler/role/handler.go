package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/role"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con roles
type Handler struct {
	roleUseCase *role.UseCase
	validate    *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de roles
func NewHandler(roleUseCase *role.UseCase) *Handler {
	return &Handler{
		roleUseCase: roleUseCase,
		validate:    validator.New(),
	}
}

// CreateRoleRequest representa la solicitud para crear un rol
type CreateRoleRequest struct {
	Name  string `json:"name" validate:"required"`
	Title string `json:"title,omitempty"`
	Scope *int   `json:"scope,omitempty"`
}

// UpdateRoleRequest representa la solicitud para actualizar un rol
type UpdateRoleRequest struct {
	Name  string `json:"name" validate:"required"`
	Title string `json:"title,omitempty"`
	Scope *int   `json:"scope,omitempty"`
}

// RoleResponse representa la respuesta de un rol
type RoleResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Title     *string `json:"title,omitempty"`
	Scope     *int    `json:"scope,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

// Create maneja la solicitud de creación de un rol
// @Summary Crear rol
// @Description Crea un nuevo rol en el sistema
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body role.CreateRoleRequest true "Datos del rol"
// @Success 201 {object} role.RoleResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 409 {string} string "Nombre de rol ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /roles [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para crear roles")
		return
	}

	// Decodificar la solicitud
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Crear el rol
	role := &role.Role{
		Name:  req.Name,
		Scope: req.Scope,
	}

	if req.Title != "" {
		title := req.Title
		role.Title = &title
	}

	if err := h.roleUseCase.Create(r.Context(), role); err != nil {
		// Verificar si el error es de nombre duplicado
		if err.Error() == "nombre de rol ya está en uso" {
			response.Error(w, http.StatusConflict, "Nombre de rol ya está en uso")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al crear el rol")
		return
	}

	// Preparar la respuesta
	resp := RoleResponse{
		ID:    role.ID,
		Name:  role.Name,
		Title: role.Title,
		Scope: role.Scope,
	}

	if role.CreatedAt != nil {
		resp.CreatedAt = role.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if role.UpdatedAt != nil {
		resp.UpdatedAt = role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Get maneja la solicitud de obtención de un rol por ID
// @Summary Obtener rol
// @Description Obtiene un rol por su ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Success 200 {object} RoleResponse
// @Failure 400 {string} string "ID de rol inválido"
// @Failure 404 {string} string "Rol no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /roles/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver roles")
		return
	}

	// Obtener el ID del rol de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de rol inválido")
		return
	}

	// Obtener el rol
	role, err := h.roleUseCase.GetByID(r.Context(), id)
	if err != nil {
		// Verificar si el error es de rol no encontrado
		if err.Error() == "rol no encontrado" {
			response.Error(w, http.StatusNotFound, "Rol no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el rol")
		return
	}

	// Preparar la respuesta
	resp := RoleResponse{
		ID:    role.ID,
		Name:  role.Name,
		Title: role.Title,
		Scope: role.Scope,
	}

	if role.CreatedAt != nil {
		resp.CreatedAt = role.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if role.UpdatedAt != nil {
		resp.UpdatedAt = role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Update maneja la solicitud de actualización de un rol
// @Summary Actualizar rol
// @Description Actualiza un rol existente
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param role body UpdateRoleRequest true "Datos del rol"
// @Success 200 {object} RoleResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 404 {string} string "Rol no encontrado"
// @Failure 409 {string} string "Nombre de rol ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /roles/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar roles")
		return
	}

	// Obtener el ID del rol de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de rol inválido")
		return
	}

	// Decodificar la solicitud
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Actualizar el rol
	role := &role.Role{
		ID:    id,
		Name:  req.Name,
		Scope: req.Scope,
	}

	if req.Title != "" {
		title := req.Title
		role.Title = &title
	}

	if err := h.roleUseCase.Update(r.Context(), role); err != nil {
		// Verificar si el error es de rol no encontrado
		if err.Error() == "rol no encontrado" {
			response.Error(w, http.StatusNotFound, "Rol no encontrado")
			return
		}
		// Verificar si el error es de nombre duplicado
		if err.Error() == "nombre de rol ya está en uso" {
			response.Error(w, http.StatusConflict, "Nombre de rol ya está en uso")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el rol")
		return
	}

	// Preparar la respuesta
	resp := RoleResponse{
		ID:    role.ID,
		Name:  role.Name,
		Title: role.Title,
		Scope: role.Scope,
	}

	if role.CreatedAt != nil {
		resp.CreatedAt = role.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if role.UpdatedAt != nil {
		resp.UpdatedAt = role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete maneja la solicitud de eliminación de un rol
// @Summary Eliminar rol
// @Description Elimina un rol
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Success 204 "No Content"
// @Failure 400 {string} string "ID de rol inválido"
// @Failure 404 {string} string "Rol no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /roles/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "delete_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para eliminar roles")
		return
	}

	// Obtener el ID del rol de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de rol inválido")
		return
	}

	// Eliminar el rol
	if err := h.roleUseCase.Delete(r.Context(), id); err != nil {
		// Verificar si el error es de rol no encontrado
		if err.Error() == "rol no encontrado" {
			response.Error(w, http.StatusNotFound, "Rol no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el rol")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List maneja la solicitud de listado de roles
// @Summary Listar roles
// @Description Obtiene una lista paginada de roles
// @Tags Roles
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param page_size query int false "Tamaño de página (por defecto: 10)"
// @Param name query string false "Filtrar por nombre"
// @Param scope query int false "Filtrar por scope"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {string} string "Parámetros de consulta inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /roles [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "list_roles") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar roles")
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

	if scopeStr := r.URL.Query().Get("scope"); scopeStr != "" {
		scope, err := strconv.Atoi(scopeStr)
		if err == nil {
			filters["scope"] = scope
		}
	}

	// Obtener la lista de roles
	roles, total, err := h.roleUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar roles")
		return
	}

	// Preparar la respuesta
	var items []RoleResponse
	for _, role := range roles {
		item := RoleResponse{
			ID:    role.ID,
			Name:  role.Name,
			Title: role.Title,
			Scope: role.Scope,
		}

		if role.CreatedAt != nil {
			item.CreatedAt = role.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if role.UpdatedAt != nil {
			item.UpdatedAt = role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	// Enviar respuesta paginada
	response.Paginated(w, items, page, pageSize, total)
}
