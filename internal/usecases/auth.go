package usecases

import (
	"context"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo   repositories.UserRepository
	jwtService *infrastructure.JWTService
}

func NewAuthUseCase(userRepo repositories.UserRepository, jwtService *infrastructure.JWTService) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, jwtService: jwtService}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Bio      string `json:"bio"`
	Picture  string `json:"picture"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	existingUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, domain.ErrUserExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Name:         input.Name,
		Phone:        input.Phone,
		Bio:          input.Bio,
		Picture:      input.Picture,
		IsAdmin:      false,
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.IsAdmin, 24)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.IsAdmin, 24)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: user}, nil
}
