package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// JSON envía una respuesta JSON
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Error envía una respuesta de error en formato JSON
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, ErrorResponse{Error: message})
}

// Paginated envía una respuesta paginada en formato JSON
func Paginated(w http.ResponseWriter, items interface{}, page, pageSize, totalItems int) {
	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}

	response := PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	JSON(w, http.StatusOK, response)
}
