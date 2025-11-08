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

func TestCreateSubscription_Success(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	planID := uuid.New()

	user := &domain.User{ID: userID, Email: "test@example.com"}
	plan := &domain.Plan{
		ID:           planID,
		Name:         "Premium",
		Price:        999,
		ValidityDays: 30,
		AccessLevel:  domain.AccessLevelPremium,
		IsActive:     true,
	}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(nil, domain.ErrSubscriptionNotFound)
	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(plan, nil)
	mockSubRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Subscription")).Return(nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)

	input := usecases.CreateSubscriptionInput{PlanID: planID}
	subscription, err := subUseCase.CreateSubscription(context.Background(), userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, userID, subscription.UserID)
	assert.Equal(t, planID, subscription.PlanID)
	assert.Equal(t, domain.SubscriptionStatusActive, subscription.Status)
	mockUserRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
	mockPlanRepo.AssertExpectations(t)
}

func TestCreateSubscription_UserNotFound(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	planID := uuid.New()

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrUserNotFound)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)

	input := usecases.CreateSubscriptionInput{PlanID: planID}
	subscription, err := subUseCase.CreateSubscription(context.Background(), userID, input)

	assert.Error(t, err)
	assert.Nil(t, subscription)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateSubscription_ActiveSubscriptionExists(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	planID := uuid.New()

	user := &domain.User{ID: userID}
	activeSubscription := &domain.Subscription{
		ID:       uuid.New(),
		UserID:   userID,
		IsActive: true,
		EndDate:  time.Now().Add(24 * time.Hour),
	}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(activeSubscription, nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)

	input := usecases.CreateSubscriptionInput{PlanID: planID}
	subscription, err := subUseCase.CreateSubscription(context.Background(), userID, input)

	assert.Error(t, err)
	assert.Nil(t, subscription)
	assert.Equal(t, domain.ErrActiveSubscriptionExists, err)
	mockUserRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
}

func TestCreateSubscription_PlanNotActive(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	planID := uuid.New()

	user := &domain.User{ID: userID}
	plan := &domain.Plan{ID: planID, IsActive: false}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(nil, domain.ErrSubscriptionNotFound)
	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(plan, nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)

	input := usecases.CreateSubscriptionInput{PlanID: planID}
	subscription, err := subUseCase.CreateSubscription(context.Background(), userID, input)

	assert.Error(t, err)
	assert.Nil(t, subscription)
	assert.Equal(t, domain.ErrPlanNotAvailable, err)
	mockUserRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
	mockPlanRepo.AssertExpectations(t)
}

func TestGetActiveSubscription_Success(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	subscription := &domain.Subscription{
		ID:       uuid.New(),
		UserID:   userID,
		IsActive: true,
		EndDate:  time.Now().Add(24 * time.Hour),
	}

	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(subscription, nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	result, err := subUseCase.GetActiveSubscription(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, subscription, result)
	mockSubRepo.AssertExpectations(t)
}

func TestGetActiveSubscription_Expired(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	subscription := &domain.Subscription{
		ID:       uuid.New(),
		UserID:   userID,
		IsActive: true,
		EndDate:  time.Now().Add(-24 * time.Hour), // expired
	}

	mockSubRepo.On("GetActiveByUserID", mock.Anything, userID).Return(subscription, nil)
	mockSubRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Subscription")).Return(nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	result, err := subUseCase.GetActiveSubscription(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrSubscriptionExpired, err)
	mockSubRepo.AssertExpectations(t)
}

func TestGetSubscriptionHistory_Success(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	subscriptions := []*domain.Subscription{
		{ID: uuid.New(), UserID: userID},
		{ID: uuid.New(), UserID: userID},
	}

	mockSubRepo.On("GetHistoryByUserID", mock.Anything, userID).Return(subscriptions, nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	result, err := subUseCase.GetSubscriptionHistory(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockSubRepo.AssertExpectations(t)
}

func TestCancelSubscription_Success(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	subID := uuid.New()
	subscription := &domain.Subscription{
		ID:       subID,
		UserID:   userID,
		IsActive: true,
	}

	mockSubRepo.On("GetByID", mock.Anything, subID).Return(subscription, nil)
	mockSubRepo.On("Cancel", mock.Anything, subID).Return(nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	err := subUseCase.CancelSubscription(context.Background(), userID, subID)

	assert.NoError(t, err)
	mockSubRepo.AssertExpectations(t)
}

func TestCancelSubscription_NotOwner(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	otherUserID := uuid.New()
	subID := uuid.New()
	subscription := &domain.Subscription{
		ID:       subID,
		UserID:   otherUserID,
		IsActive: true,
	}

	mockSubRepo.On("GetByID", mock.Anything, subID).Return(subscription, nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	err := subUseCase.CancelSubscription(context.Background(), userID, subID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)
	mockSubRepo.AssertExpectations(t)
}

func TestRenewSubscription_Success(t *testing.T) {
	mockSubRepo := new(MockSubscriptionRepository)
	mockPlanRepo := new(MockPlanRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New()
	oldSubID := uuid.New()
	planID := uuid.New()

	oldSubscription := &domain.Subscription{
		ID:     oldSubID,
		UserID: userID,
		PlanID: planID,
	}

	plan := &domain.Plan{
		ID:           planID,
		ValidityDays: 30,
		IsActive:     true,
	}

	mockSubRepo.On("GetByID", mock.Anything, oldSubID).Return(oldSubscription, nil)
	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(plan, nil)
	mockSubRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Subscription")).Return(nil)
	mockSubRepo.On("Cancel", mock.Anything, oldSubID).Return(nil)

	subUseCase := usecases.NewSubscriptionUseCase(mockSubRepo, mockPlanRepo, mockUserRepo)
	result, err := subUseCase.RenewSubscription(context.Background(), userID, oldSubID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	mockSubRepo.AssertExpectations(t)
	mockPlanRepo.AssertExpectations(t)
}
