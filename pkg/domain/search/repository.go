package search

import "context"

// Repository define las operaciones de búsqueda global
type Repository interface {
	SearchJobs(ctx context.Context, filters SearchFilters) ([]JobSearchResult, error)
	SearchCustomers(ctx context.Context, filters SearchFilters) ([]CustomerSearchResult, error)
	SearchProperties(ctx context.Context, filters SearchFilters) ([]PropertySearchResult, error)
	SearchInvoices(ctx context.Context, filters SearchFilters) ([]InvoiceSearchResult, error)
	SearchQuotes(ctx context.Context, filters SearchFilters) ([]QuoteSearchResult, error)
	SearchWarranties(ctx context.Context, filters SearchFilters) ([]WarrantySearchResult, error)
	SearchWarrantyClaims(ctx context.Context, filters SearchFilters) ([]WarrantyClaimSearchResult, error)
	SearchUsers(ctx context.Context, filters SearchFilters) ([]UserSearchResult, error)
}
