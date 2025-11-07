package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
)

func getErrorStatusCode(err error) int {
	switch err {
	case domain.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case domain.ErrUnauthorized, domain.ErrTokenExpired, domain.ErrTokenInvalid:
		return http.StatusUnauthorized
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrUserExists:
		return http.StatusConflict
	case domain.ErrUserNotFound, domain.ErrContentNotFound, domain.ErrPlanNotFound,
		domain.ErrSubscriptionNotFound, domain.ErrWatchHistoryNotFound, domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrContentNotAccessible, domain.ErrContentNotPublished:
		return http.StatusForbidden
	case domain.ErrPlanNotAvailable, domain.ErrInactivePlan:
		return http.StatusBadRequest
	case domain.ErrActiveSubscriptionExists, domain.ErrSubscriptionLimitExceeded:
		return http.StatusConflict
	case domain.ErrSubscriptionExpired, domain.ErrSubscriptionInactive:
		return http.StatusBadRequest
	case domain.ErrInvalidInput, domain.ErrValidationFailed, domain.ErrInvalidProgress:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
