package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

func (h *LogoutHandler) Logout(c *gin.Context) {

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.SetCookie("user_role", "", -1, "/", "", true, true)


	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}