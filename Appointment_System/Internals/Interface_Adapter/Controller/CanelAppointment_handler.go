package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type CancelAppointmentRequest struct {
	PatientID string `json:"patient_id" binding:"required"`
}

type CancelAppointmentHandler struct {
	cancelAppointmentUC *usecase.CancelAppointmentUseCase
}

func NewCancelAppointmentHandler(uc *usecase.CancelAppointmentUseCase) *CancelAppointmentHandler {
	return &CancelAppointmentHandler{
		cancelAppointmentUC: uc,
	}
}

func (h *CancelAppointmentHandler) CancelAppointment(c *gin.Context) {

	appointmentID := c.Param("id")
	if strings.TrimSpace(appointmentID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Appointment ID is required in URL path",
		})
		return
	}

	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.CancelAppointmentInput{
		AppointmentID: appointmentID,
		PatientID:     req.PatientID,
	}

	err := h.cancelAppointmentUC.Execute(c.Request.Context(), input)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "appointment not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
			return
		}

		if strings.Contains(errMsg, "unauthorized access") {
			c.JSON(http.StatusForbidden, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment cancelled successfully",
	})
}