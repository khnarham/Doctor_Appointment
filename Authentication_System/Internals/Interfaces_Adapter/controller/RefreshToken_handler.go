package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	jwtPkg "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"
	"my-appointment-app/Authentication_System/Internals/Usecase"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


type RefreshTokenHandler struct {
	refreshUsecase *usecase.RefreshTokenUsecase
}

func NewRefreshTokenHandler(refreshUsecase *usecase.RefreshTokenUsecase ) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		refreshUsecase: refreshUsecase,
	}
}

func (h *RefreshTokenHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"details": "refresh_token field is required",
		})
		return
	}

	tokens, err := h.refreshUsecase.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwtPkg.ErrExpiredToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "refresh token has expired, please login again",
				"data" : nil,
			})
		case errors.Is(err, jwtPkg.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "provided refresh token is malformed or invalid",
				"data": nil,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to refresh token pair",
				"data": nil,
			})
		}
		return
	}

    
	 response := RefreshTokenResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.AccessToken,
	 }
	c.JSON(http.StatusOK, gin.H{
        "message": "User LoggedIn successFully",
		"data": response,
	})
}