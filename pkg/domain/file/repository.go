package file

import "context"

// Repository define las operaciones de persistencia para archivos
type Repository interface {
	Create(ctx context.Context, file *File) (int64, error)
	GetByID(ctx context.Context, id int64) (*File, error)
	ListByFileable(ctx context.Context, fileableID int64, fileableType string) ([]*File, error)
	Delete(ctx context.Context, id int64) error
}
