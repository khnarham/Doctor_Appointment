package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
)

type AcceptInvitationRequest struct {
	Token    string `json:"token" binding:"required"`
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

type AcceptInvitationResponse struct {
	ReceptionistID string `json:"receptionist_id"`
	Email          string `json:"email"`
	DoctorID       string `json:"doctor_id"`
}

type AcceptInvitationHandler struct {
	acceptUsecase *usecase.AcceptInvitationUsecase
}

func NewAcceptInvitationHandler(acceptUsecase *usecase.AcceptInvitationUsecase) *AcceptInvitationHandler {
	return &AcceptInvitationHandler{
		acceptUsecase: acceptUsecase,
	}
}

func (h *AcceptInvitationHandler) AcceptInvitation(c *gin.Context) {
	var req AcceptInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload or missing required fields",
			"error":   err.Error(),
			"data":    nil,
		})
		return
	}

	input := usecase.AcceptInvitationInput{
		Token:    req.Token,
		FullName: req.FullName,
		Password: req.Password,
	}

	output, err := h.acceptUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvitationNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Invitation token not found",
				"data":    nil,
			})
		case errors.Is(err, usecase.ErrInvitationInvalid):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invitation token is either expired or already used",
				"data":    nil,
			})

			case errors.Is(err, usecase.ErrPasswordHashing):
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "An error occurred while processing authentication credentials",
                "data":    nil,
            })

        case errors.Is(err, usecase.ErrReceptionistCreationFailed):
            c.JSON(http.StatusNotFound, gin.H{
                "message": "Failed to create receptionist record in database",
                "data":    nil,
            })

        case errors.Is(err, usecase.ErrTokenUpdateFailed):
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "Failed to update invitation status",
                "data":    nil,
            })
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to process invitation request",
				"data":    nil,
			})
		}
		return
	}

	resp := AcceptInvitationResponse{
		ReceptionistID: output.ReceptionistID,
		Email:          output.Email,
		DoctorID:       output.DoctorID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Receptionist registered successfully",
		"data":    resp,
	})
}