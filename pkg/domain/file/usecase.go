package file

import (
	"context"
	"io"
)

// StorageService define las operaciones de almacenamiento de archivos
type StorageService interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetObject(ctx context.Context, key string) (io.ReadCloser, string, error)
}

// Service define las operaciones de negocio para archivos
type Service interface {
	Upload(ctx context.Context, fileableID int64, fileableType string, filename string, contentType string, body io.Reader) (*File, error)
	GetByID(ctx context.Context, id int64) (*File, error)
	ListByFileable(ctx context.Context, fileableID int64, fileableType string) ([]*File, error)
	Delete(ctx context.Context, id int64) error
	Download(ctx context.Context, id int64) (io.ReadCloser, string, string, error)
}

// UseCase implementa la lógica de negocio para archivos
type UseCase struct {
	repo    Repository
	storage StorageService
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository, storage StorageService) *UseCase {
	return &UseCase{
		repo:    repo,
		storage: storage,
	}
}
