package supervisor

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

type Handler struct {
	useCase supervisor.Service
}

func NewHandler(useCase supervisor.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/supervisors", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	// Sub-ruta: supervisores de un customer específico
	r.Get("/customers/{customerId}/supervisors", h.ListByCustomer)
}

type SupervisorResponse struct {
	ID         int64   `json:"id" example:"1"`
	CustomerID int64   `json:"customerId" example:"10"`
	Name       string  `json:"name" example:"John Doe"`
	Phone      *string `json:"phone,omitempty" example:"+1-555-0100"`
	Email      *string `json:"email,omitempty" example:"john@example.com"`
	CreatedAt  string  `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt  string  `json:"updatedAt,omitempty" example:"2024-01-18T14:20:00Z"`
}

type CreateSupervisorRequest struct {
	CustomerID int64   `json:"customerId" example:"10"`
	Name       string  `json:"name" example:"John Doe"`
	Phone      *string `json:"phone,omitempty" example:"+1-555-0100"`
	Email      *string `json:"email,omitempty" example:"john@example.com"`
}

type UpdateSupervisorRequest struct {
	CustomerID int64   `json:"customerId" example:"10"`
	Name       string  `json:"name" example:"John Doe"`
	Phone      *string `json:"phone,omitempty" example:"+1-555-0100"`
	Email      *string `json:"email,omitempty" example:"john@example.com"`
}

func toSupervisorResponse(s *supervisor.Supervisor) SupervisorResponse {
	resp := SupervisorResponse{
		ID:         s.ID,
		CustomerID: s.CustomerID,
		Name:       s.Name,
		Phone:      s.Phone,
		Email:      s.Email,
	}

	if s.CreatedAt != nil {
		resp.CreatedAt = s.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if s.UpdatedAt != nil {
		resp.UpdatedAt = s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if customerIDStr := r.URL.Query().Get("customerId"); customerIDStr != "" {
		if customerID, err := strconv.ParseInt(customerIDStr, 10, 64); err == nil {
			filters["customer_id"] = customerID
		}
	}

	return filters
}
