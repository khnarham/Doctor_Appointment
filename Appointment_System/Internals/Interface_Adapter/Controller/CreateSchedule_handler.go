package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type CreateScheduleRequest struct {
	DoctorID     string `json:"doctor_id" binding:"required"`
	Day          string `json:"day" binding:"required"`          
	StartTime    string `json:"start_time" binding:"required"` 
	EndTime      string `json:"end_time" binding:"required"`    
	SlotDuration int    `json:"slot_duration" binding:"required,gt=0"`
}

type ScheduleHandler struct {
	createScheduleUC *usecase.CreateScheduleUsecase
}

func NewScheduleHandler(uc *usecase.CreateScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{
		createScheduleUC: uc,
	}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req CreateScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	dayOfWeek, err := parseDayOfWeek(req.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	timeRange, err := parseTimeRange(req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	cmd := usecase.CreateScheduleInput{
		DoctorID:     req.DoctorID,
		Days:         dayOfWeek,
		TimeRange:    timeRange,
		SlotDuration: req.SlotDuration,
	}

	output, err := h.createScheduleUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Schedule created successfully",
		"data":    output,
	})
}


func parseDayOfWeek(dayStr string) (domain.DayOfWeek, error) {
	switch strings.ToLower(strings.TrimSpace(dayStr)) {
	case "sunday":
		return domain.Sunday, nil
	case "monday":
		return domain.Monday, nil
	case "tuesday":
		return domain.Tuesday, nil
	case "wednesday":
		return domain.Wednesday, nil
	case "thursday":
		return domain.Thursday, nil
	case "friday":
		return domain.Friday, nil
	case "saturday":
		return domain.Saturday, nil
	default:
		return 0, fmt.Errorf("invalid day of week: %s", dayStr)
	}
}

func parseTimeRange(startStr, endStr string) (domain.TimeRange, error) {
	startH, startM, err := parseHHMM(startStr)
	if err != nil {
		return domain.TimeRange{}, fmt.Errorf("invalid start_time: %w", err)
	}

	endH, endM, err := parseHHMM(endStr)
	if err != nil {
		return domain.TimeRange{}, fmt.Errorf("invalid end_time: %w", err)
	}

	startClock, err := domain.NewClockTime(startH, startM)
	if err != nil {
		return domain.TimeRange{}, fmt.Errorf("invalid start clock: %w", err)
	}

	endClock, err := domain.NewClockTime(endH, endM)
	if err != nil {
		return domain.TimeRange{}, fmt.Errorf("invalid end clock: %w", err)
	}

	timeRange, err := domain.NewTimeRange(startClock, endClock)
	if err != nil {
		return domain.TimeRange{}, fmt.Errorf("invalid time range: %w", err)
	}

	return timeRange, nil
}

func parseHHMM(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected format HH:MM")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid hour")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("invalid minute")
	}

	return hour, minute, nil
}