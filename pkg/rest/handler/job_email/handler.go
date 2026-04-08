package job_email

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/job_email"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service job_email.Service
}

func NewHandler(service job_email.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	Recipients []string `json:"recipients"`
	Type       string   `json:"type"`
}

type ItemResponse struct {
	ID         int64    `json:"id"`
	JobID      int64    `json:"jobId"`
	Recipients []string `json:"recipients"`
	Type       string   `json:"type"`
	CreatedAt  string   `json:"createdAt,omitempty"`
	UpdatedAt  string   `json:"updatedAt,omitempty"`
}

type ListResponse struct {
	Items      []ItemResponse `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalItems int            `json:"totalItems"`
	TotalPages int            `json:"totalPages"`
}

func toResponse(item *job_email.JobEmail) ItemResponse {
	resp := ItemResponse{ID: item.ID, JobID: item.JobID, Recipients: item.Recipients, Type: item.Type}
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
	item := &job_email.JobEmail{JobID: jobID, Recipients: req.Recipients, Type: req.Type}
	if err := h.service.Create(r.Context(), item); err != nil {
		if err == job_email.ErrJobNotFound {
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
		response.Error(w, http.StatusInternalServerError, "Error al listar emails enviados")
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
		if err == job_email.ErrNotFound {
			response.Error(w, http.StatusNotFound, "Email no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar email")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
