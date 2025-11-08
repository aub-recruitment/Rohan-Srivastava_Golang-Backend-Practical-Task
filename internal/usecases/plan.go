package usecases

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"
	"github.com/google/uuid"
)

type PlanUseCase struct {
	planRepo repositories.PlanRepository
}

func NewPlanUseCase(planRepo repositories.PlanRepository) *PlanUseCase {
	return &PlanUseCase{planRepo: planRepo}
}

type CreatePlanInput struct {
	Name              string             `json:"name" binding:"required"`
	Price             *int64             `json:"price" binding:"required,gte=0"`
	ValidityDays      int                `json:"validity_days" binding:"required"`
	AccessLevel       domain.AccessLevel `json:"access_level" binding:"required,oneof=free basic premium"`
	MaxDevicesAllowed int                `json:"max_devices_allowed" binding:"required"`
	Resolution        string             `json:"resolution"`
	Description       string             `json:"description"`
	IsActive          bool               `json:"is_active"`
}

func (uc *PlanUseCase) CreatePlan(ctx context.Context, input CreatePlanInput) (*domain.Plan, error) {
	plan := &domain.Plan{
		ID:                uuid.New(),
		Name:              input.Name,
		Price:             *input.Price,
		ValidityDays:      input.ValidityDays,
		AccessLevel:       input.AccessLevel,
		MaxDevicesAllowed: input.MaxDevicesAllowed,
		Resolution:        input.Resolution,
		Description:       input.Description,
		IsActive:          input.IsActive,
	}
	if err := uc.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (uc *PlanUseCase) GetPlan(ctx context.Context, planID uuid.UUID) (*domain.Plan, error) {
	return uc.planRepo.GetByID(ctx, planID)
}

func (uc *PlanUseCase) ListPlans(ctx context.Context, activeOnly bool) ([]*domain.Plan, error) {
	return uc.planRepo.List(ctx, activeOnly)
}

func (uc *PlanUseCase) UpdatePlan(ctx context.Context, planID uuid.UUID, input CreatePlanInput) (*domain.Plan, error) {
	plan, err := uc.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	plan.Name = input.Name
	plan.Price = *input.Price
	plan.ValidityDays = input.ValidityDays
	plan.AccessLevel = input.AccessLevel
	plan.MaxDevicesAllowed = input.MaxDevicesAllowed
	plan.Resolution = input.Resolution
	plan.Description = input.Description
	plan.IsActive = input.IsActive
	if err := uc.planRepo.Update(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (uc *PlanUseCase) DeletePlan(ctx context.Context, planID uuid.UUID) error {
	return uc.planRepo.Delete(ctx, planID)
}
