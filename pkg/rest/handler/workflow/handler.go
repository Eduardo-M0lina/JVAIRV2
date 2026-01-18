package workflow

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/workflow"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las solicitudes HTTP relacionadas con workflows
type Handler struct {
	workflowUseCase *workflow.UseCase
	validate        *validator.Validate
}

// NewHandler crea una nueva instancia del manejador de workflows
func NewHandler(workflowUseCase *workflow.UseCase) *Handler {
	return &Handler{
		workflowUseCase: workflowUseCase,
		validate:        validator.New(),
	}
}

// CreateWorkflowRequest representa la solicitud para crear un workflow
type CreateWorkflowRequest struct {
	Name     string  `json:"name" validate:"required"`
	Notes    *string `json:"notes,omitempty"`
	IsActive bool    `json:"isActive"`
	Statuses []int64 `json:"statuses,omitempty"`
}

// UpdateWorkflowRequest representa la solicitud para actualizar un workflow
type UpdateWorkflowRequest struct {
	Name     string  `json:"name" validate:"required"`
	Notes    *string `json:"notes,omitempty"`
	IsActive bool    `json:"isActive"`
	Statuses []int64 `json:"statuses,omitempty"`
}

// WorkflowResponse representa la respuesta de un workflow
type WorkflowResponse struct {
	ID       int64                    `json:"id"`
	Name     string                   `json:"name"`
	Notes    *string                  `json:"notes,omitempty"`
	IsActive bool                     `json:"isActive"`
	Statuses []WorkflowStatusResponse `json:"statuses,omitempty"`
}

// WorkflowStatusResponse representa la respuesta de un status de workflow
type WorkflowStatusResponse struct {
	JobStatusID int64  `json:"jobStatusId"`
	Order       int    `json:"order"`
	StatusName  string `json:"statusName,omitempty"`
}

// List maneja la solicitud de listado de workflows
// @Summary Listar workflows
// @Description Obtiene una lista paginada de workflows con filtros opcionales
// @Tags Workflows
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param page_size query int false "Tamaño de página" default(10)
// @Param name query string false "Filtrar por nombre"
// @Param is_active query bool false "Filtrar por estado activo"
// @Param search query string false "Búsqueda en nombre y notas"
// @Success 200 {object} response.PaginatedResponse
// @Failure 403 {string} string "No tiene permisos para listar workflows"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para listar workflows")
		return
	}

	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Obtener filtros
	filters := workflow.Filters{
		Name:   r.URL.Query().Get("name"),
		Search: r.URL.Query().Get("search"),
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive, _ := strconv.ParseBool(isActiveStr)
		filters.IsActive = &isActive
	}

	// Obtener workflows
	workflows, total, err := h.workflowUseCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		log.Printf("ERROR al listar workflows: %v", err)
		response.Error(w, http.StatusInternalServerError, "Error al listar workflows")
		return
	}

	// Preparar la respuesta
	var items []WorkflowResponse
	for _, wf := range workflows {
		items = append(items, WorkflowResponse{
			ID:       wf.ID,
			Name:     wf.Name,
			Notes:    wf.Notes,
			IsActive: wf.IsActive,
		})
	}

	resp := response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	response.JSON(w, http.StatusOK, resp)
}

// Get maneja la solicitud de obtención de un workflow por ID
// @Summary Obtener workflow
// @Description Obtiene un workflow por su ID incluyendo sus statuses
// @Tags Workflows
// @Accept json
// @Produce json
// @Param id path int true "ID del workflow"
// @Success 200 {object} WorkflowResponse
// @Failure 400 {string} string "ID de workflow inválido"
// @Failure 403 {string} string "No tiene permisos para ver workflows"
// @Failure 404 {string} string "Workflow no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "view_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para ver workflows")
		return
	}

	// Obtener el ID del workflow de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de workflow inválido")
		return
	}

	// Obtener el workflow
	wf, err := h.workflowUseCase.GetByID(r.Context(), id)
	if err != nil {
		if err == workflow.ErrWorkflowNotFound {
			response.Error(w, http.StatusNotFound, "Workflow no encontrado")
			return
		}
		log.Printf("ERROR al obtener workflow por ID %d: %v", id, err)
		response.Error(w, http.StatusInternalServerError, "Error al obtener el workflow: "+err.Error())
		return
	}

	// Preparar la respuesta
	var statuses []WorkflowStatusResponse
	for _, s := range wf.Statuses {
		statuses = append(statuses, WorkflowStatusResponse{
			JobStatusID: s.JobStatusID,
			Order:       s.Order,
			StatusName:  s.StatusName,
		})
	}

	resp := WorkflowResponse{
		ID:       wf.ID,
		Name:     wf.Name,
		Notes:    wf.Notes,
		IsActive: wf.IsActive,
		Statuses: statuses,
	}

	response.JSON(w, http.StatusOK, resp)
}

