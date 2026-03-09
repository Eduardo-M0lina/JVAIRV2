package alert

import "context"

type ListFilters struct {
	UserID     *int64
	IsRead     *bool
	AlertType  *string
	EntityType *string
}

type Repository interface {
	Create(ctx context.Context, alert *Alert) error
	GetByID(ctx context.Context, id int64) (*Alert, error)
	List(ctx context.Context, filters ListFilters, limit, offset int) ([]*Alert, int64, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllRead(ctx context.Context, userID int64) (int64, error)
	UnreadCount(ctx context.Context, userID int64) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type UserExistsChecker interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}
