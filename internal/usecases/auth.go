package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/infrastructure"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo     repositories.UserRepository
	jwtService   infrastructure.JWTServiceInterface
	cacheService infrastructure.CacheInterface
	expiration   int
}

func NewAuthUseCase(userRepo repositories.UserRepository, jwtService infrastructure.JWTServiceInterface, cache infrastructure.CacheInterface, expiration int) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, jwtService: jwtService, cacheService: cache, expiration: expiration}
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
	Token   string       `json:"token"`
	Refresh string       `json:"refresh"`
	User    *domain.User `json:"user"`
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

	token, refresh, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	uc.cacheService.Set(ctx, fmt.Sprintf("token:%s", user.ID), token, time.Duration(uc.expiration)*time.Hour)
	uc.cacheService.Set(ctx, fmt.Sprintf("refresh:%s", user.ID), refresh, time.Duration(uc.expiration)*time.Hour*24)

	return &AuthResponse{Token: token, Refresh: refresh, User: user}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, refresh, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	uc.resetTokenCache(ctx, token, refresh, user.ID)

	return &AuthResponse{Token: token, Refresh: refresh, User: user}, nil
}

func (uc *AuthUseCase) Refresh(ctx context.Context, userID uuid.UUID) (*AuthResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}
	token, refresh, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	uc.resetTokenCache(ctx, token, refresh, user.ID)

	return &AuthResponse{Token: token, Refresh: refresh, User: user}, nil
}

func (uc *AuthUseCase) resetTokenCache(ctx context.Context, token, refresh string, id uuid.UUID) {
	tokenKey := fmt.Sprintf("token:%s", id)
	refreshKey := fmt.Sprintf("refresh:%s", id)
	if _, err := uc.cacheService.Get(ctx, tokenKey); err != nil {
		uc.cacheService.Delete(ctx, tokenKey)
	}
	if _, err := uc.cacheService.Get(ctx, refreshKey); err != nil {
		uc.cacheService.Delete(ctx, refreshKey)
	}
	uc.cacheService.Set(ctx, fmt.Sprintf("token:%s", id), token, time.Duration(uc.expiration)*time.Hour)
	uc.cacheService.Set(ctx, fmt.Sprintf("refresh:%s", id), refresh, time.Duration(uc.expiration)*time.Hour*24)
}
