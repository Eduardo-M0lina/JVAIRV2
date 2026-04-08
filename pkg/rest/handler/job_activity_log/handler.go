package job_activity_log

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/job_activity_log"
	"github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service job_activity_log.Service
}

func NewHandler(service job_activity_log.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	JobID  int64  `json:"jobId"`
	Type   string `json:"type"`
	Log    string `json:"log"`
	UserID int64  `json:"userId"`
}

type ListResponse struct {
	Data  []*job_activity_log.JobActivityLog `json:"data"`
	Total int64                              `json:"total"`
	Limit int                                `json:"limit"`
	Page  int                                `json:"page"`
}

// Create godoc
// @Summary Create job activity log
// @Description Create a new activity log/note for a job
// @Tags JobActivities
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param request body CreateRequest true "Activity log data"
// @Success 201 {object} job_activity_log.JobActivityLog
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/activities [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.JobID = jobID

	log := &job_activity_log.JobActivityLog{
		JobID:  req.JobID,
		Type:   req.Type,
		Log:    req.Log,
		UserID: req.UserID,
	}

	if err := h.service.Create(r.Context(), log); err != nil {
		if err == job_activity_log.ErrJobNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job not found")
			return
		}
		if err == job_activity_log.ErrUserNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusCreated, log)
}

// List godoc
// @Summary List job activity logs
// @Description Get all activity logs for a specific job
// @Tags JobActivities
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param limit query int false "Limit" default(50)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/activities [get]
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

	logs, total, err := h.service.List(r.Context(), jobID, limit, offset)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListResponse{
		Data:  logs,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete job activity log
// @Description Delete an activity log by ID
// @Tags JobActivities
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Activity Log ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/activities/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid activity log ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_activity_log.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Activity log not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
