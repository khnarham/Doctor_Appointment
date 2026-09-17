package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domain "my-appointment-app/Authentication_System/Internals/Domain"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
)

type RegisterWalkInPatientRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Phone       string `json:"phone" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"` // YYYY-MM-DD
}

type RegisterWalkInPatientHandler struct {
	uc *usecase.RegisterWalkInPatientUsecase
}

func NewRegisterWalkInPatientHandler(uc *usecase.RegisterWalkInPatientUsecase) *RegisterWalkInPatientHandler {
	return &RegisterWalkInPatientHandler{uc: uc}
}

func (h *RegisterWalkInPatientHandler) Handle(c *gin.Context) {
	var req RegisterWalkInPatientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.RegisterWalkInPatientInput{
		Name:        req.Name,
		Phone:       req.Phone,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
	}

	output, err := h.uc.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidDateFormat),
			errors.Is(err, domain.ErrInvalidPatientName),
			errors.Is(err, domain.ErrWalkInIdentity),
			errors.Is(err, domain.ErrInvalidAge):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		case errors.Is(err, usecase.ErrPatientAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to register walk-in patient: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Walk-in patient registered successfully",
		"data":    output,
	})
}