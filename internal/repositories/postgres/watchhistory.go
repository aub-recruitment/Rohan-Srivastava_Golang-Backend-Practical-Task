package postgres

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WatchHistoryRepository struct{ db *gorm.DB }

func NewWatchHistoryRepository(db *gorm.DB) *WatchHistoryRepository {
	return &WatchHistoryRepository{db: db}
}

func (r *WatchHistoryRepository) Create(ctx context.Context, history *domain.WatchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *WatchHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WatchHistory, error) {
	var history domain.WatchHistory
	err := r.db.WithContext(ctx).Preload("Content").Where("id = ?", id).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWatchHistoryNotFound
		}
		return nil, err
	}
	return &history, nil
}

func (r *WatchHistoryRepository) GetByUserAndContent(ctx context.Context, userID, contentID uuid.UUID) (*domain.WatchHistory, error) {
	var history domain.WatchHistory
	err := r.db.WithContext(ctx).Where("user_id = ? AND content_id = ?", userID, contentID).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWatchHistoryNotFound
		}
		return nil, err
	}
	return &history, nil
}

func (r *WatchHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.WatchHistory, int64, error) {
	var histories []*domain.WatchHistory
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.WatchHistory{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Content").Limit(limit).Offset(offset).Order("last_watched_at DESC").Find(&histories).Error
	return histories, total, err
}

func (r *WatchHistoryRepository) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.WatchHistory, error) {
	var histories []*domain.WatchHistory
	err := r.db.WithContext(ctx).Preload("Content").
		Where("user_id = ? AND status != ?", userID, domain.WatchStatusCompleted).
		Where("watched_seconds > 0").
		Order("last_watched_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *WatchHistoryRepository) Update(ctx context.Context, history *domain.WatchHistory) error {
	return r.db.WithContext(ctx).Save(history).Error
}

func (r *WatchHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WatchHistory{}, "id = ?", id).Error
}
