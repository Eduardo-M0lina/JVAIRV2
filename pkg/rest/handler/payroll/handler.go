package payroll

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/email"
	"github.com/your-org/jvairv2/pkg/domain/payroll"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las peticiones HTTP para payroll
type Handler struct {
	useCase      payroll.Service
	emailService email.Service
}

// NewHandler crea una nueva instancia del handler de payroll
func NewHandler(useCase payroll.Service, emailService email.Service) *Handler {
	return &Handler{useCase: useCase, emailService: emailService}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/payroll", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{userId}", h.GetUserPayroll)
		r.Put("/{userId}/mark-paid", h.MarkPaid)
		r.Put("/{userId}/mark-held", h.MarkHeld)
		r.Get("/{userId}/paystub", h.GetPaystub)
		r.Post("/{userId}/paystub/email", h.EmailPaystub)
	})
}

// List godoc
// @Summary Listar payroll de usuarios
// @Description Obtiene la lista de usuarios con sus rates de payroll agrupados por estado (unpaid, holding, paid)
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página (default: 1)"
// @Param pageSize query int false "Tamaño de página (default: 20)"
// @Param search query string false "Búsqueda por nombre o email de usuario"
// @Param userId query int false "Filtrar por ID de usuario específico"
// @Success 200 {object} payroll.PayrollListResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := payroll.PayrollFilters{
		Page:     1,
		PageSize: 20,
	}

	// Parsear parámetros de query
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			filters.Page = p
		}
	}
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			filters.PageSize = ps
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = &search
	}
	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		if uid, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filters.UserID = &uid
		}
	}

	result, err := h.useCase.ListPayroll(r.Context(), filters)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error listing payroll", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener payroll")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// GetUserPayroll godoc
// @Summary Obtener payroll de un usuario
// @Description Obtiene los rates de payroll de un usuario específico con filtros opcionales
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID del usuario"
// @Param status query string false "Filtrar por estado: unpaid, holding, paid, all (default: all)"
// @Param page query int false "Número de página (default: 1)"
// @Param pageSize query int false "Tamaño de página (default: 50)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll/{userId} [get]
func (h *Handler) GetUserPayroll(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	page := 1
	pageSize := 50
	var status *string

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status = &statusStr
	}

	rates, total, err := h.useCase.GetUserPayroll(r.Context(), userID, status, page, pageSize)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting user payroll",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"rates":      rates,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": totalPages,
	})
}

// MarkPaid godoc
// @Summary Marcar rates como pagados
// @Description Marca los rates especificados como pagados para un usuario
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID del usuario"
// @Param body body payroll.MarkRatesRequest true "IDs de los rates a marcar como pagados"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll/{userId}/mark-paid [put]
func (h *Handler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req payroll.MarkRatesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	if len(req.RateIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "Debe especificar al menos un rate ID")
		return
	}

	if err := h.useCase.MarkAsPaid(r.Context(), userID, req.RateIDs); err != nil {
		slog.ErrorContext(r.Context(), "Error marking rates as paid",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Rates marcados como pagados exitosamente",
	})
}

// MarkHeld godoc
// @Summary Marcar rates como retenidos
// @Description Marca los rates especificados como retenidos (holding) para un usuario
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID del usuario"
// @Param body body payroll.MarkRatesRequest true "IDs de los rates a marcar como retenidos"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll/{userId}/mark-held [put]
func (h *Handler) MarkHeld(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req payroll.MarkRatesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	if len(req.RateIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "Debe especificar al menos un rate ID")
		return
	}

	if err := h.useCase.MarkAsHolding(r.Context(), userID, req.RateIDs); err != nil {
		slog.ErrorContext(r.Context(), "Error marking rates as holding",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Rates marcados como retenidos exitosamente",
	})
}

// GetPaystub godoc
// @Summary Obtener recibo de pago
// @Description Obtiene el recibo de pago (paystub) de un usuario con sus rates pendientes
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID del usuario"
// @Success 200 {object} payroll.PaystubData
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll/{userId}/paystub [get]
func (h *Handler) GetPaystub(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	paystub, err := h.useCase.GetPaystub(r.Context(), userID)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting paystub",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, paystub)
}

// EmailPaystub godoc
// @Summary Enviar recibo de pago por email
// @Description Envía el recibo de pago (paystub) de un usuario a las direcciones de email especificadas
// @Tags Payroll
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID del usuario"
// @Param body body map[string]interface{} true "Direcciones de email (campo 'emails')"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payroll/{userId}/paystub/email [post]
func (h *Handler) EmailPaystub(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req struct {
		Emails string `json:"emails"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	if req.Emails == "" {
		response.Error(w, http.StatusBadRequest, "Debe especificar al menos una dirección de email")
		return
	}

	// Separar emails por coma
	emails := strings.Split(req.Emails, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}

	// Enviar paystub por email usando el servicio de email
	err = h.emailService.SendPayStubEmail(r.Context(), userID, emails)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error sending paystub email",
			slog.String("error", err.Error()),
			slog.Int64("userId", userID),
			slog.Any("emails", emails))
		response.Error(w, http.StatusInternalServerError, "Error al enviar el paystub por email")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Paystub enviado exitosamente a " + req.Emails,
	})
}
