package controller

import (
	"log"
	"net/http"

	"my-appointment-app/Appointment_System/Internals/Usecase"

	"github.com/gin-gonic/gin"
)

type DoctorAnalyticsHandler struct {
	usecase *usecase.DoctorAnalyticsUsecase
}

func NewDoctorAnalyticsHandler(uc *usecase.DoctorAnalyticsUsecase) *DoctorAnalyticsHandler {
	 return &DoctorAnalyticsHandler{usecase: uc}

}

func (h *DoctorAnalyticsHandler) GetDashboard(c *gin.Context) {
	doctorID := c.Param("id")
	if doctorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "doctor ID path parameter is required",
		})
		return
	}

	response, err := h.usecase.GetDashboardData(c.Request.Context(), doctorID)
	if err != nil {
		// Terminal log to see precise PGX scan errors if any occur
		log.Printf("[Doctor Analytics Error]: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to fetch doctor dashboard analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}