package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ---------- Mock definitions ----------

type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- JWT Service Mock ---
type MockJWTService struct{ mock.Mock }

func (m *MockJWTService) GenerateToken(userID uuid.UUID, email string, isAdmin bool) (string, string, error) {
	args := m.Called(userID, email, isAdmin)
	return args.String(0), args.String(1), args.Error(2)
}
func (m *MockJWTService) ValidateToken(tokenString string, refresh bool) (*infrastructure.Claims, error) {
	args := m.Called(tokenString, refresh)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.Claims), args.Error(1)
}
func (m *MockJWTService) RefreshToken(refreshToken string) (string, string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

// --- Cache Mock ---
type MockCache struct{ mock.Mock }

func (c *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := c.Called(ctx, key, value, expiration)
	return args.Error(0)
}
func (c *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := c.Called(ctx, key)
	return args.String(0), args.Error(1)
}
func (c *MockCache) Delete(ctx context.Context, key string) error {
	args := c.Called(ctx, key)
	return args.Error(0)
}
func (c *MockCache) Increment(ctx context.Context, key string) (int64, error) {
	args := c.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}
func (c *MockCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := c.Called(ctx, key, expiration)
	return args.Error(0)
}
func (c *MockCache) CheckRateLimit(ctx context.Context, identifier string, maxRequests int64, window time.Duration) (bool, error) {
	args := c.Called(ctx, identifier, maxRequests, window)
	return args.Bool(0), args.Error(1)
}
func (c *MockCache) Close() error {
	args := c.Called()
	return args.Error(0)
}

// ---------- Unit tests ----------

func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockCache := new(MockCache)

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	mockJWT.On("GenerateToken", mock.Anything, "test@example.com", false).Return("valid-token", "valid-refresh-token", nil)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWT, mockCache, 24)

	input := usecases.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	response, err := authUseCase.Register(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "valid-token", response.Token)
	assert.Equal(t, "valid-refresh-token", response.Refresh)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "Test User", response.User.Name)
	mockUserRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockCache := new(MockCache)

	existingUser := &domain.User{Email: "test@example.com"}
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWT, mockCache, 24)

	input := usecases.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	response, err := authUseCase.Register(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, domain.ErrUserExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockCache := new(MockCache)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Name:         "Test User",
		IsAdmin:      false,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockJWT.On("GenerateToken", mock.Anything, "test@example.com", false).Return("valid-token", "refresh-token", nil)
	mockCache.On("Get", mock.Anything, mock.Anything).Return("", errors.New("key not found")).Twice()
	mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil).Twice()
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWT, mockCache, 24)

	input := usecases.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := authUseCase.Login(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "valid-token", response.Token)
	mockUserRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockCache := new(MockCache)

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWT, mockCache, 24)

	input := usecases.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	response, err := authUseCase.Login(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockCache := new(MockCache)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWT, mockCache, 24)

	input := usecases.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	response, err := authUseCase.Login(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	mockUserRepo.AssertExpectations(t)
}
