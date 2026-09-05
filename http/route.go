package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

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
	RefreshToken         *authController.RefreshTokenHandler
	Me                   *authController.MeHandler
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
	Auth        Handlers
	Schedule    ScheduleHandlers
	Appointment AppointmentHandlers
}

func RegisterAuthRoutes(rg *gin.RouterGroup, h Handlers) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register.Register)
		auth.POST("/login", h.Login.Login)
		auth.GET("/verify-email", h.Verification.VerifyEmail)
		auth.POST("/accept-invite", h.DoctorRegisteration.RegisterDoctor)
		auth.POST("/refresh", h.RefreshToken.RefreshToken)

		protected := auth.Group("")
		protected.Use(h.AuthMiddleware.Authenticate())
		{
			protected.GET("/me", h.Me.GetMe)
			protected.POST("/Invitation", h.DoctorInviation.InviteDoctor)
			protected.POST("/create-specialization", h.CreateSpecialization.Create)
		}
	}
}

func RegisterScheduleRoutes(rg *gin.RouterGroup, h ScheduleHandlers, authMW gin.HandlerFunc) {
	schedules := rg.Group("/schedules")
	{
		schedules.POST("/slots", h.GetSlots.GetAvailableSlots)

		protected := schedules.Group("")
		protected.Use(authMW)
		{
			protected.POST("", h.Create.CreateSchedule)
			protected.PUT("/:id", h.Update.UpdateSchedule)
			protected.PATCH("/:id/cancel", h.Cancel.CancelSchedule)
		}
	}
}

func RegisterAppointmentRoutes(rg *gin.RouterGroup, h AppointmentHandlers, authMW gin.HandlerFunc) {
	appointments := rg.Group("/appointments")
	appointments.Use(authMW) 
	{
		appointments.POST("", h.Book.BookAppointment)
		appointments.PATCH("/:id/cancel", h.Cancel.CancelAppointment)
		appointments.PATCH("/:id/status", h.UpdateStatus.UpdateAppointmentStatus)
	}
}

func SetupRouter(h AppHandlers) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authMW := h.Auth.AuthMiddleware.Authenticate()

	api := r.Group("/api/v1")
	{
		RegisterAuthRoutes(api, h.Auth)
		RegisterScheduleRoutes(api, h.Schedule, authMW)
		RegisterAppointmentRoutes(api, h.Appointment, authMW)
	}

	return r
}