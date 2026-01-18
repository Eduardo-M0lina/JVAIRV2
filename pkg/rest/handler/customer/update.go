package customer

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/customer"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update godoc
// @Summary Update a customer
// @Description Update an existing customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body UpdateCustomerRequest true "Customer information"
// @Success 200 {object} CustomerResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /customers/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var req UpdateCustomerRequest
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
		ID:                   id,
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

	if err := h.useCase.Update(r.Context(), c); err != nil {
		slog.ErrorContext(r.Context(), "Failed to update customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedCustomer, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get updated customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		response.Error(w, http.StatusInternalServerError, "Failed to get updated customer")
		return
	}

	response.JSON(w, http.StatusOK, toCustomerResponse(updatedCustomer))
}
