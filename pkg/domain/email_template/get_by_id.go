package email_template

import "context"

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*EmailTemplate, error) {
	return uc.repo.GetByID(ctx, id)
}
