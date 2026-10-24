package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type CancelScheduleRequest struct {
	DoctorID string `json:"doctor_id" binding:"required"`
}

type CancelScheduleHandler struct {
	cancelScheduleUC *usecase.CancelScheduleUsecase
}

func NewCancelScheduleHandler(uc *usecase.CancelScheduleUsecase) *CancelScheduleHandler {
	return &CancelScheduleHandler{
		cancelScheduleUC: uc,
	}
}

func (h *CancelScheduleHandler) CancelSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if strings.TrimSpace(scheduleID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Schedule ID is required in URL path",
		})
		return
	}

	var req CancelScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.CancelScheduleInput{
		ScheduleID: scheduleID,
		DoctorID:   req.DoctorID,
	}

	err := h.cancelScheduleUC.Execute(c.Request.Context(), input)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "schedule not found") {
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
		"message": "Schedule cancelled successfully",
	})
}