package postgres

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanRepository struct{ db *gorm.DB }

func NewPlanRepository(db *gorm.DB) *PlanRepository { return &PlanRepository{db: db} }

func (r *PlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *PlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPlanNotFound
		}
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) GetByName(ctx context.Context, name string) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPlanNotFound
		}
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) List(ctx context.Context, activeOnly bool) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("price ASC").Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *PlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Plan{}, "id = ?", id).Error
}
