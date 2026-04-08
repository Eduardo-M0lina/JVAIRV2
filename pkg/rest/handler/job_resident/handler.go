package job_resident

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/job_resident"
	"github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service job_resident.Service
}

func NewHandler(service job_resident.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	Name        string  `json:"name"`
	MobilePhone *string `json:"mobilePhone,omitempty"`
	HomePhone   *string `json:"homePhone,omitempty"`
	Email       *string `json:"email,omitempty"`
}

type UpdateRequest struct {
	Name        string  `json:"name"`
	MobilePhone *string `json:"mobilePhone,omitempty"`
	HomePhone   *string `json:"homePhone,omitempty"`
	Email       *string `json:"email,omitempty"`
}

type ListResponse struct {
	Data  []*job_resident.JobResident `json:"data"`
	Total int64                       `json:"total"`
	Limit int                         `json:"limit"`
	Page  int                         `json:"page"`
}

// Create godoc
// @Summary Create job resident
// @Description Create a new resident contact for a job
// @Tags JobResidents
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param request body CreateRequest true "Resident data"
// @Success 201 {object} job_resident.JobResident
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/residents [post]
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

	resident := &job_resident.JobResident{
		JobID:       jobID,
		Name:        req.Name,
		MobilePhone: req.MobilePhone,
		HomePhone:   req.HomePhone,
		Email:       req.Email,
	}

	if err := h.service.Create(r.Context(), resident); err != nil {
		if err == job_resident.ErrJobNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusCreated, resident)
}

// List godoc
// @Summary List job residents
// @Description Get all residents for a specific job
// @Tags JobResidents
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param limit query int false "Limit" default(50)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/residents [get]
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

	residents, total, err := h.service.List(r.Context(), jobID, limit, offset)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListResponse{
		Data:  residents,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}

// Update godoc
// @Summary Update job resident
// @Description Update a resident by ID
// @Tags JobResidents
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Resident ID"
// @Param request body UpdateRequest true "Resident data"
// @Success 200 {object} job_resident.JobResident
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/residents/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid resident ID")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resident := &job_resident.JobResident{
		ID:          id,
		Name:        req.Name,
		MobilePhone: req.MobilePhone,
		HomePhone:   req.HomePhone,
		Email:       req.Email,
	}

	if err := h.service.Update(r.Context(), resident); err != nil {
		if err == job_resident.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Resident not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedResident, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, updatedResident)
}

// Delete godoc
// @Summary Delete job resident
// @Description Soft delete a resident by ID
// @Tags JobResidents
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Resident ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/residents/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid resident ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_resident.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Resident not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
