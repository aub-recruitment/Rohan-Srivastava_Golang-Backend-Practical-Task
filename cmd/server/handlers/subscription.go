package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionUseCase *usecases.SubscriptionUseCase
}

// @name NewSubscriptionHandler - Creates instance of subscription handler
// @param c - gin context
// @returns - new instance of subscription handler
func NewSubscriptionHandler(subscriptionUseCase *usecases.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionUseCase: subscriptionUseCase}
}

// @name CreateSubscription - Protected API to create new subsctiption to a plan
// @param c - gin context
// @returns - Newly created subscription
// @dev - uses plan_id from plans entity, only one subscription active at a time
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	var input usecases.CreateSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subscription, err := h.subscriptionUseCase.CreateSubscription(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subscription)
}

// @name GetActiveSubscription - Protected API toget user's active subscription
// @param c - gin context
// @returns - current active subscription for user
func (h *SubscriptionHandler) GetActiveSubscription(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	subscription, err := h.subscriptionUseCase.GetActiveSubscription(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscription)
}

// @name GetSubsciptionHistory - Protected API to get user's subscription history
// @param c - gin context
// @returns - list of all user past subscriptions
func (h *SubscriptionHandler) GetSubscriptionHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	subscriptions, err := h.subscriptionUseCase.GetSubscriptionHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// @name CancelSubscription - Protected API to cancel user's current subscription
// @param c - gin context
// @returns - success message of cancellation
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	subscriptionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}
	if err := h.subscriptionUseCase.CancelSubscription(c.Request.Context(), userID, subscriptionID); err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "subscription cancelled successfully"})
}

// @name RenewSubscription - Protected API to renew previous plan subscription
// @param c - gin context
// @returns - new user subscription
func (h *SubscriptionHandler) RenewSubscription(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	subscriptionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}
	subscription, err := h.subscriptionUseCase.RenewSubscription(c.Request.Context(), userID, subscriptionID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subscription)
}
