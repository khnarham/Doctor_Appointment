package http

import (
	"os"
	"strings"

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
	AdminAuth              *authController.AdminAuthHandler
	AuthMiddleware       AuthorizingMiddleware
	DoctorInviation      *authController.AdminDoctorHandler
	DoctorRegisteration  *authController.DoctorRegistrationHandler
	CreateSpecialization *authController.SpecializationHandler
	ListSpecialization   *authController.ListSpecializationHandler
	Logout               *authController.LogoutHandler
	RefreshToken         *authController.RefreshTokenHandler
	Me                   *authController.MeHandler
	ReceptionistInvitation *authController.ReceptionistHandler
	ReceptionistAccepting *authController.AcceptInvitationHandler 
	GetDoctor              *authController.GetDoctorHandler
	SpecilizationOperation  *authController.SpecializationOperationHandler
	GetDoctors               *authController.GetDoctorsHandler
	GetPatient              *authController.GetPatientProfileHandler
	RegisterWalkinPatient  *authController.RegisterWalkInPatientHandler

}

type ScheduleHandlers struct {
	Create   *appointmentController.ScheduleHandler
	Update   *appointmentController.UpdateScheduleHandler
	Cancel   *appointmentController.CancelScheduleHandler
	GetSlots *appointmentController.GetAvailableSlotsHandler
	GetDoctorSchedules *appointmentController.GetDoctorSchedulesHandler
}

type AppointmentHandlers struct {
	Book         *appointmentController.BookAppointmentHandler
	Cancel       *appointmentController.CancelAppointmentHandler
	UpdateStatus *appointmentController.UpdateAppointmentStatusHandler
	GetPatientAppointment *appointmentController.GetPatientAppointmentsHandler
	ReviewAppointment *appointmentController.SubmitReviewHandler
	GetDoctorAppointment *appointmentController.GetDoctorAppointmentsHandler
	GetAnalytics  *appointmentController.AnalyticsHandler
	GetDoctorAnalytics *appointmentController.DoctorAnalyticsHandler
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
		auth.POST("/logout",h.Logout.Logout)
		auth.POST("/Receptionist-invite", h.ReceptionistInvitation.InviteReceptionist)
		auth.POST("/Receptionist-accept", h.ReceptionistAccepting.AcceptInvitation)
		auth.POST("/admin/login",h.AdminAuth.Login)

		auth.GET("/list-specialization" , h.ListSpecialization.ListSpecializations)
		auth.POST("/create-specialization", h.CreateSpecialization.Create)
		auth.POST("/Invitation", h.DoctorInviation.InviteDoctor)
		auth.PUT("/specialization/update/:id", h.SpecilizationOperation.Update)
		auth.DELETE("/specialization/delete/:id", h.SpecilizationOperation.Delete)
		auth.GET("/doctors" , h.GetDoctors.GetDoctors)
		auth.POST("/receptionist/patients/walk-in", h.RegisterWalkinPatient.Handle)
		protected := auth.Group("")
		protected.Use(h.AuthMiddleware.Authenticate())
		{ 
			protected.GET("/doctor" , h.GetDoctor.GetDoctorProfile)
			protected.GET("/patients/me" , h.GetPatient.GetProfile)
			protected.GET("/me", h.Me.GetMe)
		}
	}
}

func RegisterScheduleRoutes(rg *gin.RouterGroup, h ScheduleHandlers, authMW gin.HandlerFunc) {
	schedules := rg.Group("/schedules")
	{
		schedules.POST("/slots", h.GetSlots.GetAvailableSlots)

		protected := schedules.Group("")
		protected.GET("/doctor/:id" ,h.GetDoctorSchedules.GetDoctorSchedulesWithSlots)
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
	appointments.GET("doctors/:id/appointments", h.GetDoctorAppointment.GetAppointments)
	appointments.GET("/analytics/financial-overview" , h.GetAnalytics.GetFinancialOverview)
	appointments.GET("/doctors/dashboard/analytics/:id" , h.GetDoctorAnalytics.GetDashboard)
	appointments.Use(authMW) 
	{
		appointments.POST("/reviews" , h.ReviewAppointment.SubmitReview)
		appointments.GET("/my-appointment/:id" , h.GetPatientAppointment.GetAppointments)
		appointments.POST("", h.Book.BookAppointment)
		appointments.PATCH("/:id/cancel", h.Cancel.CancelAppointment)
		appointments.PATCH("/:id/status", h.UpdateStatus.UpdateAppointmentStatus)
	}
}

func SetupRouter(h AppHandlers) *gin.Engine {
	r := gin.Default()

allowOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowOriginsEnv != "" {
		origins = strings.Split(allowOriginsEnv, ",")
	} else {
		origins = []string{"http://localhost:3000"}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
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