// Create maneja la solicitud de creación de un workflow
// @Summary Crear workflow
// @Description Crea un nuevo workflow con sus statuses asociados. Ejemplo: {"name": "Workflow de Mantenimiento", "notes": "Flujo para trabajos de mantenimiento", "is_active": true, "statuses": [1, 2, 3, 4]}
// @Tags Workflows
// @Accept json
// @Produce json
// @Param workflow body CreateWorkflowRequest true "Datos del workflow a crear"
// @Success 201 {object} WorkflowResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 403 {string} string "No tiene permisos para crear workflows"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para crear workflows")
		return
	}

	// Decodificar la solicitud
	var req CreateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Crear el workflow
	wf := &workflow.Workflow{
		Name:     req.Name,
		Notes:    req.Notes,
		IsActive: req.IsActive,
	}

	if err := h.workflowUseCase.Create(r.Context(), wf, req.Statuses); err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al crear el workflow: "+err.Error())
		return
	}

	// Obtener el workflow completo con sus statuses
	created, err := h.workflowUseCase.GetByID(r.Context(), wf.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener el workflow creado")
		return
	}

	// Preparar la respuesta
	var statuses []WorkflowStatusResponse
	for _, s := range created.Statuses {
		statuses = append(statuses, WorkflowStatusResponse{
			JobStatusID: s.JobStatusID,
			Order:       s.Order,
			StatusName:  s.StatusName,
		})
	}

	resp := WorkflowResponse{
		ID:       created.ID,
		Name:     created.Name,
		Notes:    created.Notes,
		IsActive: created.IsActive,
		Statuses: statuses,
	}

	response.JSON(w, http.StatusCreated, resp)
}

// Update maneja la solicitud de actualización de un workflow
// @Summary Actualizar workflow
// @Description Actualiza un workflow existente y sus statuses asociados
// @Tags Workflows
// @Accept json
// @Produce json
// @Param id path int true "ID del workflow"
// @Param workflow body UpdateWorkflowRequest true "Datos del workflow a actualizar"
// @Success 200 {object} WorkflowResponse
// @Failure 400 {string} string "Error al decodificar la solicitud o datos inválidos"
// @Failure 403 {string} string "No tiene permisos para actualizar workflows"
// @Failure 404 {string} string "Workflow no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "update_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para actualizar workflows")
		return
	}

	// Obtener el ID del workflow de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de workflow inválido")
		return
	}

	// Decodificar la solicitud
	var req UpdateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Validar la solicitud
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	// Actualizar el workflow
	wf := &workflow.Workflow{
		ID:       id,
		Name:     req.Name,
		Notes:    req.Notes,
		IsActive: req.IsActive,
	}

	if err := h.workflowUseCase.Update(r.Context(), wf, req.Statuses); err != nil {
		if err == workflow.ErrWorkflowNotFound {
			response.Error(w, http.StatusNotFound, "Workflow no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el workflow: "+err.Error())
		return
	}

	// Obtener el workflow actualizado con sus statuses
	updated, err := h.workflowUseCase.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener el workflow actualizado")
		return
	}

	// Preparar la respuesta
	var statuses []WorkflowStatusResponse
	for _, s := range updated.Statuses {
		statuses = append(statuses, WorkflowStatusResponse{
			JobStatusID: s.JobStatusID,
			Order:       s.Order,
			StatusName:  s.StatusName,
		})
	}

	resp := WorkflowResponse{
		ID:       updated.ID,
		Name:     updated.Name,
		Notes:    updated.Notes,
		IsActive: updated.IsActive,
		Statuses: statuses,
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete maneja la solicitud de eliminación de un workflow
// @Summary Eliminar workflow
// @Description Elimina un workflow y sus relaciones con job_statuses
// @Tags Workflows
// @Accept json
// @Produce json
// @Param id path int true "ID del workflow"
// @Success 200 {string} string "Workflow eliminado exitosamente"
// @Failure 400 {string} string "ID de workflow inválido"
// @Failure 403 {string} string "No tiene permisos para eliminar workflows"
// @Failure 404 {string} string "Workflow no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "delete_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para eliminar workflows")
		return
	}

	// Obtener el ID del workflow de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de workflow inválido")
		return
	}

	// Eliminar el workflow
	if err := h.workflowUseCase.Delete(r.Context(), id); err != nil {
		if err == workflow.ErrWorkflowNotFound {
			response.Error(w, http.StatusNotFound, "Workflow no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el workflow")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Workflow eliminado exitosamente"})
}

// Duplicate maneja la solicitud de duplicación de un workflow
// @Summary Duplicar workflow
// @Description Duplica un workflow existente con todos sus statuses
// @Tags Workflows
// @Accept json
// @Produce json
// @Param id path int true "ID del workflow a duplicar"
// @Success 201 {object} WorkflowResponse
// @Failure 400 {string} string "ID de workflow inválido"
// @Failure 403 {string} string "No tiene permisos para duplicar workflows"
// @Failure 404 {string} string "Workflow no encontrado"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/v1/workflows/{id}/duplicate [post]
// @Security BearerAuth
func (h *Handler) Duplicate(w http.ResponseWriter, r *http.Request) {
	// Verificar permisos
	if !middleware.HasAbility(r.Context(), "create_workflow") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para duplicar workflows")
		return
	}

	// Obtener el ID del workflow de la URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de workflow inválido")
		return
	}

	// Duplicar el workflow
	duplicated, err := h.workflowUseCase.Duplicate(r.Context(), id)
	if err != nil {
		if err == workflow.ErrWorkflowNotFound {
			response.Error(w, http.StatusNotFound, "Workflow no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al duplicar el workflow")
		return
	}

	// Preparar la respuesta
	var statuses []WorkflowStatusResponse
	for _, s := range duplicated.Statuses {
		statuses = append(statuses, WorkflowStatusResponse{
			JobStatusID: s.JobStatusID,
			Order:       s.Order,
			StatusName:  s.StatusName,
		})
	}

	resp := WorkflowResponse{
		ID:       duplicated.ID,
		Name:     duplicated.Name,
		Notes:    duplicated.Notes,
		IsActive: duplicated.IsActive,
		Statuses: statuses,
	}

	response.JSON(w, http.StatusCreated, resp)
}
