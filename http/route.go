package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"


	appointmentController "my-appointment-app/Appointment_System/Internals/Interface_Adapter/Controller"

	authController "my-appointment-app/Authentication_System/Internals/Interfaces_Adapter/controller"
)


type AuthorizingMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type Handlers struct {
	Register             *authController.RegisterUserHandler
	Verification         *authController.EmailVerificationGinHandler
	Login                *authController.LoginHandler
	AuthMiddleware       AuthorizingMiddleware
	DoctorInviation      *authController.AdminDoctorHandler
	DoctorRegisteration  *authController.DoctorRegistrationHandler
	CreateSpecialization *authController.SpecializationHandler
}


type ScheduleHandlers struct {
	Create   *appointmentController.ScheduleHandler
	Update   *appointmentController.UpdateScheduleHandler
	Cancel   *appointmentController.CancelScheduleHandler
	GetSlots *appointmentController.GetAvailableSlotsHandler
}

type AppointmentHandlers struct {
	Book         *appointmentController.BookAppointmentHandler
	Cancel       *appointmentController.CancelAppointmentHandler
	UpdateStatus *appointmentController.UpdateAppointmentStatusHandler
}


type AppHandlers struct {
	Auth         Handlers
	Schedule     ScheduleHandlers
	Appointment  AppointmentHandlers
}


func RegisterAuthRoutes(rg *gin.RouterGroup, h Handlers) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register.Register)
		auth.POST("/login", h.Login.Login)
		auth.GET("/verify-email", h.Verification.VerifyEmail)
		auth.POST("/Invitation", h.DoctorInviation.InviteDoctor)
		auth.POST("/accept-invite", h.DoctorRegisteration.RegisterDoctor)
		auth.POST("/create-specialization", h.CreateSpecialization.Create)
	}
}

func RegisterScheduleRoutes(rg *gin.RouterGroup, h ScheduleHandlers) {
	schedules := rg.Group("/schedules")
	{
		schedules.POST("/slots", h.GetSlots.GetAvailableSlots)
		schedules.POST("", h.Create.CreateSchedule)
		schedules.PUT("/:id", h.Update.UpdateSchedule)
		schedules.PATCH("/:id/cancel", h.Cancel.CancelSchedule)
	}
}

func RegisterAppointmentRoutes(rg *gin.RouterGroup, h AppointmentHandlers) {
	appointments := rg.Group("/appointments")
	{
		appointments.POST("", h.Book.BookAppointment)
		appointments.PATCH("/:id/cancel", h.Cancel.CancelAppointment)
		appointments.PATCH("/:id/status", h.UpdateStatus.UpdateAppointmentStatus)
	}
}


func SetupRouter(h AppHandlers) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
     }))

	api := r.Group("/api/v1")
	{
		RegisterAuthRoutes(api, h.Auth)
		RegisterScheduleRoutes(api, h.Schedule)
		RegisterAppointmentRoutes(api, h.Appointment)
	}

	return r
}