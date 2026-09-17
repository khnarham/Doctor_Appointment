package controller

import (
	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDoctorSchedulesHandler struct {
	getDoctorSchedulesUC *usecase.GetDoctorSchedulesUseCase
}

func NewGetDoctorSchedulesHandler(uc *usecase.GetDoctorSchedulesUseCase) *GetDoctorSchedulesHandler {
	return &GetDoctorSchedulesHandler{
		getDoctorSchedulesUC: uc,
	}
}

func (h *GetDoctorSchedulesHandler) GetDoctorSchedulesWithSlots(c *gin.Context) {
	doctorID := c.Param("id")
	if doctorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "doctor ID is required in URL path parameter",
		})
		return
	}

	output, err := h.getDoctorSchedulesUC.Execute(c.Request.Context(), doctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}