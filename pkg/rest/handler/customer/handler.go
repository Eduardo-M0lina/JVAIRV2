package customer

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/customer"
)

type Handler struct {
	useCase customer.Service
}

func NewHandler(useCase customer.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/customers", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

type CustomerResponse struct {
	ID             int64           `json:"id" example:"1"`
	Name           string          `json:"name" example:"ACME Corporation"`
	Email          *string         `json:"email,omitempty" example:"contact@acme.com"`
	Phone          *string         `json:"phone,omitempty" example:"+1-555-0100"`
	Mobile         *string         `json:"mobile,omitempty" example:"+1-555-0101"`
	Fax            *string         `json:"fax,omitempty" example:"+1-555-0102"`
	PhoneOther     *string         `json:"phoneOther,omitempty" example:"+1-555-0103"`
	Website        *string         `json:"website,omitempty" example:"https://acme.com"`
	ContactName    *string         `json:"contactName,omitempty" example:"John Doe"`
	ContactEmail   *string         `json:"contactEmail,omitempty" example:"john@acme.com"`
	ContactPhone   *string         `json:"contactPhone,omitempty" example:"+1-555-0104"`
	BillingAddress *BillingAddress `json:"billingAddress,omitempty"`
	WorkflowID     int64           `json:"workflowId" example:"5"`
	Notes          *string         `json:"notes,omitempty" example:"Important client notes"`
	CreatedAt      string          `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt      string          `json:"updatedAt,omitempty" example:"2024-01-18T14:20:00Z"`
}

type BillingAddress struct {
	Street *string `json:"street,omitempty" example:"123 Main St"`
	City   *string `json:"city,omitempty" example:"New York"`
	State  *string `json:"state,omitempty" example:"NY"`
	Zip    *string `json:"zip,omitempty" example:"10001"`
}

type CreateCustomerRequest struct {
	Name                 string  `json:"name" validate:"required" example:"ACME Corporation"`
	Email                *string `json:"email,omitempty" validate:"omitempty,email" example:"contact@acme.com"`
	Phone                *string `json:"phone,omitempty" example:"+1-555-0100"`
	Mobile               *string `json:"mobile,omitempty" example:"+1-555-0101"`
	Fax                  *string `json:"fax,omitempty" example:"+1-555-0102"`
	PhoneOther           *string `json:"phoneOther,omitempty" example:"+1-555-0103"`
	Website              *string `json:"website,omitempty" example:"https://acme.com"`
	ContactName          *string `json:"contactName,omitempty" example:"John Doe"`
	ContactEmail         *string `json:"contactEmail,omitempty" example:"john@acme.com"`
	ContactPhone         *string `json:"contactPhone,omitempty" example:"+1-555-0104"`
	BillingAddressStreet *string `json:"billingAddressStreet,omitempty" example:"123 Main St"`
	BillingAddressCity   *string `json:"billingAddressCity,omitempty" example:"New York"`
	BillingAddressState  *string `json:"billingAddressState,omitempty" example:"NY"`
	BillingAddressZip    *string `json:"billingAddressZip,omitempty" example:"10001"`
	WorkflowID           int64   `json:"workflowId" validate:"required,gt=0" example:"5"`
	Notes                *string `json:"notes,omitempty" example:"Important client notes"`
}

type UpdateCustomerRequest struct {
	Name                 string  `json:"name" validate:"required" example:"ACME Corporation"`
	Email                *string `json:"email,omitempty" validate:"omitempty,email" example:"contact@acme.com"`
	Phone                *string `json:"phone,omitempty" example:"+1-555-0100"`
	Mobile               *string `json:"mobile,omitempty" example:"+1-555-0101"`
	Fax                  *string `json:"fax,omitempty" example:"+1-555-0102"`
	PhoneOther           *string `json:"phoneOther,omitempty" example:"+1-555-0103"`
	Website              *string `json:"website,omitempty" example:"https://acme.com"`
	ContactName          *string `json:"contactName,omitempty" example:"John Doe"`
	ContactEmail         *string `json:"contactEmail,omitempty" example:"john@acme.com"`
	ContactPhone         *string `json:"contactPhone,omitempty" example:"+1-555-0104"`
	BillingAddressStreet *string `json:"billingAddressStreet,omitempty" example:"123 Main St"`
	BillingAddressCity   *string `json:"billingAddressCity,omitempty" example:"New York"`
	BillingAddressState  *string `json:"billingAddressState,omitempty" example:"NY"`
	BillingAddressZip    *string `json:"billingAddressZip,omitempty" example:"10001"`
	WorkflowID           int64   `json:"workflowId" validate:"required,gt=0" example:"5"`
	Notes                *string `json:"notes,omitempty" example:"Important client notes"`
}

func toCustomerResponse(c *customer.Customer) CustomerResponse {
	resp := CustomerResponse{
		ID:           c.ID,
		Name:         c.Name,
		Email:        c.Email,
		Phone:        c.Phone,
		Mobile:       c.Mobile,
		Fax:          c.Fax,
		PhoneOther:   c.PhoneOther,
		Website:      c.Website,
		ContactName:  c.ContactName,
		ContactEmail: c.ContactEmail,
		ContactPhone: c.ContactPhone,
		WorkflowID:   c.WorkflowID,
		Notes:        c.Notes,
	}

	if c.BillingAddressStreet != nil || c.BillingAddressCity != nil ||
		c.BillingAddressState != nil || c.BillingAddressZip != nil {
		resp.BillingAddress = &BillingAddress{
			Street: c.BillingAddressStreet,
			City:   c.BillingAddressCity,
			State:  c.BillingAddressState,
			Zip:    c.BillingAddressZip,
		}
	}

	if c.CreatedAt != nil {
		resp.CreatedAt = c.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if c.UpdatedAt != nil {
		resp.UpdatedAt = c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if workflowIDStr := r.URL.Query().Get("workflow_id"); workflowIDStr != "" {
		if workflowID, err := strconv.ParseInt(workflowIDStr, 10, 64); err == nil {
			filters["workflow_id"] = workflowID
		}
	}

	return filters
}
