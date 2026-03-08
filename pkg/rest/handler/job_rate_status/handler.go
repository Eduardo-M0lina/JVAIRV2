package job_rate_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_rate_status"
)

type Handler struct {
	service job_rate_status.Service
}

func NewHandler(service job_rate_status.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	Label string  `json:"label"`
	Class *string `json:"class,omitempty"`
	Order int     `json:"order"`
}

type UpdateRequest struct {
	Label string  `json:"label"`
	Class *string `json:"class,omitempty"`
	Order int     `json:"order"`
}

// Create godoc
// @Summary Create job rate status
// @Description Create a new job rate status
// @Tags JobRateStatuses
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Job rate status data"
// @Success 201 {object} job_rate_status.JobRateStatus
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/job-rate-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := &job_rate_status.JobRateStatus{
		Label: req.Label,
		Class: req.Class,
		Order: req.Order,
	}

	if err := h.service.Create(r.Context(), status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(status)
}

// List godoc
// @Summary List job rate statuses
// @Description Get all job rate statuses
// @Tags JobRateStatuses
// @Accept json
// @Produce json
// @Success 200 {array} job_rate_status.JobRateStatus
// @Failure 500 {object} map[string]string
// @Router /api/v1/job-rate-statuses [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(statuses)
}

// GetByID godoc
// @Summary Get job rate status by ID
// @Description Get a job rate status by ID
// @Tags JobRateStatuses
// @Accept json
// @Produce json
// @Param id path int true "Job Rate Status ID"
// @Success 200 {object} job_rate_status.JobRateStatus
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/job-rate-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == job_rate_status.ErrNotFound {
			http.Error(w, "Job rate status not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// Update godoc
// @Summary Update job rate status
// @Description Update a job rate status by ID
// @Tags JobRateStatuses
// @Accept json
// @Produce json
// @Param id path int true "Job Rate Status ID"
// @Param request body UpdateRequest true "Job rate status data"
// @Success 200 {object} job_rate_status.JobRateStatus
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/job-rate-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := &job_rate_status.JobRateStatus{
		ID:    id,
		Label: req.Label,
		Class: req.Class,
		Order: req.Order,
	}

	if err := h.service.Update(r.Context(), status); err != nil {
		if err == job_rate_status.ErrNotFound {
			http.Error(w, "Job rate status not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedStatus, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedStatus)
}

// Delete godoc
// @Summary Delete job rate status
// @Description Soft delete a job rate status by ID
// @Tags JobRateStatuses
// @Accept json
// @Produce json
// @Param id path int true "Job Rate Status ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/job-rate-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_rate_status.ErrNotFound {
			http.Error(w, "Job rate status not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
