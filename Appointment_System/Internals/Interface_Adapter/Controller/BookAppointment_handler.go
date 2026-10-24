package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type BookAppointmentRequest struct {
	PatientID string `json:"patient_id" binding:"required"`
	DoctorID  string `json:"doctor_id" binding:"required"`
	Date      string `json:"date" binding:"required"`      
	StartTime string `json:"start_time" binding:"required"` 
}

type BookAppointmentHandler struct {
	bookAppointmentUC *usecase.BookAppointmentUseCase
}

func NewBookAppointmentHandler(uc *usecase.BookAppointmentUseCase) *BookAppointmentHandler {
	return &BookAppointmentHandler{
		bookAppointmentUC: uc,
	}
}

func (h *BookAppointmentHandler) BookAppointment(c *gin.Context) {
	var req BookAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.BookAppointmentInput{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		Date:      req.Date,
		StartTime: req.StartTime,
	}

	output, err := h.bookAppointmentUC.Execute(c.Request.Context(), input)
	if err != nil {

		if errors.Is(err, domain.ErrSlotAlreadyBooked) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "This slot is already booked",
			})
			return
		}

		if errors.Is(err, domain.ErrSlotNotInSchedule) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Requested time slot does not exist in doctor's active schedule",
			})
			return
		}

		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid date format") || strings.Contains(errMsg, "cannot book appointment for past dates") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment booked successfully",
		"data":    output,
	})
}