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

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	expectedUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockSubRepo)
	user, err := userUseCase.GetProfile(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestGetProfile_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrUserNotFound)

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockSubRepo)
	user, err := userUseCase.GetProfile(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	originalUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Original Name",
		Bio:   "Old Bio",
	}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(originalUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockSubRepo)

	input := usecases.UpdateProfileInput{
		Name: "Updated Name",
		Bio:  "New Bio",
	}

	user, err := userUseCase.UpdateProfile(context.Background(), userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Updated Name", user.Name)
	assert.Equal(t, "New Bio", user.Bio)
	mockUserRepo.AssertExpectations(t)
}

func TestGetSubscriptionHistory_UserUseCase_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSubRepo := new(MockSubscriptionRepository)

	userID := uuid.New()
	subscriptions := []*domain.Subscription{
		{ID: uuid.New(), UserID: userID},
		{ID: uuid.New(), UserID: userID},
	}

	mockSubRepo.On("GetHistoryByUserID", mock.Anything, userID).Return(subscriptions, nil)

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockSubRepo)
	result, err := userUseCase.GetSubscriptionHistory(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockSubRepo.AssertExpectations(t)
}
