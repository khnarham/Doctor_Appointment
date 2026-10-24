package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
)

type DoctorResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	SpecializationID string    `json:"specialization_id"`
	LicenseNumber    string    `json:"license_number"`
	ConsultationFee  float64   `json:"consultation_fee"`
	Status           string    `json:"status"`
	IsAvailable      bool      `json:"is_available"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type GetDoctorHandler struct {
	doctorUsecase *usecase.DoctorUsecase
}

func NewGetDoctorHandler(doctorUsecase *usecase.DoctorUsecase) *GetDoctorHandler {
	return &GetDoctorHandler{
		doctorUsecase: doctorUsecase,
	}
}

func (h *GetDoctorHandler) GetDoctorProfile(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access: user context missing",
			"data":    nil,
		})
		return
	}

	userID, ok := val.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user identity in request context",
			"data":    nil,
		})
		return
	}

	doctor, err := h.doctorUsecase.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "User ID is invalid",
				"data":    nil,
			})

		case errors.Is(err, usecase.ErrDoctorNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Doctor profile not found",
				"data":    nil,
			})

		case errors.Is(err, usecase.ErrDoctorNotAvailable):
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Doctor account is currently unavailable",
				"data":    nil,
			})

		default:
			log.Printf("[ERROR] GetDoctorProfile internal failure: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to retrieve doctor profile",
				"data":    nil,
			})
		}
		return
	}

	// 3. Response DTO Construction
	resp := DoctorResponse{
		ID:               doctor.ID(),
		UserID:           doctor.UserID(),
		SpecializationID: doctor.SpecializationID(),
		LicenseNumber:    doctor.LicenseNumber(),
		ConsultationFee:  doctor.ConsultationFee(),
		Status:           string(doctor.Status()),
		IsAvailable:      doctor.IsAvailable(),
		CreatedAt:        doctor.CreatedAt(),
		UpdatedAt:        doctor.UpdatedAt(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Doctor profile fetched successfully",
		"data":    resp,
	})
}