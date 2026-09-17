package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type GetDoctorAppointmentsHandler struct {
	uc *usecase.GetDoctorAppointmentsUseCase
}

func NewGetDoctorAppointmentsHandler(uc *usecase.GetDoctorAppointmentsUseCase) *GetDoctorAppointmentsHandler {
	return &GetDoctorAppointmentsHandler{uc: uc}
}

func (h *GetDoctorAppointmentsHandler) GetAppointments(c *gin.Context) {
	doctorID := strings.TrimSpace(c.Param("id"))

	if doctorID == "" {
		doctorID = strings.TrimSpace(c.Query("doctor_id"))
	}

	if doctorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Doctor ID is required in URL parameters",
		})
		return
	}

	input := usecase.GetDoctorAppointmentsInput{
		DoctorID: doctorID,
	}

	output, err := h.uc.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrDoctorIDRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch doctor appointments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor appointments retrieved successfully",
		"data":    output,
	})
}