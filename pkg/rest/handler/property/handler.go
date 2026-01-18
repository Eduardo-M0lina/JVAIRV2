package property

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/property"
)

// Handler maneja las peticiones HTTP para propiedades
type Handler struct {
	useCase *property.UseCase
}

// NewHandler crea una nueva instancia del handler de propiedades
func NewHandler(useCase *property.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// CreatePropertyRequest representa la solicitud para crear una propiedad
type CreatePropertyRequest struct {
	CustomerID   int64   `json:"customerId"`
	PropertyCode *string `json:"propertyCode,omitempty"`
	Street       string  `json:"street"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Zip          string  `json:"zip"`
	Notes        *string `json:"notes,omitempty"`
}

// UpdatePropertyRequest representa la solicitud para actualizar una propiedad
type UpdatePropertyRequest struct {
	CustomerID   int64   `json:"customerId"`
	PropertyCode *string `json:"propertyCode,omitempty"`
	Street       string  `json:"street"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Zip          string  `json:"zip"`
	Notes        *string `json:"notes,omitempty"`
}

// PropertyResponse representa la respuesta de una propiedad
type PropertyResponse struct {
	ID           int64   `json:"id"`
	CustomerID   int64   `json:"customerId"`
	PropertyCode *string `json:"propertyCode,omitempty"`
	Street       string  `json:"street"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Zip          string  `json:"zip"`
	Notes        *string `json:"notes,omitempty"`
	Address      string  `json:"address"`
	Name         string  `json:"name"`
	CreatedAt    string  `json:"createdAt,omitempty"`
	UpdatedAt    string  `json:"updatedAt,omitempty"`
}

func toPropertyResponse(p *property.Property) PropertyResponse {
	resp := PropertyResponse{
		ID:           p.ID,
		CustomerID:   p.CustomerID,
		PropertyCode: p.PropertyCode,
		Street:       p.Street,
		City:         p.City,
		State:        p.State,
		Zip:          p.Zip,
		Notes:        p.Notes,
		Address:      p.GetAddress(),
		Name:         p.GetName(),
	}

	if p.CreatedAt != nil {
		resp.CreatedAt = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if p.UpdatedAt != nil {
		resp.UpdatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
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
