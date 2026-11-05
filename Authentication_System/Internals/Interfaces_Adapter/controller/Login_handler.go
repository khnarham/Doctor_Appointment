package controller

import (
	"errors"
	"net/http"
	"os"
	"time"

	"my-appointment-app/Authentication_System/Internals/Usecase"

	"github.com/gin-gonic/gin"
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

isProduction := os.Getenv("GIN_MODE") == "release" || os.Getenv("ENV") == "production"
cookieDomain := os.Getenv("COOKIE_DOMAIN") // Vercel ke liye "" (blank) rahe ga

// Cross-site (Railway backend + Vercel frontend) ke liye SameSite None lazmi hai
sameSiteMode := http.SameSiteLaxMode
if isProduction {
    sameSiteMode = http.SameSiteNoneMode
}

// 1. Access Token Cookie
http.SetCookie(c.Writer, &http.Cookie{
    Name:     "access_token",
    Value:    output.AccessToken,
    Path:     "/",
    Domain:   cookieDomain,
    MaxAge:   900,
    HttpOnly: true,
    Secure:   isProduction,
    SameSite: sameSiteMode,
})

// 2. Refresh Token Cookie
http.SetCookie(c.Writer, &http.Cookie{
    Name:     "refresh_token",
    Value:    output.RefreshToken,
    Path:     "/",
    Domain:   cookieDomain,
    MaxAge:   604800,
    HttpOnly: true,
    Secure:   isProduction,
    SameSite: sameSiteMode,
})

// 3. User Role Cookie
http.SetCookie(c.Writer, &http.Cookie{
    Name:     "user_role",
    Value:    output.User.Role,
    Path:     "/",
    Domain:   cookieDomain,
    MaxAge:   604800,
    HttpOnly: false,
    Secure:   isProduction,
    SameSite: sameSiteMode,
})
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