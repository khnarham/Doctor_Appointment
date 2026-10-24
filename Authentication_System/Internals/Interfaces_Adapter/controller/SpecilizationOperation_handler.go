package controller

import (
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SpecializationOperationHandler struct {
	usecase *usecase.SpecializationUsecase
}

func NewSpecializationOperationHandler(uc *usecase.SpecializationUsecase) *SpecializationOperationHandler {
	return &SpecializationOperationHandler{usecase: uc}
}

type updateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type specializationResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (h *SpecializationOperationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	spec, err := h.usecase.Update(c.Request.Context(), usecase.UpdateSpecializationInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidSpecializationName):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrSpecializationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrSpecializationNameExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "specialization updated successfully",
		"data": specializationResponse{
			ID:          spec.ID(),
			Name:        spec.Name(),
			Description: spec.Description(),
			CreatedAt:   spec.CreatedAt().Format(time.RFC3339),
			UpdatedAt:   spec.UpdatedAt().Format(time.RFC3339),
		},
	})
}

func (h *SpecializationOperationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSpecializationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrSpecializationInUse):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "specialization deleted successfully"})
}