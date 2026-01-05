package permission

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/permission"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con permisos
type Handler struct {
	permissionUseCase *permission.UseCase
	validate          *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de permisos
func NewHandler(permissionUseCase *permission.UseCase) *Handler {
	return &Handler{
		permissionUseCase: permissionUseCase,
		validate:          validator.New(),
	}
}

// CreatePermissionRequest representa la solicitud para crear un permiso
type CreatePermissionRequest struct {
	AbilityID  int64   `json:"ability_id" validate:"required"`
	EntityID   int64   `json:"entity_id" validate:"required"`
	EntityType string  `json:"entity_type" validate:"required"`
	Forbidden  bool    `json:"forbidden"`
	Conditions *string `json:"conditions,omitempty"`
}

// UpdatePermissionRequest representa la solicitud para actualizar un permiso
type UpdatePermissionRequest struct {
	AbilityID  int64   `json:"ability_id" validate:"required"`
	EntityID   int64   `json:"entity_id" validate:"required"`
	EntityType string  `json:"entity_type" validate:"required"`
	Forbidden  bool    `json:"forbidden"`
	Conditions *string `json:"conditions,omitempty"`
}

// PermissionResponse representa la respuesta de un permiso
type PermissionResponse struct {
	ID         int64   `json:"id"`
	AbilityID  int64   `json:"ability_id"`
	EntityID   int64   `json:"entity_id"`
	EntityType string  `json:"entity_type"`
	Forbidden  bool    `json:"forbidden"`
	Conditions *string `json:"conditions,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

// Create maneja la solicitud de creación de un permiso
// @Summary Crear permiso
// @Description Crea un nuevo permiso en el sistema
// @Tags Permissions
// @Accept json
// @Produce json
// @Param permission body permission.CreatePermissionRequest true "Datos del permiso"
// @Success 201 {object} permission.PermissionResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 409 {string} string "El permiso ya existe para esta entidad y ability"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_permission") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para crear permisos")
		return
	}

	// Decodificar la solicitud
	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Crear el permiso
	permission := &permission.Permission{
		AbilityID:  req.AbilityID,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Forbidden:  req.Forbidden,
		Conditions: req.Conditions,
	}

	if err := h.permissionUseCase.Create(r.Context(), permission); err != nil {
		// Verificar si el error es de permiso duplicado
		if err.Error() == "el permiso ya existe para esta entidad y ability" {
			response.Error(w, http.StatusConflict, "El permiso ya existe para esta entidad y ability")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al crear el permiso")
		return
	}

	// Preparar la respuesta
	resp := PermissionResponse{
		ID:         permission.ID,
		AbilityID:  permission.AbilityID,
		EntityID:   permission.EntityID,
		EntityType: permission.EntityType,
		Forbidden:  permission.Forbidden,
		Conditions: permission.Conditions,
	}

	if permission.CreatedAt != nil {
		resp.CreatedAt = permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if permission.UpdatedAt != nil {
		resp.UpdatedAt = permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Get maneja la solicitud de obtención de un permiso por ID
// @Summary Obtener permiso
// @Description Obtiene un permiso por su ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "ID del permiso"
// @Success 200 {object} PermissionResponse
// @Failure 400 {string} string "ID de permiso inválido"
// @Failure 404 {string} string "Permiso no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_permission") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver permisos")
		return
	}

	// Obtener el ID del permiso de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de permiso inválido")
		return
	}

	// Obtener el permiso
	permission, err := h.permissionUseCase.GetByID(r.Context(), id)
	if err != nil {
		// Verificar si el error es de permiso no encontrado
		if err.Error() == "permiso no encontrado" {
			response.Error(w, http.StatusNotFound, "Permiso no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el permiso")
		return
	}

	// Preparar la respuesta
	resp := PermissionResponse{
		ID:         permission.ID,
		AbilityID:  permission.AbilityID,
		EntityID:   permission.EntityID,
		EntityType: permission.EntityType,
		Forbidden:  permission.Forbidden,
		Conditions: permission.Conditions,
	}

	if permission.CreatedAt != nil {
		resp.CreatedAt = permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if permission.UpdatedAt != nil {
		resp.UpdatedAt = permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// GetByEntity maneja la solicitud de obtención de permisos por entidad
// @Summary Obtener permisos por entidad
// @Description Obtiene todos los permisos para una entidad específica
// @Tags Permissions
// @Accept json
// @Produce json
// @Param entity_type path string true "Tipo de entidad"
// @Param entity_id path int true "ID de la entidad"
// @Success 200 {array} PermissionResponse
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/entity/{entity_type}/{entity_id} [get]
// @Security BearerAuth
func (h *Handler) GetByEntity(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_entity_permissions") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver permisos de entidades")
		return
	}

	// Obtener los parámetros de la URL
	entityType := chi.URLParam(r, "entity_type")
	entityIDStr := chi.URLParam(r, "entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de entidad inválido")
		return
	}

	// Obtener los permisos
	permissions, err := h.permissionUseCase.GetByEntity(r.Context(), entityType, entityID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener los permisos")
		return
	}

	// Preparar la respuesta
	var items []PermissionResponse
	for _, p := range permissions {
		item := PermissionResponse{
			ID:         p.ID,
			AbilityID:  p.AbilityID,
			EntityID:   p.EntityID,
			EntityType: p.EntityType,
			Forbidden:  p.Forbidden,
			Conditions: p.Conditions,
		}

		if p.CreatedAt != nil {
			item.CreatedAt = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if p.UpdatedAt != nil {
			item.UpdatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	response.JSON(w, http.StatusOK, items)
}

// GetByAbility maneja la solicitud de obtención de permisos por ability
// @Summary Obtener permisos por ability
// @Description Obtiene todos los permisos para una ability específica
// @Tags Permissions
// @Accept json
// @Produce json
// @Param ability_id path int true "ID de la ability"
// @Success 200 {array} PermissionResponse
// @Failure 400 {string} string "ID de ability inválido"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/ability/{ability_id} [get]
// @Security BearerAuth
func (h *Handler) GetByAbility(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_ability_permissions") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver permisos de abilities")
		return
	}

	// Obtener el ID de la ability de la URL
	abilityIDStr := chi.URLParam(r, "ability_id")
	abilityID, err := strconv.ParseInt(abilityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de ability inválido")
		return
	}

	// Obtener los permisos
	permissions, err := h.permissionUseCase.GetByAbility(r.Context(), abilityID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener los permisos")
		return
	}

	// Preparar la respuesta
	var items []PermissionResponse
	for _, p := range permissions {
		item := PermissionResponse{
			ID:         p.ID,
			AbilityID:  p.AbilityID,
			EntityID:   p.EntityID,
			EntityType: p.EntityType,
			Forbidden:  p.Forbidden,
			Conditions: p.Conditions,
		}

		if p.CreatedAt != nil {
			item.CreatedAt = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if p.UpdatedAt != nil {
			item.UpdatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	response.JSON(w, http.StatusOK, items)
}

// Update maneja la solicitud de actualización de un permiso
// @Summary Actualizar permiso
// @Description Actualiza un permiso existente
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "ID del permiso"
// @Param permission body UpdatePermissionRequest true "Datos del permiso"
// @Success 200 {object} PermissionResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 404 {string} string "Permiso no encontrado"
// @Failure 409 {string} string "El permiso ya existe para esta entidad y ability"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_permission") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar permisos")
		return
	}

	// Obtener el ID del permiso de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de permiso inválido")
		return
	}

	// Decodificar la solicitud
	var req UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Actualizar el permiso
	permission := &permission.Permission{
		ID:         id,
		AbilityID:  req.AbilityID,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Forbidden:  req.Forbidden,
		Conditions: req.Conditions,
	}

	if err := h.permissionUseCase.Update(r.Context(), permission); err != nil {
		// Verificar si el error es de permiso no encontrado
		if err.Error() == "permiso no encontrado" {
			response.Error(w, http.StatusNotFound, "Permiso no encontrado")
			return
		}
		// Verificar si el error es de permiso duplicado
		if err.Error() == "el permiso ya existe para esta entidad y ability" {
			response.Error(w, http.StatusConflict, "El permiso ya existe para esta entidad y ability")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el permiso")
		return
	}

	// Preparar la respuesta
	resp := PermissionResponse{
		ID:         permission.ID,
		AbilityID:  permission.AbilityID,
		EntityID:   permission.EntityID,
		EntityType: permission.EntityType,
		Forbidden:  permission.Forbidden,
		Conditions: permission.Conditions,
	}

	if permission.CreatedAt != nil {
		resp.CreatedAt = permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if permission.UpdatedAt != nil {
		resp.UpdatedAt = permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete maneja la solicitud de eliminación de un permiso
// @Summary Eliminar permiso
// @Description Elimina un permiso
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path int true "ID del permiso"
// @Success 204 "No Content"
// @Failure 400 {string} string "ID de permiso inválido"
// @Failure 404 {string} string "Permiso no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "delete_permission") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para eliminar permisos")
		return
	}

	// Obtener el ID del permiso de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de permiso inválido")
		return
	}

	// Eliminar el permiso
	if err := h.permissionUseCase.Delete(r.Context(), id); err != nil {
		// Verificar si el error es de permiso no encontrado
		if err.Error() == "permiso no encontrado" {
			response.Error(w, http.StatusNotFound, "Permiso no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el permiso")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Exists maneja la solicitud de verificación si existe un permiso específico
// @Summary Verificar permiso
// @Description Verifica si existe un permiso específico
// @Tags Permissions
// @Accept json
// @Produce json
// @Param ability_id path int true "ID de la ability"
// @Param entity_type path string true "Tipo de entidad"
// @Param entity_id path int true "ID de la entidad"
// @Success 200 {object} map[string]bool
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions/check/{ability_id}/{entity_type}/{entity_id} [get]
// @Security BearerAuth
func (h *Handler) Exists(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "check_permission") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para verificar permisos")
		return
	}

	// Obtener los parámetros de la URL
	abilityIDStr := chi.URLParam(r, "ability_id")
	abilityID, err := strconv.ParseInt(abilityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de ability inválido")
		return
	}

	entityType := chi.URLParam(r, "entity_type")
	entityIDStr := chi.URLParam(r, "entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de entidad inválido")
		return
	}

	// Verificar si existe el permiso
	exists, err := h.permissionUseCase.Exists(r.Context(), abilityID, entityID, entityType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al verificar el permiso")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

// List maneja la solicitud de listado de permisos
// @Summary Listar permisos
// @Description Obtiene una lista paginada de permisos
// @Tags Permissions
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param page_size query int false "Tamaño de página (por defecto: 10)"
// @Param ability_id query int false "Filtrar por ID de ability"
// @Param entity_type query string false "Filtrar por tipo de entidad"
// @Param entity_id query int false "Filtrar por ID de entidad"
// @Param forbidden query bool false "Filtrar por prohibición"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {string} string "Parámetros de consulta inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /permissions [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "list_permissions") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar permisos")
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

	if abilityIDStr := r.URL.Query().Get("ability_id"); abilityIDStr != "" {
		abilityID, err := strconv.ParseInt(abilityIDStr, 10, 64)
		if err == nil {
			filters["ability_id"] = abilityID
		}
	}

	if entityType := r.URL.Query().Get("entity_type"); entityType != "" {
		filters["entity_type"] = entityType
	}

	if entityIDStr := r.URL.Query().Get("entity_id"); entityIDStr != "" {
		entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
		if err == nil {
			filters["entity_id"] = entityID
		}
	}

	if forbiddenStr := r.URL.Query().Get("forbidden"); forbiddenStr != "" {
		switch forbiddenStr {
		case "true":
			filters["forbidden"] = true
		case "false":
			filters["forbidden"] = false
		}
	}

	// Obtener la lista de permisos
	permissions, total, err := h.permissionUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar permisos")
		return
	}

	// Preparar la respuesta
	var items []PermissionResponse
	for _, p := range permissions {
		item := PermissionResponse{
			ID:         p.ID,
			AbilityID:  p.AbilityID,
			EntityID:   p.EntityID,
			EntityType: p.EntityType,
			Forbidden:  p.Forbidden,
			Conditions: p.Conditions,
		}

		if p.CreatedAt != nil {
			item.CreatedAt = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if p.UpdatedAt != nil {
			item.UpdatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	// Enviar respuesta paginada
	response.Paginated(w, items, page, pageSize, total)
}
