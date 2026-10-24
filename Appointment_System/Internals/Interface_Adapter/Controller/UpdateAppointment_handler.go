package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type UpdateAppointmentStatusRequest struct {
	DoctorID  string `json:"doctor_id" binding:"required"`
	NewStatus string `json:"new_status" binding:"required"` 
}

type UpdateAppointmentStatusHandler struct {
	updateStatusUC *usecase.UpdateAppointmentStatusUseCase
}

func NewUpdateAppointmentStatusHandler(uc *usecase.UpdateAppointmentStatusUseCase) *UpdateAppointmentStatusHandler {
	return &UpdateAppointmentStatusHandler{
		updateStatusUC: uc,
	}
}

func (h *UpdateAppointmentStatusHandler) UpdateAppointmentStatus(c *gin.Context) {

	appointmentID := c.Param("id")
	if strings.TrimSpace(appointmentID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Appointment ID is required in URL path",
		})
		return
	}

	var req UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.UpdateAppointmentStatusInput{
		AppointmentID: appointmentID,
		DoctorID:      req.DoctorID,
		NewStatus:     domain.AppointmentStatus(req.NewStatus),
	}

	err := h.updateStatusUC.Execute(c.Request.Context(), input)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "appointment not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
			return
		}

		if strings.Contains(errMsg, "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment status updated successfully",
	})
}