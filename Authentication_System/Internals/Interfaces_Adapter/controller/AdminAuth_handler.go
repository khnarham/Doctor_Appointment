package controller

import (
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminAuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AdminAuthHandler {
	return &AdminAuthHandler{authUC: authUC}
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminloginResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	output, err := h.authUC.Login(c.Request.Context(), usecase.AdminLoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusUnauthorized
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
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
		604800, // 7 days
		"/",
		"localhost",    
		false,
		true, // HttpOnly
	)

	c.SetCookie(
		"user_role",
		string(output.Role), 
		604800,
		"/",    
		"",    
		false, 
		false,
	)

	response := loginResponseDTO{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Logged In Successfully",
		"data":    response,
	})
}