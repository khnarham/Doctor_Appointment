package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type UpdateScheduleRequest struct {
	DoctorID     string `json:"doctor_id" binding:"required"`
	Day          string `json:"day" binding:"required"`      
	StartTime    string `json:"start_time" binding:"required"`   
	EndTime      string `json:"end_time" binding:"required"`     
	SlotDuration int    `json:"slot_duration" binding:"required,gt=0"`
}

type UpdateScheduleHandler struct {
	updateScheduleUC *usecase.UpdateScheduleUsecase
}

func NewUpdateScheduleHandler(uc *usecase.UpdateScheduleUsecase) *UpdateScheduleHandler {
	return &UpdateScheduleHandler{
		updateScheduleUC: uc,
	}
}

func (h *UpdateScheduleHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
	if strings.TrimSpace(scheduleID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Schedule ID is required in URL path",
		})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	dayOfWeek, err := parseDayOfWeek(req.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startH, startM, err := parseHHMM(req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format: expected HH:MM"})
		return
	}

	endH, endM, err := parseHHMM(req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format: expected HH:MM"})
		return
	}

	input := usecase.UpdateScheduleInput{
		ScheduleID:   scheduleID,
		DoctorID:     req.DoctorID,
		Day:          dayOfWeek,
		StartHour:    startH,
		StartMinute:  startM,
		EndHour:      endH,
		EndMinute:    endM,
		SlotDuration: time.Duration(req.SlotDuration) * time.Minute,
	}

	output, err := h.updateScheduleUC.Execute(c.Request.Context(), input)
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
		"message": "Schedule updated successfully",
		"data":    output,
	})
}


