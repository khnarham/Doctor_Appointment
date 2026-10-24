package controller

import (
	"my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"

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


	c.SetCookie(
		"access_token",
		output.AccessToken,
		900,
		"/",
		"", 
		false, 
		true, 
	)

	c.SetCookie(
		"refresh_token",
		output.RefreshToken,
		604800,
		"/api/v1/auth/refresh",
		"",    
		false,
		true,  
	)

      
c.SetCookie(
    "user_role",
    output.Role,
    604800,
    "/",   
    "",   
    false, 
    false, 
)


	c.JSON(http.StatusCreated, gin.H{
		"message": "We have Sent an verificaiton link to you Email",
		"data":    output,
	})

}