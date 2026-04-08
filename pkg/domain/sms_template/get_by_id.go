package sms_template

import "context"

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*SMSTemplate, error) {
	return uc.repo.GetByID(ctx, id)
}
