package alert

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/alert"
	"github.com/your-org/jvairv2/pkg/domain/user"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

type Handler struct {
	service alert.Service
}

func NewHandler(service alert.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	UserID       *int64 `json:"userId,omitempty"`
	AlertType    string `json:"alertType"`
	EntityID     int64  `json:"entityId"`
	EntityType   string `json:"entityType"`
	MessageLevel string `json:"messageLevel"`
	Message      string `json:"message"`
}

type ListResponse struct {
	Data  []*alert.Alert `json:"data"`
	Total int64          `json:"total"`
	Limit int            `json:"limit"`
	Page  int            `json:"page"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type MarkAllReadResponse struct {
	Updated int64 `json:"updated"`
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Put("/read-all", h.MarkAllRead)
		r.Get("/unread-count", h.UnreadCount)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}/read", h.MarkAsRead)
		r.Delete("/{id}", h.Delete)
	})
}

// List godoc
// @Summary List alerts
// @Description Get all alerts with optional filters and pagination
// @Tags Alerts
// @Accept json
// @Produce json
// @Param userId query int false "Filter by user ID"
// @Param isRead query bool false "Filter by read status (true/false)"
// @Param alertType query string false "Filter by alert type"
// @Param entityType query string false "Filter by entity type"
// @Param limit query int false "Limit" default(15)
// @Param page query int false "Page" default(1)
// @Success 200 {object} ListResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var filters alert.ListFilters

	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filters.UserID = &userID
		}
	}

	if isReadStr := r.URL.Query().Get("isRead"); isReadStr != "" {
		isRead := isReadStr == "true"
		filters.IsRead = &isRead
	}

	if alertType := r.URL.Query().Get("alertType"); alertType != "" {
		filters.AlertType = &alertType
	}

	if entityType := r.URL.Query().Get("entityType"); entityType != "" {
		filters.EntityType = &entityType
	}

	limit := 15
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

	alerts, total, err := h.service.List(r.Context(), filters, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListResponse{
		Data:  alerts,
		Total: total,
		Limit: limit,
		Page:  page,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GetByID godoc
// @Summary Get alert by ID
// @Description Get a single alert by its ID
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} alert.Alert
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	a, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == alert.ErrNotFound {
			http.Error(w, "Alert not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a)
}

// Create godoc
// @Summary Create alert
// @Description Create a new alert (typically used by system events)
// @Tags Alerts
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Alert data"
// @Success 201 {object} alert.Alert
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	a := &alert.Alert{
		UserID:       req.UserID,
		AlertType:    req.AlertType,
		EntityID:     req.EntityID,
		EntityType:   req.EntityType,
		MessageLevel: req.MessageLevel,
		Message:      req.Message,
		IsRead:       false,
	}

	if err := h.service.Create(r.Context(), a); err != nil {
		if err == alert.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

// MarkAsRead godoc
// @Summary Mark alert as read
// @Description Mark a single alert as read
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts/{id}/read [put]
// @Security BearerAuth
func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkAsRead(r.Context(), id); err != nil {
		if err == alert.ErrNotFound {
			http.Error(w, "Alert not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkAllRead godoc
// @Summary Mark all alerts as read
// @Description Mark all unread alerts for the authenticated user as read
// @Tags Alerts
// @Accept json
// @Produce json
// @Success 200 {object} MarkAllReadResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts/read-all [put]
// @Security BearerAuth
func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(middleware.UserContextKey).(*user.User)
	if !ok || u == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	updated, err := h.service.MarkAllRead(r.Context(), u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := MarkAllReadResponse{Updated: updated}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// UnreadCount godoc
// @Summary Get unread alerts count
// @Description Get the count of unread alerts for the authenticated user
// @Tags Alerts
// @Accept json
// @Produce json
// @Success 200 {object} UnreadCountResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts/unread-count [get]
// @Security BearerAuth
func (h *Handler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(middleware.UserContextKey).(*user.User)
	if !ok || u == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := h.service.UnreadCount(r.Context(), u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UnreadCountResponse{Count: count}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Delete godoc
// @Summary Delete alert
// @Description Delete an alert by ID (hard delete)
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/alerts/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == alert.ErrNotFound {
			http.Error(w, "Alert not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
