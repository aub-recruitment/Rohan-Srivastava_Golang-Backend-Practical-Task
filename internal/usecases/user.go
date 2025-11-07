package usecases

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"

	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo         repositories.UserRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, subscriptionRepo repositories.SubscriptionRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, subscriptionRepo: subscriptionRepo}
}

type UpdateProfileInput struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Picture string `json:"picture"`
	Phone   string `json:"phone"`
}

func (uc *UserUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Bio != "" {
		user.Bio = input.Bio
	}
	if input.Picture != "" {
		user.Picture = input.Picture
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) GetSubscriptionHistory(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	return uc.subscriptionRepo.GetHistoryByUserID(ctx, userID)
}
