package unit

import (
	"context"
	"testing"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockContentRepository struct {
	mock.Mock
}

func (m *MockContentRepository) Create(ctx context.Context, content *domain.Content) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func (m *MockContentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Content, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Content), args.Error(1)
}

func (m *MockContentRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Content, int64, error) {
	args := m.Called(ctx, filters, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return args.Get(0).([]*domain.Content), total, args.Error(2)
}

func (m *MockContentRepository) Update(ctx context.Context, content *domain.Content) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func (m *MockContentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateContent_Success(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	mockContentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Content")).Return(nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)

	input := usecases.CreateContentInput{
		Title:           "Test Movie",
		Description:     "A test movie",
		AccessLevel:     domain.AccessLevelPremium,
		DurationSeconds: 7200,
		Published:       true,
	}

	content, err := contentUseCase.CreateContent(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.Equal(t, "Test Movie", content.Title)
	assert.Equal(t, domain.AccessLevelPremium, content.AccessLevel)
	mockContentRepo.AssertExpectations(t)
}

func TestGetContent_Published_FreeAccess_NoAuth(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contentID := uuid.New()
	content := &domain.Content{
		ID:          contentID,
		Title:       "Free Content",
		AccessLevel: domain.AccessLevelFree,
		Published:   true,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	result, err := contentUseCase.GetContent(context.Background(), contentID, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Free Content", result.Title)
	mockContentRepo.AssertExpectations(t)
}

func TestGetContent_Published_PremiumAccess_NoAuth(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contentID := uuid.New()
	content := &domain.Content{
		ID:          contentID,
		Title:       "Premium Content",
		AccessLevel: domain.AccessLevelPremium,
		Published:   true,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	result, err := contentUseCase.GetContent(context.Background(), contentID, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrContentNotAccessible, err)
	mockContentRepo.AssertExpectations(t)
}

func TestGetContent_NotPublished(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contentID := uuid.New()
	content := &domain.Content{
		ID:          contentID,
		Title:       "Draft Content",
		AccessLevel: domain.AccessLevelFree,
		Published:   false,
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	result, err := contentUseCase.GetContent(context.Background(), contentID, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrContentNotPublished, err)
	mockContentRepo.AssertExpectations(t)
}

func TestListContent_Success(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contents := []*domain.Content{
		{ID: uuid.New(), Title: "Movie 1", Published: true},
		{ID: uuid.New(), Title: "Movie 2", Published: true},
	}

	mockContentRepo.On("List", mock.Anything, mock.MatchedBy(func(filters map[string]interface{}) bool {
		val, exists := filters["published"]
		return exists && val == true
	}), 20, 0).Return(contents, int64(2), nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	result, total, err := contentUseCase.ListContent(context.Background(), true, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockContentRepo.AssertExpectations(t)
}

func TestUpdateContent_Success(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contentID := uuid.New()
	content := &domain.Content{
		ID:    contentID,
		Title: "Original Title",
	}

	mockContentRepo.On("GetByID", mock.Anything, contentID).Return(content, nil)
	mockContentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Content")).Return(nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)

	input := usecases.CreateContentInput{
		Title:       "Updated Title",
		AccessLevel: domain.AccessLevelBasic,
	}

	result, err := contentUseCase.UpdateContent(context.Background(), contentID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Title", result.Title)
	mockContentRepo.AssertExpectations(t)
}

func TestDeleteContent_Success(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	contentID := uuid.New()
	mockContentRepo.On("Delete", mock.Anything, contentID).Return(nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	err := contentUseCase.DeleteContent(context.Background(), contentID)

	assert.NoError(t, err)
	mockContentRepo.AssertExpectations(t)
}

func TestCheckAccess_FreeContent(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)

	hasAccess, err := contentUseCase.CheckAccess(context.Background(), userID, domain.AccessLevelFree)

	assert.NoError(t, err)
	assert.True(t, hasAccess)
}

func TestCheckAccess_PremiumContent_NoSubscription(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(nil, domain.ErrSubscriptionNotFound)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	hasAccess, err := contentUseCase.CheckAccess(context.Background(), userID, domain.AccessLevelPremium)

	assert.NoError(t, err)
	assert.False(t, hasAccess)
	mockSubRepo.AssertExpectations(t)
}

func TestCheckAccess_PremiumContent_WithValidSubscription(t *testing.T) {
	mockContentRepo := new(MockContentRepository)
	mockSubRepo := new(MockSubscriptionRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	subscription := &domain.Subscription{
		ID:       uuid.New(),
		UserID:   userID,
		IsActive: true,
		EndDate:  time.Now().Add(24 * time.Hour),
		Plan:     &domain.Plan{AccessLevel: domain.AccessLevelPremium},
	}

	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(subscription, nil)

	contentUseCase := usecases.NewContentUseCase(mockContentRepo, mockSubRepo, mockUserRepo)
	hasAccess, err := contentUseCase.CheckAccess(context.Background(), userID, domain.AccessLevelPremium)

	assert.NoError(t, err)
	assert.True(t, hasAccess)
	mockSubRepo.AssertExpectations(t)
}
