package job_rate

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_rate"
	"github.com/your-org/jvairv2/pkg/rest/handler"
)

type Handler struct {
	service job_rate.Service
}

func NewHandler(service job_rate.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	UserID          int64   `json:"userId"`
	JobRateStatusID int64   `json:"jobRateStatusId"`
	SalePrice       float64 `json:"salePrice"`
	RatePercent     float64 `json:"ratePercent"`
	RateFlat        float64 `json:"rateFlat"`
	TechParts       float64 `json:"techParts"`
	CompanyParts    float64 `json:"companyParts"`
	PartsReplaced   *string `json:"partsReplaced,omitempty"`
	Deduction       float64 `json:"deduction"`
	Notes           *string `json:"notes,omitempty"`
}

type UpdateRequest struct {
	UserID          int64   `json:"userId"`
	JobRateStatusID int64   `json:"jobRateStatusId"`
	SalePrice       float64 `json:"salePrice"`
	RatePercent     float64 `json:"ratePercent"`
	RateFlat        float64 `json:"rateFlat"`
	TechParts       float64 `json:"techParts"`
	CompanyParts    float64 `json:"companyParts"`
	PartsReplaced   *string `json:"partsReplaced,omitempty"`
	Deduction       float64 `json:"deduction"`
	Paid            bool    `json:"paid"`
	Notes           *string `json:"notes,omitempty"`
}

type CalculatePaymentRequest struct {
	SalePrice    float64 `json:"salePrice"`
	RatePercent  float64 `json:"ratePercent"`
	RateFlat     float64 `json:"rateFlat"`
	TechParts    float64 `json:"techParts"`
	CompanyParts float64 `json:"companyParts"`
	Deduction    float64 `json:"deduction"`
}

type CalculatePaymentResponse struct {
	Payment float64 `json:"payment"`
}

type ListResponse struct {
	Data  []*job_rate.JobRate `json:"data"`
	Total int64               `json:"total"`
	Limit int                 `json:"limit"`
	Page  int                 `json:"page"`
}

// Create godoc
// @Summary Create job rate
// @Description Create a new rate/commission for a job
// @Tags JobRates
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param request body CreateRequest true "Rate data"
// @Success 201 {object} job_rate.JobRate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/rates [post]
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

	rate := &job_rate.JobRate{
		JobID:           jobID,
		UserID:          req.UserID,
		JobRateStatusID: req.JobRateStatusID,
		SalePrice:       req.SalePrice,
		RatePercent:     req.RatePercent,
		RateFlat:        req.RateFlat,
		TechParts:       req.TechParts,
		CompanyParts:    req.CompanyParts,
		PartsReplaced:   req.PartsReplaced,
		Deduction:       req.Deduction,
		Notes:           req.Notes,
	}

	if err := h.service.Create(r.Context(), rate); err != nil {
		if err == job_rate.ErrJobNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job not found")
			return
		}
		if err == job_rate.ErrUserNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if err == job_rate.ErrJobRateStatusNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job rate status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusCreated, rate)
}

// List godoc
// @Summary List job rates
// @Description Get all rates for a specific job
// @Tags JobRates
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param limit query int false "Limit" default(50)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/rates [get]
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

	rates, total, err := h.service.List(r.Context(), jobID, limit, offset)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListResponse{
		Data:  rates,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}

// Update godoc
// @Summary Update job rate
// @Description Update a rate by ID
// @Tags JobRates
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Rate ID"
// @Param request body UpdateRequest true "Rate data"
// @Success 200 {object} job_rate.JobRate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/rates/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid rate ID")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rate := &job_rate.JobRate{
		ID:              id,
		UserID:          req.UserID,
		JobRateStatusID: req.JobRateStatusID,
		SalePrice:       req.SalePrice,
		RatePercent:     req.RatePercent,
		RateFlat:        req.RateFlat,
		TechParts:       req.TechParts,
		CompanyParts:    req.CompanyParts,
		PartsReplaced:   req.PartsReplaced,
		Deduction:       req.Deduction,
		Paid:            req.Paid,
		Notes:           req.Notes,
	}

	if err := h.service.Update(r.Context(), rate); err != nil {
		if err == job_rate.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Rate not found")
			return
		}
		if err == job_rate.ErrUserNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if err == job_rate.ErrJobRateStatusNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Job rate status not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedRate, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, updatedRate)
}

// Delete godoc
// @Summary Delete job rate
// @Description Soft delete a rate by ID
// @Tags JobRates
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Rate ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/rates/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid rate ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_rate.ErrNotFound {
			handler.RespondWithError(w, http.StatusNotFound, "Rate not found")
			return
		}
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CalculatePayment godoc
// @Summary Calculate rate payment
// @Description Calculate payment amount based on rate formula
// @Tags JobRates
// @Accept json
// @Produce json
// @Param request body CalculatePaymentRequest true "Payment calculation data"
// @Success 200 {object} CalculatePaymentResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/calculate-rate-payment [post]
// @Security BearerAuth
func (h *Handler) CalculatePayment(w http.ResponseWriter, r *http.Request) {
	var req CalculatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment := job_rate.CalculatePayment(
		req.SalePrice,
		req.RatePercent,
		req.RateFlat,
		req.TechParts,
		req.CompanyParts,
		req.Deduction,
	)

	response := CalculatePaymentResponse{
		Payment: payment,
	}

	handler.RespondWithJSON(w, http.StatusOK, response)
}
