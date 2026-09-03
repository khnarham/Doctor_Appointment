package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type AdminDoctorHandler struct {
	inviteUseCase *usecase.InviteDoctorUsecase
}
func NewAdminDoctorHandler(inviteUsecase *usecase.InviteDoctorUsecase) *AdminDoctorHandler {
    return &AdminDoctorHandler{
        inviteUseCase: inviteUsecase,
    }
}

type InviteDoctorRequest struct {
	Email            string `json:"email" binding:"required,email"`
	SpecializationID string `json:"specialization_id" binding:"required,uuid"`
}

func (h *AdminDoctorHandler) InviteDoctor(c *gin.Context) {
	var req InviteDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input parameters", "details": err.Error()})
		return
	}

	input := usecase.InviteDoctorInput{
		Email:            req.Email,
		SpecializationID: req.SpecializationID,
	}

	output, err := h.inviteUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrInviteAlreadyActive):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err , usecase.ErrSpecializationAlreadyExists):
		    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Doctor invitation sent successfully",
		"data":    output,
	})
}