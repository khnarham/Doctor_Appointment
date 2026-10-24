package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"my-appointment-app/Authentication_System/Internals/Usecase"

	"github.com/gin-gonic/gin"
)

type RegisterRequestDTO struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Gender      string `json:"gender" binding:"required,oneof=Male Female"`
	Phone       string `json:"phone" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
}

type RegisterResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type ErrorResponseDTO struct {
	Error string `json:"error"`
}

var (
	ErrInvalidInputData   = errors.New("provided input data is invalid")
)

type RegisterUserHandler struct {
	registerUsecase *usecase.RegisterUserUsecase
}

func NewRegisterUserHandler(uc *usecase.RegisterUserUsecase) *RegisterUserHandler {
	return &RegisterUserHandler{registerUsecase: uc}
}

func (h *RegisterUserHandler) Register(c *gin.Context) {
	fmt.Println("====== USECASE EXECUTE STARTED ======")
	var req RegisterRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponseDTO{Error: "invalid request payload or missing required fields"})
		return
	}
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
    if err != nil {
    }

	input := usecase.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Gender: domain.Gender(req.Gender),
		Phone: req.Phone,
		DateOfBirth: dob,
	}

	output, err := h.registerUsecase.Execute(c.Request.Context(), input)
	log.Printf("[REGISTER ERROR TRACE]: %+v\n", err)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "data": nil})
		case errors.Is(err, ErrInvalidInputData):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": nil})
		}
		return
	}

	response := RegisterResponseDTO{
		ID:        output.Id,
		Name:      output.Name,
		Email:     output.Email,
		CreatedAt: output.CreatedAt,
	}

    
		c.JSON(http.StatusCreated, gin.H{
		"message": "We have Sent an verificaiton link to you Email",
		"data":    response,
	})

}