package controller

import (
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetDoctorsHandler struct {
	usecase usecase.GetdoctorUsecase
}

func NewDoctorHandler(u usecase.GetdoctorUsecase) *GetDoctorsHandler {
	return &GetDoctorsHandler{usecase: u}
}

func (h *GetDoctorsHandler) GetDoctors(c *gin.Context) {
	specializationID := c.Query("specializationID")
	search := c.Query("search")

	maxFee, _ := strconv.ParseFloat(c.Query("maxFee"), 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	filter := usecase.DoctorFilter{
		SpecializationID: specializationID,
		MaxFee:           maxFee,
		Search:           search,
		Limit:            limit,
		Offset:           offset,
	}

	doctors, err := h.usecase.ListDoctors(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch doctors list",
			"error":   err.Error(),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"doctors": doctors,
		},
	})
}