package search

import (
	"context"
	"sync"
)

// Service define las operaciones del caso de uso de búsqueda global
type Service interface {
	GlobalSearch(ctx context.Context, filters SearchFilters) (*GlobalSearchResponse, error)
}

// UseCase implementa la lógica de negocio de búsqueda global
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// GlobalSearch realiza una búsqueda en todas las entidades en paralelo
func (uc *UseCase) GlobalSearch(ctx context.Context, filters SearchFilters) (*GlobalSearchResponse, error) {
	if filters.Limit <= 0 {
		filters.Limit = 10
	}

	response := &GlobalSearchResponse{
		Query:          filters.Query,
		Jobs:           []JobSearchResult{},
		Customers:      []CustomerSearchResult{},
		Properties:     []PropertySearchResult{},
		Invoices:       []InvoiceSearchResult{},
		Quotes:         []QuoteSearchResult{},
		Warranties:     []WarrantySearchResult{},
		WarrantyClaims: []WarrantyClaimSearchResult{},
		Users:          []UserSearchResult{},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 8)

	// Search Jobs
	wg.Add(1)
	go func() {
		defer wg.Done()
		jobs, err := uc.repo.SearchJobs(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Jobs = jobs
		mu.Unlock()
	}()

	// Search Customers
	wg.Add(1)
	go func() {
		defer wg.Done()
		customers, err := uc.repo.SearchCustomers(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Customers = customers
		mu.Unlock()
	}()

	// Search Properties
	wg.Add(1)
	go func() {
		defer wg.Done()
		properties, err := uc.repo.SearchProperties(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Properties = properties
		mu.Unlock()
	}()

	// Search Invoices
	wg.Add(1)
	go func() {
		defer wg.Done()
		invoices, err := uc.repo.SearchInvoices(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Invoices = invoices
		mu.Unlock()
	}()

	// Search Quotes
	wg.Add(1)
	go func() {
		defer wg.Done()
		quotes, err := uc.repo.SearchQuotes(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Quotes = quotes
		mu.Unlock()
	}()

	// Search Warranties
	wg.Add(1)
	go func() {
		defer wg.Done()
		warranties, err := uc.repo.SearchWarranties(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.Warranties = warranties
		mu.Unlock()
	}()

	// Search Warranty Claims
	wg.Add(1)
	go func() {
		defer wg.Done()
		claims, err := uc.repo.SearchWarrantyClaims(ctx, filters)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		response.WarrantyClaims = claims
		mu.Unlock()
	}()

	// Search Users (solo si no tiene restricción job_view_user_only)
	if !filters.UserOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			users, err := uc.repo.SearchUsers(ctx, filters)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			response.Users = users
			mu.Unlock()
		}()
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()
	close(errChan)

	// Verificar si hubo errores
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Calcular total de resultados
	response.TotalResults = len(response.Jobs) + len(response.Customers) +
		len(response.Properties) + len(response.Invoices) +
		len(response.Quotes) + len(response.Warranties) +
		len(response.WarrantyClaims) + len(response.Users)

	return response, nil
}
