package sms_template

import "context"

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*SMSTemplate, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	return uc.repo.List(ctx, filters, page, pageSize)
}
