package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTServiceInterface defines the interface for JWT operations
type JWTServiceInterface interface {
	GenerateToken(userID uuid.UUID, email string, isAdmin bool) (string, string, error)
	ValidateToken(tokenString string, refresh bool) (*Claims, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type JWTService struct {
	secretKey   string
	secretSauce string
	expiration  int
}

type Claims struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	IsAdmin bool      `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey, secretSauce string, expiration int) *JWTService {
	return &JWTService{secretKey: secretKey, secretSauce: secretSauce, expiration: expiration}
}

func (j *JWTService) GenerateToken(userID uuid.UUID, email string, isAdmin bool) (string, string, error) {
	claims := Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiration) * time.Hour * 24))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshString, err := token.SignedString([]byte(j.secretSauce))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshString, nil
}

func (j *JWTService) ValidateToken(tokenString string, refresh bool) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if refresh {
			return []byte(j.secretSauce), nil
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (j *JWTService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := j.ValidateToken(refreshToken, true)
	if err != nil {
		return "", "", err
	}
	return j.GenerateToken(claims.UserID, claims.Email, claims.IsAdmin)
}
