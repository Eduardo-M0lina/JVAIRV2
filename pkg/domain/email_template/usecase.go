package email_template

import "context"

type Service interface {
	Create(ctx context.Context, template *EmailTemplate) error
	GetByID(ctx context.Context, id int64) (*EmailTemplate, error)
	GetByLabel(ctx context.Context, label string) (*EmailTemplate, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*EmailTemplate, int, error)
	Update(ctx context.Context, template *EmailTemplate) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}
