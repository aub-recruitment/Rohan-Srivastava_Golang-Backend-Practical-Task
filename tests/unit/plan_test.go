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

type MockPlanRepository struct {
	mock.Mock
}

func (m *MockPlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func (m *MockPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Plan), args.Error(1)
}

func (m *MockPlanRepository) GetByName(ctx context.Context, name string) (*domain.Plan, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Plan), args.Error(1)
}

func (m *MockPlanRepository) List(ctx context.Context, activeOnly bool) ([]*domain.Plan, error) {
	args := m.Called(ctx, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Plan), args.Error(1)
}

func (m *MockPlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func (m *MockPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreatePlan_Success(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	mockPlanRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Plan")).Return(nil)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)

	price := int64(999)
	input := usecases.CreatePlanInput{
		Name:              "Premium",
		Price:             &price,
		ValidityDays:      30,
		AccessLevel:       domain.AccessLevelPremium,
		MaxDevicesAllowed: 4,
		Resolution:        "4K",
		IsActive:          true,
	}

	plan, err := planUseCase.CreatePlan(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, "Premium", plan.Name)
	assert.Equal(t, int64(999), plan.Price)
	mockPlanRepo.AssertExpectations(t)
}

func TestGetPlan_Success(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	planID := uuid.New()
	expectedPlan := &domain.Plan{
		ID:    planID,
		Name:  "Premium",
		Price: 999,
	}

	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(expectedPlan, nil)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)
	plan, err := planUseCase.GetPlan(context.Background(), planID)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, expectedPlan, plan)
	mockPlanRepo.AssertExpectations(t)
}

func TestGetPlan_NotFound(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	planID := uuid.New()
	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(nil, domain.ErrPlanNotFound)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)
	plan, err := planUseCase.GetPlan(context.Background(), planID)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, domain.ErrPlanNotFound, err)
	mockPlanRepo.AssertExpectations(t)
}

func TestListPlans_Active(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	plans := []*domain.Plan{
		{ID: uuid.New(), Name: "Basic", IsActive: true},
		{ID: uuid.New(), Name: "Premium", IsActive: true},
	}

	mockPlanRepo.On("List", mock.Anything, true).Return(plans, nil)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)
	result, err := planUseCase.ListPlans(context.Background(), true)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockPlanRepo.AssertExpectations(t)
}

func TestUpdatePlan_Success(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	planID := uuid.New()
	plan := &domain.Plan{
		ID:    planID,
		Name:  "Basic",
		Price: 499,
	}

	mockPlanRepo.On("GetByID", mock.Anything, planID).Return(plan, nil)
	mockPlanRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Plan")).Return(nil)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)

	price := int64(599)
	input := usecases.CreatePlanInput{
		Name:              "Basic Updated",
		Price:             &price,
		ValidityDays:      30,
		AccessLevel:       domain.AccessLevelBasic,
		MaxDevicesAllowed: 2,
	}

	result, err := planUseCase.UpdatePlan(context.Background(), planID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockPlanRepo.AssertExpectations(t)
}

func TestDeletePlan_Success(t *testing.T) {
	mockPlanRepo := new(MockPlanRepository)

	planID := uuid.New()
	mockPlanRepo.On("Delete", mock.Anything, planID).Return(nil)

	planUseCase := usecases.NewPlanUseCase(mockPlanRepo)
	err := planUseCase.DeletePlan(context.Background(), planID)

	assert.NoError(t, err)
	mockPlanRepo.AssertExpectations(t)
}
