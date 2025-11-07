package usecases

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"

	"github.com/google/uuid"
)

type ContentUseCase struct {
	contentRepo      repositories.ContentRepository
	subscriptionRepo repositories.SubscriptionRepository
	userRepo         repositories.UserRepository
}

func NewContentUseCase(contentRepo repositories.ContentRepository, subscriptionRepo repositories.SubscriptionRepository, userRepo repositories.UserRepository) *ContentUseCase {
	return &ContentUseCase{contentRepo: contentRepo, subscriptionRepo: subscriptionRepo, userRepo: userRepo}
}

type CreateContentInput struct {
	Title           string             `json:"title" binding:"required"`
	Description     string             `json:"description"`
	AccessLevel     domain.AccessLevel `json:"access_level" binding:"required"`
	DurationSeconds int                `json:"duration_seconds" binding:"required"`
	ThumbnailURL    string             `json:"thumbnail_url"`
	VideoURL        string             `json:"video_url"`
	Published       bool               `json:"published"`
}

func (uc *ContentUseCase) CreateContent(ctx context.Context, input CreateContentInput) (*domain.Content, error) {
	content := &domain.Content{
		ID:              uuid.New(),
		Title:           input.Title,
		Description:     input.Description,
		AccessLevel:     input.AccessLevel,
		DurationSeconds: input.DurationSeconds,
		ThumbnailURL:    input.ThumbnailURL,
		VideoURL:        input.VideoURL,
		Published:       input.Published,
	}
	if err := uc.contentRepo.Create(ctx, content); err != nil {
		return nil, err
	}
	return content, nil
}

func (uc *ContentUseCase) GetContent(ctx context.Context, contentID uuid.UUID, userID *uuid.UUID) (*domain.Content, error) {
	content, err := uc.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, err
	}
	if !content.Published {
		return nil, domain.ErrContentNotPublished
	}
	if userID != nil {
		hasAccess, err := uc.CheckAccess(ctx, *userID, content.AccessLevel)
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, domain.ErrContentNotAccessible
		}
	} else {
		if content.AccessLevel != domain.AccessLevelFree {
			return nil, domain.ErrContentNotAccessible
		}
	}
	return content, nil
}

func (uc *ContentUseCase) ListContent(ctx context.Context, publishedOnly bool, limit, offset int) ([]*domain.Content, int64, error) {
	filters := make(map[string]interface{})
	if publishedOnly {
		filters["published"] = true
	}
	return uc.contentRepo.List(ctx, filters, limit, offset)
}

func (uc *ContentUseCase) UpdateContent(ctx context.Context, contentID uuid.UUID, input CreateContentInput) (*domain.Content, error) {
	content, err := uc.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, err
	}
	content.Title = input.Title
	content.Description = input.Description
	content.AccessLevel = input.AccessLevel
	content.DurationSeconds = input.DurationSeconds
	content.ThumbnailURL = input.ThumbnailURL
	content.VideoURL = input.VideoURL
	content.Published = input.Published
	if err := uc.contentRepo.Update(ctx, content); err != nil {
		return nil, err
	}
	return content, nil
}

func (uc *ContentUseCase) DeleteContent(ctx context.Context, contentID uuid.UUID) error {
	return uc.contentRepo.Delete(ctx, contentID)
}

func (uc *ContentUseCase) CheckAccess(ctx context.Context, userID uuid.UUID, requiredLevel domain.AccessLevel) (bool, error) {
	if requiredLevel == domain.AccessLevelFree {
		return true, nil
	}
	subscription, err := uc.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		if err == domain.ErrSubscriptionNotFound {
			return false, nil
		}
		return false, err
	}
	if !subscription.IsActive || subscription.IsExpired() {
		return false, nil
	}
	if requiredLevel == domain.AccessLevelBasic {
		return subscription.Plan.AccessLevel == domain.AccessLevelBasic ||
			subscription.Plan.AccessLevel == domain.AccessLevelPremium, nil
	}
	if requiredLevel == domain.AccessLevelPremium {
		return subscription.Plan.AccessLevel == domain.AccessLevelPremium, nil
	}
	return false, nil
}
