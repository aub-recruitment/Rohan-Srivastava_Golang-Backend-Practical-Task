package handlers

import (
	"net/http"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler struct {
	planUseCase *usecases.PlanUseCase
}

// @name NewPlanHandler - Creates new instance of plan handler
// @param planUseCase - plan usecase (service)
// @returns - new instance of plan handler
func NewPlanHandler(planUseCase *usecases.PlanUseCase) *PlanHandler {
	return &PlanHandler{planUseCase: planUseCase}
}

// @name CreatePlan - Admin API creates new plan
// @param c - gin context
// @returns - newly created plan
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var input usecases.CreatePlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.planUseCase.CreatePlan(c.Request.Context(), input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

// @name GetPlan - Open API to get plan with given ID
// @param c - gin context
// @returns - plan with id
func (h *PlanHandler) GetPlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}
	plan, err := h.planUseCase.GetPlan(c.Request.Context(), planID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// @name ListPlans - Open API to get list of all plans
// @param c - gin context
// @query active - boolean value to only return active plans
// @returns - plans slice
func (h *PlanHandler) ListPlans(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "true") == "true"
	plans, err := h.planUseCase.ListPlans(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// @name UpdatePlan - Admin API to update plan with given ID
// @param c - gin context
// @returns - updated plan with id
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}
	var input usecases.CreatePlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.planUseCase.UpdatePlan(c.Request.Context(), planID, input)
	if err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

// @name GetContent - Admin API to delete plan with given ID
// @param c - gin context
// @returns - deletion successful message
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}
	if err := h.planUseCase.DeletePlan(c.Request.Context(), planID); err != nil {
		c.JSON(getErrorStatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "plan deleted successfully"})
}
