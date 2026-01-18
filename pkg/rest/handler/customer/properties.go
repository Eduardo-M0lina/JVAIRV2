package customer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// GetProperties godoc
// @Summary Get customer properties
// @Description Get all properties for a specific customer
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /customers/{id}/properties [get]
// @Security BearerAuth
func (h *Handler) GetProperties(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	customerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	// Verificar que el customer existe
	_, err = h.useCase.GetByID(r.Context(), customerID)
	if err != nil {
		if err.Error() == "not found" {
			response.Error(w, http.StatusNotFound, "Customer not found")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get customer",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}

	// Parsear parámetros de paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Crear filtros con customerId
	filters := map[string]interface{}{
		"customer_id": customerID,
	}

	// Agregar búsqueda si existe
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	// Obtener propiedades usando el use case de property
	properties, total, err := h.propertyUC.List(r.Context(), filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to list customer properties",
			slog.String("error", err.Error()),
			slog.Int64("customerId", customerID))
		response.Error(w, http.StatusInternalServerError, "Failed to list properties")
		return
	}

	// Convertir a response (necesitamos importar el handler de property para usar su función)
	items := make([]map[string]interface{}, len(properties))
	for i, p := range properties {
		items[i] = map[string]interface{}{
			"id":           p.ID,
			"customerId":   p.CustomerID,
			"propertyCode": p.PropertyCode,
			"street":       p.Street,
			"city":         p.City,
			"state":        p.State,
			"zip":          p.Zip,
			"notes":        p.Notes,
			"address":      p.GetAddress(),
			"name":         p.GetName(),
		}
		if p.CreatedAt != nil {
			items[i]["createdAt"] = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if p.UpdatedAt != nil {
			items[i]["updatedAt"] = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
