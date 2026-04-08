package email_template

import "context"

func (uc *UseCase) GetByLabel(ctx context.Context, label string) (*EmailTemplate, error) {
	return uc.repo.GetByLabel(ctx, label)
}
