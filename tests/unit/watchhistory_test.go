package unit

import (
	"context"
	"testing"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWatchHistoryRepository struct {
	mock.Mock
}

func (m *MockWatchHistoryRepository) Create(ctx context.Context, history *domain.WatchHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockWatchHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WatchHistory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WatchHistory), args.Error(1)
}

func (m *MockWatchHistoryRepository) GetByUserAndContent(ctx context.Context, userID, contentID uuid.UUID) (*domain.WatchHistory, error) {
	args := m.Called(ctx, userID, contentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WatchHistory), args.Error(1)
}

func (m *MockWatchHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.WatchHistory, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return args.Get(0).([]*domain.WatchHistory), total, args.Error(2)
}

func (m *MockWatchHistoryRepository) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.WatchHistory, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WatchHistory), args.Error(1)
}

func (m *MockWatchHistoryRepository) Update(ctx context.Context, history *domain.WatchHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockWatchHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateWatchHistory_FreeContent(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	contentID := uuid.New()

	content := &domain.Content{
		ID:              contentID,
		Title:           "Free Movie",
		AccessLevel:     domain.AccessLevelFree,
		DurationSeconds: 7200,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)
	mockWatchRepo.On("GetByUserAndContent", mock.Anything, userID, contentID).Return(nil, domain.ErrWatchHistoryNotFound)
	mockWatchRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.WatchHistory")).Return(nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)

	watchedSeconds := 1800
	input := usecases.WatchHistoryInput{
		ContentID:      contentID,
		WatchedSeconds: &watchedSeconds,
	}

	result, err := watchUseCase.CreateOrUpdateWatchHistory(context.Background(), userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, contentID, result.ContentID)
	assert.Equal(t, 1800, result.WatchedSeconds)
	mockContentRepo.AssertExpectations(t)
	mockWatchRepo.AssertExpectations(t)
}

func TestCreateWatchHistory_PremiumContent_NoSubscription(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	contentID := uuid.New()

	content := &domain.Content{
		ID:              contentID,
		AccessLevel:     domain.AccessLevelPremium,
		DurationSeconds: 7200,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)
	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(nil, domain.ErrSubscriptionNotFound)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)

	watchedSeconds := 1800
	input := usecases.WatchHistoryInput{
		ContentID:      contentID,
		WatchedSeconds: &watchedSeconds,
	}

	result, err := watchUseCase.CreateOrUpdateWatchHistory(context.Background(), userID, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrContentNotAccessible, err)
	mockContentRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
}

func TestUpdateWatchHistory_ExistingEntry(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	contentID := uuid.New()

	content := &domain.Content{
		ID:              contentID,
		AccessLevel:     domain.AccessLevelFree,
		DurationSeconds: 7200,
	}

	existingHistory := &domain.WatchHistory{
		ID:             uuid.New(),
		UserID:         userID,
		ContentID:      contentID,
		WatchedSeconds: 1800,
		TotalSeconds:   7200,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)
	mockWatchRepo.On("GetByUserAndContent", mock.Anything, userID, contentID).Return(existingHistory, nil)
	mockWatchRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.WatchHistory")).Return(nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)

	watchedSeconds := 3600
	input := usecases.WatchHistoryInput{
		ContentID:      contentID,
		WatchedSeconds: &watchedSeconds,
	}

	result, err := watchUseCase.CreateOrUpdateWatchHistory(context.Background(), userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3600, result.WatchedSeconds)
	mockContentRepo.AssertExpectations(t)
	mockWatchRepo.AssertExpectations(t)
}

func TestGetWatchHistory_Success(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	histories := []*domain.WatchHistory{
		{ID: uuid.New(), UserID: userID, ContentID: uuid.New()},
		{ID: uuid.New(), UserID: userID, ContentID: uuid.New()},
	}

	mockWatchRepo.On("GetByUserID", mock.Anything, userID, 20, 0).Return(histories, int64(2), nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)
	result, total, err := watchUseCase.GetWatchHistory(context.Background(), userID, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockWatchRepo.AssertExpectations(t)
}

func TestGetContinueWatching_Success(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	histories := []*domain.WatchHistory{
		{
			ID:             uuid.New(),
			UserID:         userID,
			WatchedSeconds: 1800,
			TotalSeconds:   7200,
			Status:         domain.WatchStatusPaused,
		},
	}

	mockWatchRepo.On("GetContinueWatching", mock.Anything, userID, 10).Return(histories, nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)
	result, err := watchUseCase.GetContinueWatching(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockWatchRepo.AssertExpectations(t)
}

func TestUpdateProgress_Success(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	historyID := uuid.New()

	history := &domain.WatchHistory{
		ID:             historyID,
		UserID:         userID,
		WatchedSeconds: 1800,
		TotalSeconds:   7200,
		Status:         domain.WatchStatusPaused,
	}

	mockWatchRepo.On("GetByID", mock.Anything, historyID).Return(history, nil)
	mockWatchRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.WatchHistory")).Return(nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)
	result, err := watchUseCase.UpdateProgress(context.Background(), userID, historyID, 6500)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 6500, result.WatchedSeconds)
	assert.Equal(t, domain.WatchStatusCompleted, result.Status) // 90% watched = completed
	mockWatchRepo.AssertExpectations(t)
}

func TestUpdateProgress_NotOwner(t *testing.T) {
	mockWatchRepo := new(MockWatchHistoryRepository)
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	otherUserID := uuid.New()
	historyID := uuid.New()

	history := &domain.WatchHistory{
		ID:     historyID,
		UserID: otherUserID,
	}

	mockWatchRepo.On("GetByID", mock.Anything, historyID).Return(history, nil)

	watchUseCase := usecases.NewWatchHistoryUseCase(mockWatchRepo, mockContentRepo, mockSubRepo)
	result, err := watchUseCase.UpdateProgress(context.Background(), userID, historyID, 3600)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrForbidden, err)
	mockWatchRepo.AssertExpectations(t)
}

func TestProgressPercentage_Calculation(t *testing.T) {
	history := &domain.WatchHistory{
		WatchedSeconds: 3600,
		TotalSeconds:   7200,
	}

	percentage := history.ProgressPercentage()
	assert.Equal(t, 50.0, percentage)
}

func TestIsCompleted_AtThreshold(t *testing.T) {
	history := &domain.WatchHistory{
		WatchedSeconds: 6480,
		TotalSeconds:   7200,
		Status:         domain.WatchStatusPaused,
	}

	completed := history.IsCompleted()
	assert.True(t, completed)
}

func TestIsCompleted_NotReached(t *testing.T) {
	history := &domain.WatchHistory{
		WatchedSeconds: 3600,
		TotalSeconds:   7200,
		Status:         domain.WatchStatusPaused,
	}

	completed := history.IsCompleted()
	assert.False(t, completed)
}
