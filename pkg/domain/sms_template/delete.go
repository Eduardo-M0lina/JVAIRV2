package sms_template

import "context"

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
