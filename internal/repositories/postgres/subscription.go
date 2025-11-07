package postgres

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct{ db *gorm.DB }

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var subscription domain.Subscription
	err := r.db.WithContext(ctx).Preload("User").Preload("Plan").Where("id = ?", id).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	var subscription domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Plan").
		Where("user_id = ? AND is_active = ? AND status = ?",
			userID, true, domain.SubscriptionStatusActive).
		Order("end_date DESC").
		First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) GetHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	var subscriptions []*domain.Subscription
	err := r.db.WithContext(ctx).Preload("Plan").Where("user_id = ?", userID).Order("created_at DESC").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *SubscriptionRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Subscription{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": false,
		"status":    domain.SubscriptionStatusCancelled,
	}).Error
}
