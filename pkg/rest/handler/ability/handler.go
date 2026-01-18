package ability

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/ability"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con abilities
type Handler struct {
	abilityUseCase *ability.UseCase
	validate       *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de abilities
func NewHandler(abilityUseCase *ability.UseCase) *Handler {
	return &Handler{
		abilityUseCase: abilityUseCase,
		validate:       validator.New(),
	}
}

// CreateAbilityRequest representa la solicitud para crear una ability
type CreateAbilityRequest struct {
	Name       string  `json:"name" validate:"required"`
	Title      *string `json:"title,omitempty"`
	EntityID   *int64  `json:"entityId,omitempty"`
	EntityType *string `json:"entityType,omitempty"`
	OnlyOwned  bool    `json:"onlyOwned"`
	Options    *string `json:"options,omitempty"`
	Scope      *int    `json:"scope,omitempty"`
}

// UpdateAbilityRequest representa la solicitud para actualizar una ability
type UpdateAbilityRequest struct {
	Name       string  `json:"name" validate:"required"`
	Title      *string `json:"title,omitempty"`
	EntityID   *int64  `json:"entityId,omitempty"`
	EntityType *string `json:"entityType,omitempty"`
	OnlyOwned  bool    `json:"onlyOwned"`
	Options    *string `json:"options,omitempty"`
	Scope      *int    `json:"scope,omitempty"`
}

// AbilityResponse representa la respuesta de una ability
type AbilityResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Title      *string `json:"title,omitempty"`
	EntityID   *int64  `json:"entityId,omitempty"`
	EntityType *string `json:"entityType,omitempty"`
	OnlyOwned  bool    `json:"onlyOwned"`
	Options    *string `json:"options,omitempty"`
	Scope      *int    `json:"scope,omitempty"`
	CreatedAt  string  `json:"createdAt,omitempty"`
	UpdatedAt  string  `json:"updatedAt,omitempty"`
}

