package job_task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/email"
	"github.com/your-org/jvairv2/pkg/domain/job_task"
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
		http.Error(w, `{"error":"invalid job ID"}`, http.StatusBadRequest)
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if err == job_task.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err == job_task.ErrTaskStatusNotFound {
			http.Error(w, "Task status not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
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
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListResponse{
		Data:  tasks,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListResponse{
		Data:  tasks,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
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
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if err == job_task.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err == job_task.ErrTaskStatusNotFound {
			http.Error(w, "Task status not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedTask, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedTask)
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
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_task.ErrNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
