package customer

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/your-org/jvairv2/pkg/domain/customer"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

var validate = validator.New()

// Create godoc
// @Summary Create a new customer
// @Description Create a new customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body CreateCustomerRequest true "Customer information"
// @Success 201 {object} CustomerResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /customers [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "Failed to decode request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		slog.WarnContext(r.Context(), "Validation failed",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c := &customer.Customer{
		Name:                 req.Name,
		Email:                req.Email,
		Phone:                req.Phone,
		Mobile:               req.Mobile,
		Fax:                  req.Fax,
		PhoneOther:           req.PhoneOther,
		Website:              req.Website,
		ContactName:          req.ContactName,
		ContactEmail:         req.ContactEmail,
		ContactPhone:         req.ContactPhone,
		BillingAddressStreet: req.BillingAddressStreet,
		BillingAddressCity:   req.BillingAddressCity,
		BillingAddressState:  req.BillingAddressState,
		BillingAddressZip:    req.BillingAddressZip,
		WorkflowID:           req.WorkflowID,
		Notes:                req.Notes,
	}

	if err := h.useCase.Create(r.Context(), c); err != nil {
		slog.ErrorContext(r.Context(), "Failed to create customer",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toCustomerResponse(c))
}
