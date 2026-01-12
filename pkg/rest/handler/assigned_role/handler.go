package assigned_role

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/assigned_role"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con asignaciones de roles
type Handler struct {
	assignedRoleUseCase *assigned_role.UseCase
	validate            *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de asignaciones de roles
func NewHandler(assignedRoleUseCase *assigned_role.UseCase) *Handler {
	return &Handler{
		assignedRoleUseCase: assignedRoleUseCase,
		validate:            validator.New(),
	}
}

// AssignRoleRequest representa la solicitud para asignar un rol a una entidad
type AssignRoleRequest struct {
	RoleID     int64  `json:"role_id" validate:"required"`
	EntityID   int64  `json:"entity_id" validate:"required"`
	EntityType string `json:"entity_type" validate:"required"`
	Restricted bool   `json:"restricted"`
	Scope      *int   `json:"scope,omitempty"`
}

// AssignedRoleResponse representa la respuesta de una asignación de rol
type AssignedRoleResponse struct {
	ID         int64  `json:"id"`
	RoleID     int64  `json:"role_id"`
	EntityID   int64  `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Restricted bool   `json:"restricted"`
	Scope      *int   `json:"scope,omitempty"`
}

// Assign maneja la solicitud de asignación de un rol a una entidad
// @Summary Asignar rol
// @Description Asigna un rol a una entidad
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param role body assigned_role.AssignRoleRequest true "Datos de la asignación de rol"
// @Success 201 {object} assigned_role.AssignedRoleResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 409 {string} string "El rol ya está asignado a esta entidad"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles [post]
// @Security BearerAuth
func (h *Handler) Assign(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "assign_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para asignar roles")
		return
	}

	// Decodificar la solicitud
	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Asignar el rol
	assignedRole := &assigned_role.AssignedRole{
		RoleID:     req.RoleID,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Restricted: req.Restricted,
		Scope:      req.Scope,
	}

	if err := h.assignedRoleUseCase.Assign(r.Context(), assignedRole); err != nil {
		// Verificar si el error es de asignación duplicada
		if err.Error() == "el rol ya está asignado a esta entidad" {
			response.Error(w, http.StatusConflict, "El rol ya está asignado a esta entidad")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al asignar el rol")
		return
	}

	// Preparar la respuesta
	resp := AssignedRoleResponse{
		ID:         assignedRole.ID,
		RoleID:     assignedRole.RoleID,
		EntityID:   assignedRole.EntityID,
		EntityType: assignedRole.EntityType,
		Restricted: assignedRole.Restricted,
		Scope:      assignedRole.Scope,
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Get maneja la solicitud de obtención de una asignación de rol por ID
// @Summary Obtener asignación de rol
// @Description Obtiene una asignación de rol por su ID
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param id path int true "ID de la asignación de rol"
// @Success 200 {object} AssignedRoleResponse
// @Failure 400 {string} string "ID de asignación de rol inválido"
// @Failure 404 {string} string "Asignación de rol no encontrada"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_assigned_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver asignaciones de roles")
		return
	}

	// Obtener el ID de la asignación de rol de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de asignación de rol inválido")
		return
	}

	// Obtener la asignación de rol
	assignedRole, err := h.assignedRoleUseCase.GetByID(r.Context(), id)
	if err != nil {
		// Verificar si el error es de asignación de rol no encontrada
		if err.Error() == "asignación de rol no encontrada" {
			response.Error(w, http.StatusNotFound, "Asignación de rol no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener la asignación de rol")
		return
	}

	// Preparar la respuesta
	resp := AssignedRoleResponse{
		ID:         assignedRole.ID,
		RoleID:     assignedRole.RoleID,
		EntityID:   assignedRole.EntityID,
		EntityType: assignedRole.EntityType,
		Restricted: assignedRole.Restricted,
		Scope:      assignedRole.Scope,
	}

	response.JSON(w, http.StatusOK, resp)
}

// GetByEntity maneja la solicitud de obtención de asignaciones de roles por entidad
// @Summary Obtener asignaciones de roles por entidad
// @Description Obtiene todas las asignaciones de roles para una entidad específica
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param entity_type path string true "Tipo de entidad"
// @Param entity_id path int true "ID de la entidad"
// @Success 200 {array} AssignedRoleResponse
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles/entity/{entity_type}/{entity_id} [get]
// @Security BearerAuth
func (h *Handler) GetByEntity(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_entity_roles") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver roles de entidades")
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

	// Obtener las asignaciones de roles
	assignedRoles, err := h.assignedRoleUseCase.GetByEntity(r.Context(), entityType, entityID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener las asignaciones de roles")
		return
	}

	// Preparar la respuesta
	var items []AssignedRoleResponse
	for _, ar := range assignedRoles {
		item := AssignedRoleResponse{
			ID:         ar.ID,
			RoleID:     ar.RoleID,
			EntityID:   ar.EntityID,
			EntityType: ar.EntityType,
			Restricted: ar.Restricted,
			Scope:      ar.Scope,
		}

		items = append(items, item)
	}

	response.JSON(w, http.StatusOK, items)
}

// Revoke maneja la solicitud de revocación de un rol de una entidad
// @Summary Revocar rol
// @Description Revoca un rol de una entidad
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param role_id path int true "ID del rol"
// @Param entity_type path string true "Tipo de entidad"
// @Param entity_id path int true "ID de la entidad"
// @Success 204 "No Content"
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 404 {string} string "Asignación de rol no encontrada"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles/revoke/{role_id}/{entity_type}/{entity_id} [delete]
// @Security BearerAuth
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "revoke_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para revocar roles")
		return
	}

	// Obtener los parámetros de la URL
	roleIDStr := chi.URLParam(r, "role_id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de rol inválido")
		return
	}

	entityType := chi.URLParam(r, "entity_type")
	entityIDStr := chi.URLParam(r, "entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de entidad inválido")
		return
	}

	// Revocar el rol
	if err := h.assignedRoleUseCase.Revoke(r.Context(), roleID, entityID, entityType); err != nil {
		// Verificar si el error es de asignación de rol no encontrada
		if err.Error() == "asignación de rol no encontrada" {
			response.Error(w, http.StatusNotFound, "Asignación de rol no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al revocar el rol")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HasRole maneja la solicitud de verificación si una entidad tiene un rol específico
// @Summary Verificar rol
// @Description Verifica si una entidad tiene un rol específico
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param role_id path int true "ID del rol"
// @Param entity_type path string true "Tipo de entidad"
// @Param entity_id path int true "ID de la entidad"
// @Success 200 {object} map[string]bool
// @Failure 400 {string} string "Parámetros inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles/check/{role_id}/{entity_type}/{entity_id} [get]
// @Security BearerAuth
func (h *Handler) HasRole(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "check_role") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para verificar roles")
		return
	}

	// Obtener los parámetros de la URL
	roleIDStr := chi.URLParam(r, "role_id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de rol inválido")
		return
	}

	entityType := chi.URLParam(r, "entity_type")
	entityIDStr := chi.URLParam(r, "entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de entidad inválido")
		return
	}

	// Verificar si la entidad tiene el rol
	hasRole, err := h.assignedRoleUseCase.HasRole(r.Context(), roleID, entityID, entityType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al verificar el rol")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"has_role": hasRole})
}

// List maneja la solicitud de listado de asignaciones de roles
// @Summary Listar asignaciones de roles
// @Description Obtiene una lista paginada de asignaciones de roles
// @Tags AssignedRoles
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param page_size query int false "Tamaño de página (por defecto: 10)"
// @Param role_id query int false "Filtrar por ID de rol"
// @Param entity_type query string false "Filtrar por tipo de entidad"
// @Param entity_id query int false "Filtrar por ID de entidad"
// @Param restricted query bool false "Filtrar por restricción"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {string} string "Parámetros de consulta inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/assigned-roles [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "list_assigned_roles") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar asignaciones de roles")
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

	if roleIDStr := r.URL.Query().Get("role_id"); roleIDStr != "" {
		roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
		if err == nil {
			filters["role_id"] = roleID
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

	if restrictedStr := r.URL.Query().Get("restricted"); restrictedStr != "" {
		switch restrictedStr {
		case "true":
			filters["restricted"] = true
		case "false":
			filters["restricted"] = false
		}
	}

	// Obtener la lista de asignaciones de roles
	assignedRoles, total, err := h.assignedRoleUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		log.Printf("ERROR al listar asignaciones de roles: %v", err)
		response.Error(w, http.StatusInternalServerError, "Error al listar asignaciones de roles: "+err.Error())
		return
	}

	// Preparar la respuesta
	var items []AssignedRoleResponse
	for _, ar := range assignedRoles {
		item := AssignedRoleResponse{
			ID:         ar.ID,
			RoleID:     ar.RoleID,
			EntityID:   ar.EntityID,
			EntityType: ar.EntityType,
			Restricted: ar.Restricted,
			Scope:      ar.Scope,
		}

		items = append(items, item)
	}

	// Enviar respuesta paginada
	response.Paginated(w, items, page, pageSize, total)
}
