package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

func (h *LogoutHandler) Logout(c *gin.Context) {
	isProduction := os.Getenv("GIN_MODE") == "release" || os.Getenv("ENV") == "production"

	sameSiteMode := http.SameSiteLaxMode
	if isProduction {
		sameSiteMode = http.SameSiteNoneMode
	}

	// 1. Clear Access Token Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1, // Deletes the cookie immediately
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: sameSiteMode,
	})

	// 2. Clear Refresh Token Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: sameSiteMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "user_role",
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   isProduction,
		SameSite: sameSiteMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}