// Create maneja la solicitud de creación de una ability
// @Summary Crear ability
// @Description Crea una nueva ability en el sistema
// @Tags Abilities
// @Accept json
// @Produce json
// @Param ability body ability.CreateAbilityRequest true "Datos de la ability"
// @Success 201 {object} ability.AbilityResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 409 {string} string "Nombre de ability ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/abilities [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_ability") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para crear abilities")
		return
	}

	// Decodificar la solicitud
	var req CreateAbilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Crear la ability
	ability := &ability.Ability{
		Name:       req.Name,
		Title:      req.Title,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		OnlyOwned:  req.OnlyOwned,
		Options:    req.Options,
		Scope:      req.Scope,
	}

	if err := h.abilityUseCase.Create(r.Context(), ability); err != nil {
		// Verificar si el error es de nombre duplicado
		if err.Error() == "nombre de ability ya está en uso" {
			response.Error(w, http.StatusConflict, "Nombre de ability ya está en uso")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al crear la ability")
		return
	}

	// Preparar la respuesta
	resp := AbilityResponse{
		ID:         ability.ID,
		Name:       ability.Name,
		Title:      ability.Title,
		EntityID:   ability.EntityID,
		EntityType: ability.EntityType,
		OnlyOwned:  ability.OnlyOwned,
		Options:    ability.Options,
		Scope:      ability.Scope,
	}

	if ability.CreatedAt != nil {
		resp.CreatedAt = ability.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if ability.UpdatedAt != nil {
		resp.UpdatedAt = ability.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Get maneja la solicitud de obtención de una ability por ID
// @Summary Obtener ability
// @Description Obtiene una ability por su ID
// @Tags Abilities
// @Accept json
// @Produce json
// @Param id path int true "ID de la ability"
// @Success 200 {object} AbilityResponse
// @Failure 400 {string} string "ID de ability inválido"
// @Failure 404 {string} string "Ability no encontrada"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/abilities/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_ability") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver abilities")
		return
	}

	// Obtener el ID de la ability de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de ability inválido")
		return
	}

	// Obtener la ability
	ability, err := h.abilityUseCase.GetByID(r.Context(), id)
	if err != nil {
		// Verificar si el error es de ability no encontrada
		if err.Error() == "ability no encontrada" {
			response.Error(w, http.StatusNotFound, "Ability no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener la ability")
		return
	}

	// Preparar la respuesta
	resp := AbilityResponse{
		ID:         ability.ID,
		Name:       ability.Name,
		Title:      ability.Title,
		EntityID:   ability.EntityID,
		EntityType: ability.EntityType,
		OnlyOwned:  ability.OnlyOwned,
		Options:    ability.Options,
		Scope:      ability.Scope,
	}

	if ability.CreatedAt != nil {
		resp.CreatedAt = ability.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if ability.UpdatedAt != nil {
		resp.UpdatedAt = ability.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Update maneja la solicitud de actualización de una ability
// @Summary Actualizar ability
// @Description Actualiza una ability existente
// @Tags Abilities
// @Accept json
// @Produce json
// @Param id path int true "ID de la ability"
// @Param ability body UpdateAbilityRequest true "Datos de la ability"
// @Success 200 {object} AbilityResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 404 {string} string "Ability no encontrada"
// @Failure 409 {string} string "Nombre de ability ya está en uso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/abilities/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_ability") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar abilities")
		return
	}

	// Obtener el ID de la ability de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de ability inválido")
		return
	}

	// Decodificar la solicitud
	var req UpdateAbilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Actualizar la ability
	ability := &ability.Ability{
		ID:         id,
		Name:       req.Name,
		Title:      req.Title,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		OnlyOwned:  req.OnlyOwned,
		Options:    req.Options,
		Scope:      req.Scope,
	}

	if err := h.abilityUseCase.Update(r.Context(), ability); err != nil {
		// Verificar si el error es de ability no encontrada
		if err.Error() == "ability no encontrada" {
			response.Error(w, http.StatusNotFound, "Ability no encontrada")
			return
		}
		// Verificar si el error es de nombre duplicado
		if err.Error() == "nombre de ability ya está en uso" {
			response.Error(w, http.StatusConflict, "Nombre de ability ya está en uso")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar la ability")
		return
	}

	// Preparar la respuesta
	resp := AbilityResponse{
		ID:         ability.ID,
		Name:       ability.Name,
		Title:      ability.Title,
		EntityID:   ability.EntityID,
		EntityType: ability.EntityType,
		OnlyOwned:  ability.OnlyOwned,
		Options:    ability.Options,
		Scope:      ability.Scope,
	}

	if ability.CreatedAt != nil {
		resp.CreatedAt = ability.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if ability.UpdatedAt != nil {
		resp.UpdatedAt = ability.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete maneja la solicitud de eliminación de una ability
// @Summary Eliminar ability
// @Description Elimina una ability
// @Tags Abilities
// @Accept json
// @Produce json
// @Param id path int true "ID de la ability"
// @Success 204 "No Content"
// @Failure 400 {string} string "ID de ability inválido"
// @Failure 404 {string} string "Ability no encontrada"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/abilities/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "delete_ability") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para eliminar abilities")
		return
	}

	// Obtener el ID de la ability de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de ability inválido")
		return
	}

	// Eliminar la ability
	if err := h.abilityUseCase.Delete(r.Context(), id); err != nil {
		// Verificar si el error es de ability no encontrada
		if err.Error() == "ability no encontrada" {
			response.Error(w, http.StatusNotFound, "Ability no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar la ability")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List maneja la solicitud de listado de abilities
// @Summary Listar abilities
// @Description Obtiene una lista paginada de abilities
// @Tags Abilities
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param pageSize query int false "Tamaño de página (por defecto: 10)"
// @Param name query string false "Filtrar por nombre"
// @Param entityType query string false "Filtrar por tipo de entidad"
// @Param scope query int false "Filtrar por scope"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {string} string "Parámetros de consulta inválidos"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/abilities [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "list_abilities") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar abilities")
		return
	}

	// Obtener parámetros de consulta
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Construir filtros
	filters := make(map[string]interface{})

	if name := r.URL.Query().Get("name"); name != "" {
		filters["name"] = name
	}

	if entityType := r.URL.Query().Get("entityType"); entityType != "" {
		filters["entity_type"] = entityType
	}

	if scopeStr := r.URL.Query().Get("scope"); scopeStr != "" {
		scope, err := strconv.Atoi(scopeStr)
		if err == nil {
			filters["scope"] = scope
		}
	}

	// Obtener la lista de abilities
	abilities, total, err := h.abilityUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar abilities")
		return
	}

	// Preparar la respuesta
	var items []AbilityResponse
	for _, ability := range abilities {
		item := AbilityResponse{
			ID:         ability.ID,
			Name:       ability.Name,
			Title:      ability.Title,
			EntityID:   ability.EntityID,
			EntityType: ability.EntityType,
			OnlyOwned:  ability.OnlyOwned,
			Options:    ability.Options,
			Scope:      ability.Scope,
		}

		if ability.CreatedAt != nil {
			item.CreatedAt = ability.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		if ability.UpdatedAt != nil {
			item.UpdatedAt = ability.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	// Enviar respuesta paginada
	response.Paginated(w, items, page, pageSize, total)
}
