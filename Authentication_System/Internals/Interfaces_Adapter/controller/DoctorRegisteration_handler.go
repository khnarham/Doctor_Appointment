package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type DoctorRegistrationHandler struct {
	registerUseCase *usecase.RegisterDoctorWithInviteUsecase
}

func NewDoctorRegistrationHandler(registerUC *usecase.RegisterDoctorWithInviteUsecase) *DoctorRegistrationHandler {
	return &DoctorRegistrationHandler{
		registerUseCase: registerUC,
	}
}

type RegisterDoctorRequest struct {
	Token           string  `json:"token" binding:"required"`
	Name            string  `json:"name" binding:"required,min=2"`
	Password        string  `json:"password" binding:"required,min=8"`
	Gender          string  `json:"gender" binding:"required"`
	LicenseNumber   string  `json:"license_number" binding:"required"`
	ConsultationFee float64 `json:"consultation_fee" binding:"required,gte=0"`
}

func (h *DoctorRegistrationHandler) RegisterDoctor(c *gin.Context) {
	var req RegisterDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration payload", "details": err.Error()})
		return
	}

	input := usecase.RegisterDoctorWithInviteInput{
		Token:           req.Token,
		Name:            req.Name,
		Password:        req.Password,
		Gender:          domain.Gender(req.Gender),
		LicenseNumber:   req.LicenseNumber,
		ConsultationFee: req.ConsultationFee,
	}

	output, err := h.registerUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInviteExpired), errors.Is(err, domain.ErrInviteUsed):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.InvalidGender):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Doctor registered successfully",
		"data":    output,
	})
}