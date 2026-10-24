package controller

import (
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin" // Standard Gin engine integration
)

type AnalyticsHandler struct {
	usecase usecase.AnalyticsUsecase
}

func NewAnalyticsHandler(u usecase.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{usecase: u}
}

// GetFinancialOverview GODOC
// @Summary Get Admin Financial Overview & Top Doctors
// @Tags Admin Analytics
// @Produce json
// @Param start_date query string false "YYYY-MM-DD"
// @Param end_date query string false "YYYY-MM-DD"
// @Param limit query int false "Top doctors limit"
// @Success 200 {object} domain.FinancialOverviewResponse
// @Router /api/v1/analytics/financial-overview [get]
func (h *AnalyticsHandler) GetFinancialOverview(c *gin.Context) {
	var filter domain.FinancialFilterDTO

	// Query params parsing
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.Query("limit")

	if startDateStr != "" {
		if parsedStart, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &parsedStart
		}
	}

	if endDateStr != "" {
		if parsedEnd, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &parsedEnd
		}
	}

	if limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limitVal
		}
	}

	overview, err := h.usecase.GetFinancialOverview(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch financial analytics: " + err.Error(),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overview,
	})
}