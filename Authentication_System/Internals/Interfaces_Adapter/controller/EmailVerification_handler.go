package controller

import (
	"my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type EmailVerificationGinHandler struct {
	usecase *usecase.EmailVerificationUsecase
}

func NewEmailVerificationGinHandler(u *usecase.EmailVerificationUsecase) *EmailVerificationGinHandler {
	return &EmailVerificationGinHandler{usecase: u}
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *EmailVerificationGinHandler) VerifyEmail(c *gin.Context) {
	 token := c.Query("token")
     if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token is required"})
		return
	}
	verificationInput := usecase.VerficationInput{
		Token: token,
	}
	output, err := h.usecase.Execute(c.Request.Context(), verificationInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() , "data": nil})
		return
	}

isProduction := os.Getenv("GIN_MODE") == "release" || os.Getenv("ENV") == "production"
cookieDomain := os.Getenv("COOKIE_DOMAIN")

// Cross-site (Vercel frontend + separate backend) ke liye SameSite None lazmi hai
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
    Value:    output.Role,
    Path:     "/",
    Domain:   cookieDomain,
    MaxAge:   604800,
    HttpOnly: false,
    Secure:   isProduction,
    SameSite: sameSiteMode,
})

	c.JSON(http.StatusCreated, gin.H{
		"message": "We have Sent an verificaiton link to you Email",
		"data":    output,
	})

}