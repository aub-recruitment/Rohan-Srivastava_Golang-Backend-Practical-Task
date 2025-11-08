package usecases

import (
	"context"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"
	"github.com/google/uuid"
)

type WatchHistoryUseCase struct {
	watchHistoryRepo repositories.WatchHistoryRepository
	contentRepo      repositories.ContentRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewWatchHistoryUseCase(watchHistoryRepo repositories.WatchHistoryRepository, contentRepo repositories.ContentRepository, subscriptionRepo repositories.SubscriptionRepository) *WatchHistoryUseCase {
	return &WatchHistoryUseCase{
		watchHistoryRepo: watchHistoryRepo,
		contentRepo:      contentRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

type WatchHistoryInput struct {
	ContentID      uuid.UUID `json:"content_id" binding:"required"`
	WatchedSeconds *int      `json:"watched_seconds" binding:"required,gte=0"`
}

func (uc *WatchHistoryUseCase) CreateOrUpdateWatchHistory(ctx context.Context, userID uuid.UUID, input WatchHistoryInput) (*domain.WatchHistory, error) {
	content, err := uc.contentRepo.GetByID(ctx, input.ContentID)
	if err != nil {
		return nil, err
	}
	if content.AccessLevel != domain.AccessLevelFree {
		subscription, err := uc.subscriptionRepo.GetActiveByUserID(ctx, userID)
		if err != nil || subscription.IsExpired() {
			return nil, domain.ErrContentNotAccessible
		}
	}
	existingHistory, err := uc.watchHistoryRepo.GetByUserAndContent(ctx, userID, input.ContentID)
	if err != nil && err != domain.ErrWatchHistoryNotFound {
		return nil, err
	}
	if existingHistory != nil {
		existingHistory.WatchedSeconds = *input.WatchedSeconds
		existingHistory.LastWatchedAt = time.Now()
		if existingHistory.ProgressPercentage() >= 90.0 {
			existingHistory.Status = domain.WatchStatusCompleted
		} else if *input.WatchedSeconds > 0 {
			existingHistory.Status = domain.WatchStatusPaused
		}
		if err := uc.watchHistoryRepo.Update(ctx, existingHistory); err != nil {
			return nil, err
		}
		existingHistory.Content = content
		return existingHistory, nil
	}
	watchHistory := &domain.WatchHistory{
		ID:             uuid.New(),
		UserID:         userID,
		ContentID:      input.ContentID,
		WatchedSeconds: *input.WatchedSeconds,
		TotalSeconds:   content.DurationSeconds,
		Status:         domain.WatchStatusStarted,
		LastWatchedAt:  time.Now(),
	}
	if watchHistory.ProgressPercentage() >= 90.0 {
		watchHistory.Status = domain.WatchStatusCompleted
	}
	if err := uc.watchHistoryRepo.Create(ctx, watchHistory); err != nil {
		return nil, err
	}
	watchHistory.Content = content
	return watchHistory, nil
}

func (uc *WatchHistoryUseCase) GetWatchHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.WatchHistory, int64, error) {
	return uc.watchHistoryRepo.GetByUserID(ctx, userID, limit, offset)
}

func (uc *WatchHistoryUseCase) GetContinueWatching(ctx context.Context, userID uuid.UUID) ([]*domain.WatchHistory, error) {
	return uc.watchHistoryRepo.GetContinueWatching(ctx, userID, 10)
}

func (uc *WatchHistoryUseCase) UpdateProgress(ctx context.Context, userID, historyID uuid.UUID, watchedSeconds int) (*domain.WatchHistory, error) {
	history, err := uc.watchHistoryRepo.GetByID(ctx, historyID)
	if err != nil {
		return nil, err
	}
	if history.UserID != userID {
		return nil, domain.ErrForbidden
	}
	history.WatchedSeconds = watchedSeconds
	history.LastWatchedAt = time.Now()
	if history.ProgressPercentage() >= 90.0 {
		history.Status = domain.WatchStatusCompleted
	} else {
		history.Status = domain.WatchStatusPaused
	}
	if err := uc.watchHistoryRepo.Update(ctx, history); err != nil {
		return nil, err
	}
	return history, nil
}
