package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type GetAvailableSlotsQueryRequest struct {
    DoctorID   string `json:"doctor_id" binding:"required"`
    TargetDate string `json:"target_date" binding:"required"`
}

type GetAvailableSlotsHandler struct {
	getAvailableSlotsUC *usecase.GetAvailableSlotsUseCase
}

func NewGetAvailableSlotsHandler(uc *usecase.GetAvailableSlotsUseCase) *GetAvailableSlotsHandler {
	return &GetAvailableSlotsHandler{
		getAvailableSlotsUC: uc,
	}
}

func (h *GetAvailableSlotsHandler) GetAvailableSlots(c *gin.Context) {
	var req GetAvailableSlotsQueryRequest


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters: " + err.Error(),
		})
		return
	}

	input := usecase.GetAvailableSlotsInput{
		DoctorID:   req.DoctorID,
		TargetDate: req.TargetDate,
	}

	output, err := h.getAvailableSlotsUC.Execute(c.Request.Context(), input)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "invalid date format") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Available slots fetched successfully",
		"data":    output,
	})
}