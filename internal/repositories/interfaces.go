package repositories

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ContentRepository interface {
	Create(ctx context.Context, content *domain.Content) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Content, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Content, int64, error)
	Update(ctx context.Context, content *domain.Content) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error)
	GetByName(ctx context.Context, name string) (*domain.Plan, error)
	List(ctx context.Context, activeOnly bool) ([]*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error)
	GetHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error)
	Update(ctx context.Context, subscription *domain.Subscription) error
	Cancel(ctx context.Context, id uuid.UUID) error
}

type WatchHistoryRepository interface {
	Create(ctx context.Context, history *domain.WatchHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WatchHistory, error)
	GetByUserAndContent(ctx context.Context, userID, contentID uuid.UUID) (*domain.WatchHistory, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.WatchHistory, int64, error)
	GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.WatchHistory, error)
	Update(ctx context.Context, history *domain.WatchHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
}
