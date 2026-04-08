package job_rate_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_rate_status"
	"github.com/your-org/jvairv2/pkg/rest/handler"
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
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status := &job_rate_status.JobRateStatus{
		Label: req.Label,
		Class: req.Class,
		Order: req.Order,
	}

	if err := h.service.Create(r.Context(), status); err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusCreated, status)
}

// List godoc
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.service.List(r.Context())
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, statuses)
}

// GetByID godoc
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid status ID")
		return
	}

	status, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == job_rate_status.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job rate status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, status)
}

// Update godoc
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid status ID")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
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
			handler.RespondWithError(w, http.StatusNotFound, "Job rate status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedStatus, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, updatedStatus)
}

// Delete godoc
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid status ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_rate_status.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job rate status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
