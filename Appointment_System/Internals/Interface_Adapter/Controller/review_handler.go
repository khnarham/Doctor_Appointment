// Controller / Handler Layer
package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	usecase "my-appointment-app/Appointment_System/Internals/Usecase"
)

type SubmitReviewRequest struct {
	PatientID     string `json:"patient_id" binding:"required"`
	AppointmentID string `json:"appointment_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required"`
	Comment       string `json:"comment"`
}

type SubmitReviewHandler struct {
	uc *usecase.SubmitReviewUsecase
}

func NewSubmitReviewHandler(uc *usecase.SubmitReviewUsecase) *SubmitReviewHandler {
	return &SubmitReviewHandler{uc: uc}
}

func (h *SubmitReviewHandler) SubmitReview(c *gin.Context) {
	var req SubmitReviewRequest

	// JSON Request binding
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	input := usecase.SubmitReviewInput{
		AuthenticatedPatientID: req.PatientID, // Direct Request Body se ID pass kar rahe hain
		AppointmentID:          req.AppointmentID,
		Rating:                 req.Rating,
		Comment:                req.Comment,
	}

	err := h.uc.Execute(c.Request.Context(), input)
	if err != nil {
		// Specific Domain Error Handling
		switch {
		case errors.Is(err, usecase.ErrInvalidRatingRange):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrAppointmentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrUnauthorizedReview):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrAppointmentNotInCompletedState):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrReviewAlreadySubmitted):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrReviewPeriodExpired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to submit review: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review submitted successfully",
	})
}