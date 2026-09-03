package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type loginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponseDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	UserID       string    `json:"user_id"`
	Role  domain.Role       `josn:"role"`
	LoggedInAt   time.Time `json:"logged_in_at"`
}

type LoginHandler struct {
	loginUsecase *usecase.LoginUsecase
}

func NewLoginHandler(lu *usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		loginUsecase: lu,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req loginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid input parameters",
			"details": err.Error(),
		})
		return
	}

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

    fmt.Printf("Email: %s | Password: %s\n", input.Email, input.Password)
	output, err := h.loginUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	response := loginResponseDTO{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		UserID:       output.UserID,
		Role: output.Role,
		LoggedInAt:   output.LoggedInAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *LoginHandler) handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrEmailNotVerified):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}