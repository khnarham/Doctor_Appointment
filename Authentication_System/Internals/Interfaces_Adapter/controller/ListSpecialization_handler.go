package controller

import (
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"
		"github.com/gin-gonic/gin"

)

type ListSpecializationHandler struct {
	listUC *usecase.ListSpecializationsUsecase
}

func NewListSpecializationHandler(listUC *usecase.ListSpecializationsUsecase) *ListSpecializationHandler {
	return &ListSpecializationHandler{
		listUC: listUC,
	}
}

func (h *ListSpecializationHandler) ListSpecializations(c *gin.Context) {
	output, err := h.listUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch specializations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}