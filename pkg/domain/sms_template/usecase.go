package sms_template

import "context"

type Service interface {
	Create(ctx context.Context, template *SMSTemplate) error
	GetByID(ctx context.Context, id int64) (*SMSTemplate, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*SMSTemplate, int, error)
	Update(ctx context.Context, template *SMSTemplate) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}
