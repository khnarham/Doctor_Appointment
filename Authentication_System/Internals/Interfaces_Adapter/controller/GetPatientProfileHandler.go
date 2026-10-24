package controller

import (
	"errors"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPatientProfileHandler struct {
	getPatientProfileUC *usecase.GetPatientProfileUseCase
}

func NewGetPatientProfileHandler(uc *usecase.GetPatientProfileUseCase) *GetPatientProfileHandler {
	return &GetPatientProfileHandler{
		getPatientProfileUC: uc,
	}
}

func (h *GetPatientProfileHandler) GetProfile(c *gin.Context) {
	// 1. Middleware context se Auth ID pull karna
	// Assumed key: "user_id" (Apne middleware ki key ke mutabiq adjust kar lena)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user identity not found in context",
		})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid user identity type",
		})
		return
	}

	// 2. Execute Usecase
	input := usecase.GetPatientProfileInput{
		UserID: userID,
	}

	output, err := h.getPatientProfileUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrPatientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Patient profile not found for this account",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Response Return
	c.JSON(http.StatusOK, gin.H{
		"message": "Patient profile retrieved successfully",
		"data":    output,
	})
}