package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type GetPatientAppointmentsHandler struct {
	uc *usecase.GetPatientAppointmentsUseCase
}

func NewGetPatientAppointmentsHandler(uc *usecase.GetPatientAppointmentsUseCase) *GetPatientAppointmentsHandler {
	return &GetPatientAppointmentsHandler{uc: uc}
}

func (h *GetPatientAppointmentsHandler) GetAppointments(c *gin.Context) {
	patientID := strings.TrimSpace(c.Param("id"))

	if patientID == "" {
		patientID = strings.TrimSpace(c.Query("patient_id"))
	}

	if patientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Patient ID is required in URL parameters",
		})
		return
	}

	input := usecase.GetPatientAppointmentsInput{
		PatientID: patientID,
	}

	output, err := h.uc.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch appointments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointments retrieved successfully",
		"data":    output,
	})
}