package usecases

import (
	"context"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"
	"github.com/google/uuid"
)

type SubscriptionUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
	planRepo         repositories.PlanRepository
	userRepo         repositories.UserRepository
}

func NewSubscriptionUseCase(subscriptionRepo repositories.SubscriptionRepository, planRepo repositories.PlanRepository, userRepo repositories.UserRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		userRepo:         userRepo,
	}
}

type CreateSubscriptionInput struct {
	PlanID uuid.UUID `json:"plan_id" binding:"required"`
}

func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, userID uuid.UUID, input CreateSubscriptionInput) (*domain.Subscription, error) {
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	activeSubscription, _ := uc.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if activeSubscription != nil && !activeSubscription.IsExpired() {
		return nil, domain.ErrActiveSubscriptionExists
	}
	plan, err := uc.planRepo.GetByID(ctx, input.PlanID)
	if err != nil {
		return nil, err
	}
	if !plan.IsActive {
		return nil, domain.ErrPlanNotAvailable
	}
	now := time.Now()
	subscription := &domain.Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		PlanID:    plan.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.ValidityDays),
		IsActive:  true,
		Status:    domain.SubscriptionStatusActive,
	}
	if err := uc.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}
	subscription.Plan = plan
	return subscription, nil
}

func (uc *SubscriptionUseCase) GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	subscription, err := uc.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if subscription.IsExpired() {
		subscription.IsActive = false
		subscription.Status = domain.SubscriptionStatusExpired
		_ = uc.subscriptionRepo.Update(ctx, subscription)
		return nil, domain.ErrSubscriptionExpired
	}
	return subscription, nil
}

func (uc *SubscriptionUseCase) GetSubscriptionHistory(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	return uc.subscriptionRepo.GetHistoryByUserID(ctx, userID)
}

func (uc *SubscriptionUseCase) CancelSubscription(ctx context.Context, userID, subscriptionID uuid.UUID) error {
	subscription, err := uc.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if subscription.UserID != userID {
		return domain.ErrForbidden
	}
	if !subscription.IsActive {
		return domain.ErrSubscriptionInactive
	}
	return uc.subscriptionRepo.Cancel(ctx, subscriptionID)
}

func (uc *SubscriptionUseCase) RenewSubscription(ctx context.Context, userID, subscriptionID uuid.UUID) (*domain.Subscription, error) {
	oldSubscription, err := uc.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	if oldSubscription.UserID != userID {
		return nil, domain.ErrForbidden
	}
	plan, err := uc.planRepo.GetByID(ctx, oldSubscription.PlanID)
	if err != nil {
		return nil, err
	}
	if !plan.IsActive {
		return nil, domain.ErrPlanNotAvailable
	}
	now := time.Now()
	newSubscription := &domain.Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		PlanID:    plan.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.ValidityDays),
		IsActive:  true,
		Status:    domain.SubscriptionStatusActive,
	}
	if err := uc.subscriptionRepo.Create(ctx, newSubscription); err != nil {
		return nil, err
	}
	_ = uc.subscriptionRepo.Cancel(ctx, subscriptionID)
	newSubscription.Plan = plan
	return newSubscription, nil
}
