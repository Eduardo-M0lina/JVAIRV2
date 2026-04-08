package job_sms

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_sms"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	service job_sms.Service
}

func NewHandler(service job_sms.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	Recipients []string `json:"recipients"`
	Type       string   `json:"type"`
	Message    string   `json:"message"`
}

type ItemResponse struct {
	ID         int64    `json:"id"`
	JobID      int64    `json:"jobId"`
	Recipients []string `json:"recipients"`
	Type       string   `json:"type"`
	Message    string   `json:"message"`
	CreatedAt  string   `json:"createdAt,omitempty"`
	UpdatedAt  string   `json:"updatedAt,omitempty"`
}

func toResponse(item *job_sms.JobSMS) ItemResponse {
	resp := ItemResponse{ID: item.ID, JobID: item.JobID, Recipients: item.Recipients, Type: item.Type, Message: item.Message}
	if item.CreatedAt != nil {
		resp.CreatedAt = item.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if item.UpdatedAt != nil {
		resp.UpdatedAt = item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

// Create godoc
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de trabajo inválido")
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}
	item := &job_sms.JobSMS{JobID: jobID, Recipients: req.Recipients, Type: req.Type, Message: req.Message}
	if err := h.service.Create(r.Context(), item); err != nil {
		if err == job_sms.ErrJobNotFound {
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, toResponse(item))
}

// List godoc
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de trabajo inválido")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 15
	}
	offset := (page - 1) * limit
	items, total, err := h.service.List(r.Context(), jobID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar SMS enviados")
		return
	}
	result := make([]ItemResponse, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	response.Paginated(w, result, page, limit, int(total))
}

// Delete godoc
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == job_sms.ErrNotFound {
			response.Error(w, http.StatusNotFound, "SMS no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar SMS")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
