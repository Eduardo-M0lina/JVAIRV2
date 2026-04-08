package search

// SearchResult representa un resultado de búsqueda individual
type SearchResult struct {
	ID          int64   `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Subtitle    *string `json:"subtitle,omitempty"`
	Description *string `json:"description,omitempty"`
}

// JobSearchResult representa un job encontrado en la búsqueda
type JobSearchResult struct {
	ID             int64   `json:"id"`
	WorkOrder      *string `json:"workOrder,omitempty"`
	PropertyStreet *string `json:"propertyStreet,omitempty"`
	PropertyCity   *string `json:"propertyCity,omitempty"`
	CustomerName   *string `json:"customerName,omitempty"`
	StatusName     *string `json:"statusName,omitempty"`
	TechnicianName *string `json:"technicianName,omitempty"`
	Closed         bool    `json:"closed"`
}

// CustomerSearchResult representa un customer encontrado en la búsqueda
type CustomerSearchResult struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Email   *string `json:"email,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Website *string `json:"website,omitempty"`
}

// PropertySearchResult representa una property encontrada en la búsqueda
type PropertySearchResult struct {
	ID           int64   `json:"id"`
	PropertyCode *string `json:"propertyCode,omitempty"`
	Street       *string `json:"street,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	Zip          *string `json:"zip,omitempty"`
	CustomerName *string `json:"customerName,omitempty"`
}

// InvoiceSearchResult representa una invoice encontrada en la búsqueda
type InvoiceSearchResult struct {
	ID            int64   `json:"id"`
	InvoiceNumber *string `json:"invoiceNumber,omitempty"`
	Total         float64 `json:"total"`
	WorkOrder     *string `json:"workOrder,omitempty"`
	CustomerName  *string `json:"customerName,omitempty"`
}

// QuoteSearchResult representa una quote encontrada en la búsqueda
type QuoteSearchResult struct {
	ID           int64   `json:"id"`
	QuoteNumber  *string `json:"quoteNumber,omitempty"`
	Total        float64 `json:"total"`
	WorkOrder    *string `json:"workOrder,omitempty"`
	CustomerName *string `json:"customerName,omitempty"`
	StatusName   *string `json:"statusName,omitempty"`
}

// WarrantySearchResult representa una warranty encontrada en la búsqueda
type WarrantySearchResult struct {
	ID              int64   `json:"id"`
	WarrantyNumber  *string `json:"warrantyNumber,omitempty"`
	AgreementNumber *string `json:"agreementNumber,omitempty"`
	WorkOrder       *string `json:"workOrder,omitempty"`
	CustomerName    *string `json:"customerName,omitempty"`
	TypeName        *string `json:"typeName,omitempty"`
	StatusName      *string `json:"statusName,omitempty"`
}

// WarrantyClaimSearchResult representa un warranty claim encontrado en la búsqueda
type WarrantyClaimSearchResult struct {
	ID                  int64   `json:"id"`
	ClaimNumber         *string `json:"claimNumber,omitempty"`
	InternalClaimNumber *string `json:"internalClaimNumber,omitempty"`
	WarrantyPart        *string `json:"warrantyPart,omitempty"`
	Manufacturer        *string `json:"manufacturer,omitempty"`
	WorkOrder           *string `json:"workOrder,omitempty"`
	StatusName          *string `json:"statusName,omitempty"`
}

// UserSearchResult representa un user encontrado en la búsqueda
type UserSearchResult struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	RoleName *string `json:"roleName,omitempty"`
	IsActive bool    `json:"isActive"`
}

// GlobalSearchResponse respuesta de la búsqueda global
type GlobalSearchResponse struct {
	Query          string                      `json:"query"`
	Jobs           []JobSearchResult           `json:"jobs"`
	Customers      []CustomerSearchResult      `json:"customers"`
	Properties     []PropertySearchResult      `json:"properties"`
	Invoices       []InvoiceSearchResult       `json:"invoices"`
	Quotes         []QuoteSearchResult         `json:"quotes"`
	Warranties     []WarrantySearchResult      `json:"warranties"`
	WarrantyClaims []WarrantyClaimSearchResult `json:"warrantyClaims"`
	Users          []UserSearchResult          `json:"users"`
	TotalResults   int                         `json:"totalResults"`
}

// SearchFilters filtros para la búsqueda global
type SearchFilters struct {
	Query    string
	Limit    int
	UserID   *int64 // Para filtrar por usuario (job_view_user_only)
	UserOnly bool   // Si el usuario tiene permiso job_view_user_only
}
