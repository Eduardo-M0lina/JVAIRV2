package job_task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/email"
	"github.com/angumol/jvairv2/pkg/domain/job_task"
	"github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      job_task.Service
	emailService email.Service
}

func NewHandler(service job_task.Service, emailService email.Service) *Handler {
	return &Handler{service: service, emailService: emailService}
}

type CreateRequest struct {
	UserID       int64   `json:"userId"`
	Task         string  `json:"task"`
	TaskStatusID int64   `json:"taskStatusId"`
	DueDate      *string `json:"dueDate,omitempty"`
}

type UpdateRequest struct {
	UserID       int64   `json:"userId"`
	Task         string  `json:"task"`
	TaskStatusID int64   `json:"taskStatusId"`
	DueDate      *string `json:"dueDate,omitempty"`
}

type ListResponse struct {
	Data  []*job_task.JobTask `json:"data"`
	Total int64               `json:"total"`
	Limit int                 `json:"limit"`
	Page  int                 `json:"page"`
}

// Create godoc
// @Summary Create job task
// @Description Create a new task for a job
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param request body CreateRequest true "Task data"
// @Success 201 {object} job_task.JobTask
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/tasks [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "invalid job ID")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task := &job_task.JobTask{
		JobID:        jobID,
		UserID:       req.UserID,
		Task:         req.Task,
		TaskStatusID: req.TaskStatusID,
	}

	if err := h.service.Create(ctx, task); err != nil {
		if err == job_task.ErrJobNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job not found")
			return
		}
		if err == job_task.ErrUserNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if err == job_task.ErrTaskStatusNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Task status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusCreated, task)
}

// List godoc
// @Summary List job tasks
// @Description Get all tasks for a specific job
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param limit query int false "Limit" default(50)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/tasks [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	tasks, total, err := h.service.List(r.Context(), jobID, limit, offset)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListResponse{
		Data:  tasks,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}

// ListAll godoc
// @Summary List all tasks
// @Description Get all tasks across all jobs (global view)
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks [get]
// @Security BearerAuth
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	tasks, total, err := h.service.ListAll(r.Context(), limit, offset)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListResponse{
		Data:  tasks,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}

// Update godoc
// @Summary Update job task
// @Description Update a task by ID
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Task ID"
// @Param request body UpdateRequest true "Task data"
// @Success 200 {object} job_task.JobTask
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/tasks/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task := &job_task.JobTask{
		ID:           id,
		UserID:       req.UserID,
		Task:         req.Task,
		TaskStatusID: req.TaskStatusID,
	}

	if err := h.service.Update(r.Context(), task); err != nil {
		if err == job_task.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		if err == job_task.ErrUserNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if err == job_task.ErrTaskStatusNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Task status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedTask, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, updatedTask)
}

// Delete godoc
// @Summary Delete job task
// @Description Soft delete a task by ID
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/tasks/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_task.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
