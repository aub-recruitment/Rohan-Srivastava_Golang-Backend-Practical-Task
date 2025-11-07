package postgres

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentRepository struct{ db *gorm.DB }

func NewContentRepository(db *gorm.DB) *ContentRepository { return &ContentRepository{db: db} }

func (r *ContentRepository) Create(ctx context.Context, content *domain.Content) error {
	return r.db.WithContext(ctx).Create(content).Error
}

func (r *ContentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Content, error) {
	var content domain.Content
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&content).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrContentNotFound
		}
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Content, int64, error) {
	var contents []*domain.Content
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Content{})
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&contents).Error
	if err != nil {
		return nil, 0, err
	}
	return contents, total, nil
}

func (r *ContentRepository) Update(ctx context.Context, content *domain.Content) error {
	return r.db.WithContext(ctx).Save(content).Error
}

func (r *ContentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Content{}, "id = ?", id).Error
}
