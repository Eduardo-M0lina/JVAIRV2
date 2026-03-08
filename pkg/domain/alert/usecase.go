package alert

import "context"

type Service interface {
	Create(ctx context.Context, alert *Alert) error
	GetByID(ctx context.Context, id int64) (*Alert, error)
	List(ctx context.Context, filters ListFilters, limit, offset int) ([]*Alert, int64, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllRead(ctx context.Context, userID int64) (int64, error)
	UnreadCount(ctx context.Context, userID int64) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo        Repository
	userChecker UserExistsChecker
}

func NewUseCase(repo Repository, userChecker UserExistsChecker) Service {
	return &service{
		repo:        repo,
		userChecker: userChecker,
	}
}
