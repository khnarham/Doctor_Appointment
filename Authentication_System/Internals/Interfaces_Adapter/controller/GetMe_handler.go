package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type meResponseDTO struct {
	User usecase.GetMeOutput `json:"user"`
}

type MeHandler struct {
	getMeUsecase *usecase.GetMeUsecase
}

func NewMeHandler(gmu *usecase.GetMeUsecase) *MeHandler {
	return &MeHandler{
		getMeUsecase: gmu,
	}
}

func (h *MeHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user session"})
		return
	}

	output, err := h.getMeUsecase.Execute(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "user profile not found" , "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error" , "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Get Successfully",
		"data": output,
	})
}