package controller

import (
	"errors"
	"net/http"
	"time"

	"my-appointment-app/Authentication_System/Internals/Usecase" 

	"github.com/gin-gonic/gin"
)

type CreateSpecializationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type SpecializationResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SpecializationHandler struct {
	createUsecase *usecase.CreateSpecializationUsecase
}

func NewSpecializationHandler(createUsecase *usecase.CreateSpecializationUsecase) *SpecializationHandler {
	return &SpecializationHandler{
		createUsecase: createUsecase,
	}
}

func (h *SpecializationHandler) Create(c *gin.Context) {
	var req CreateSpecializationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload: name and description are required",
		})
		return
	}

	input := usecase.CreateSpecializationInput{
		Name:        req.Name,
		Description: req.Description,
	}

	output, err := h.createUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrSpecializationAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create specialization",
		})
		return
	}

	res := SpecializationResponse{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "specialization created successfully",
		"data":    res,
	})
}