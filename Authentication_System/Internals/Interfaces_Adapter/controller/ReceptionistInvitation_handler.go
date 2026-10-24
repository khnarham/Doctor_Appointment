package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
    usecase "my-appointment-app/Authentication_System/Internals/Usecase"
)

type InviteReceptionistRequest struct {
	DoctorID string `json:"doctor_id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type InviteReceptionistResponse struct {
	InvitationID string `json:"invitation_id"`
	Status       string `json:"status"`
}

type ReceptionistHandler struct {
	inviteUsecase *usecase.InviteReceptionistUsecase
}

func NewReceptionistHandler(inviteUsecase *usecase.InviteReceptionistUsecase) *ReceptionistHandler {
	return &ReceptionistHandler{
		inviteUsecase: inviteUsecase,
	}
}

func (h *ReceptionistHandler) InviteReceptionist(c *gin.Context) {
	var req InviteReceptionistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"details": err.Error(),
		})
		return
	}

	input := usecase.InviteReceptionistInput{
		DoctorID: req.DoctorID,
		Email:    req.Email,
	}

	output, err := h.inviteUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error() , "data": nil})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	resp := InviteReceptionistResponse{
		InvitationID: output.InvitationID,
		Status:       output.Status,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invitation Send Successfully",
		"data": resp,
	})
}