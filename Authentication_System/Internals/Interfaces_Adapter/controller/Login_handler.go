package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type loginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponseDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User        usecase.UserOutput  `json:"user"`
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
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUsecase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	} 


	c.SetCookie(
		"access_token",
		output.AccessToken,
		900,
		"/",
		"localhost", 
		false, 
		true, 
	)

	c.SetCookie(
		"refresh_token",
		output.RefreshToken,
		604800,
		"/",
		"localhost",    
		false,
		true,  
	)

      
c.SetCookie(
    "user_role",
    output.User.Role,
    604800,
    "/",   
    "",   
    false, 
    false, 
)

	response := loginResponseDTO{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User: output.User,
		LoggedInAt:   output.LoggedInAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Logged In Successfully",
		"data": response,
	})
}

func (h *LoginHandler) handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() , "data": nil})
	case errors.Is(err, usecase.ErrEmailNotVerified):
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error() , "data": nil})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}