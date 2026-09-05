package main

import (
	"log"
	Infra "my-appointment-app/Authentication_System/Internals/Infrastructure/Crypto"
	database "my-appointment-app/Authentication_System/Internals/Infrastructure/Database"
	jwt "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"
	"my-appointment-app/Authentication_System/Internals/Infrastructure/Mailer"
	"my-appointment-app/Authentication_System/Internals/Interfaces_Adapter/controller"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	appointmentController "my-appointment-app/Appointment_System/Internals/Interface_Adapter/Controller"
	appointmentDatabase "my-appointment-app/Appointment_System/Internals/Infrastructure/Database" 
	appointmentUsecase "my-appointment-app/Appointment_System/Internals/Usecase"
	http "my-appointment-app/http"
	"os"
	"time"
)

func main() {
	// DB Initilized //
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	// env Initlizzation //
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		log.Fatal("RESEND_API_KEY environment variable is missing")
	}

	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "onboarding@resend.dev" 
	}
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessExpiry, _ := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRY"))
	refreshExpiry, _ := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRY"))

	txManager := database.NewPgxTransactionManager(db)
	userRepo := database.NewUserRegisterRepository(db)
	patientRepo := database.NewPatientRegisterRepository(db)
	verRepo := database.NewEmailSaveRepository(db)
	emVerRepo := database.NewEmailVerificationRepository(db)
	existedRepo := database.NewExistingUserRegistrationRepo(db)
	exitedPatientRepo := database.NewPatientUpdateRepository(db)
	doctorInviteRepo := database.NewDoctorInvitationRepository(db)
	sepecilizationRepo := database.NewSpecializationRepository(db)
	registerDoctorRepo := database.NewDoctorRepository(db)
	scheduleRepo := appointmentDatabase.NewScheduleRepository(db)
	appontmentRepo := appointmentDatabase.NewAppointmentRepository(db)
	bookappointmentRepo := appointmentDatabase.NewPgxBookAppointmentRepository(db)
     
    // Infrastrcuture
	hasher := Infra.NewBcryptHasher()
	tokenGen := Infra.NewCryptoCodeGenerator()
	emailSender := Mailer.NewResendMailer(resendAPIKey, senderEmail)
	JWTService := jwt.NewJWTService(accessSecret , refreshSecret , accessExpiry , refreshExpiry)
	authMiddleware := http.NewAuthMiddleware(JWTService)
     
	// Usecase // 
	registerUsecase := usecase.NewRegisterUserUsecase(
     userRepo,
	 hasher,
	 patientRepo,
	 verRepo,
	 tokenGen,
	 txManager,
	 emailSender,
	 existedRepo,
	 exitedPatientRepo,
	)
	
	verifyEmailUsecase := usecase.NewEmailVerificationUsecase(
		emVerRepo,
		JWTService,
	)
	loginUserUsecase := usecase.NewLoginUsecase(
		userRepo,
		hasher, 
		JWTService,
	)
	specilizationUescase := usecase.NewCreateSpecializationUsecase(sepecilizationRepo)

    doctorInvitationUsecase := usecase.NewInviteDoctorUsecase(
	doctorInviteRepo,
	userRepo,
	sepecilizationRepo,
	emailSender,
    tokenGen,
    )
	acceptInviteUsecase := usecase.NewRegisterDoctorWithInviteUsecase(
		doctorInviteRepo, 
		userRepo,
		registerDoctorRepo,
		hasher,
		txManager,
	)

	createScheduleUsecase := appointmentUsecase.NewCreateScheduleUsecase(
      scheduleRepo,
	)
     
    getAvailableSlotUsecase := appointmentUsecase.NewGetAvailableSlotsUseCase(
     scheduleRepo,
	 appontmentRepo,
	)

	updateScheduleUsecase := appointmentUsecase.NewUpdateScheduleUsecase(
      scheduleRepo,
	)

	bookAppointmentUsecase := appointmentUsecase.NewBookAppointmentUseCase(
      scheduleRepo,
	  bookappointmentRepo,
	)
	RefreshTokenUsecase := usecase.NewRefreshTokenUsecase(
		JWTService,
	)
	 
    // Handlers // 
	registerHandler := controller.NewRegisterUserHandler(registerUsecase)
	verificationHandler := controller.NewEmailVerificationGinHandler(verifyEmailUsecase)
	loginHandler := controller.NewLoginHandler(loginUserUsecase)
	specilizationHandler := controller.NewSpecializationHandler(specilizationUescase)
	doctorInvitationHandler := controller.NewAdminDoctorHandler(doctorInvitationUsecase)
	doctorRegisterationHandler := controller.NewDoctorRegistrationHandler(acceptInviteUsecase)
	createScheduleHandler := appointmentController.NewScheduleHandler(createScheduleUsecase)
	getAvailableSlotHandler := appointmentController.NewGetAvailableSlotsHandler(getAvailableSlotUsecase)
	updateScheduleHandler := appointmentController.NewUpdateScheduleHandler(updateScheduleUsecase)
    bookAppointmentHandler := appointmentController.NewBookAppointmentHandler(bookAppointmentUsecase)
	RefreshTokenHandler := controller.NewRefreshTokenHandler(RefreshTokenUsecase)



	// Authentication Handlers //
	Authenticationhandlers := http.Handlers{
		Register:     registerHandler,
		Verification: verificationHandler,
		Login: loginHandler,
		AuthMiddleware: authMiddleware,
		CreateSpecialization: specilizationHandler,
		DoctorInviation: doctorInvitationHandler,
		DoctorRegisteration: doctorRegisterationHandler,
		RefreshToken: RefreshTokenHandler,

	}
    
	// Appointment Handler // 
	AppointmentHandler := http.ScheduleHandlers{
         Create: createScheduleHandler,
		 GetSlots: getAvailableSlotHandler,
		 Update: updateScheduleHandler,
	}

	BookingAppointmentHandler := http.AppointmentHandlers{
		Book: bookAppointmentHandler,
	}

    
	// Mian Handler //
	handlers := http.AppHandlers{
		Auth: Authenticationhandlers,
		Schedule: AppointmentHandler,
		Appointment: BookingAppointmentHandler,
	}
	r := http.SetupRouter(handlers)
    
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
    
	// Server Start //
	log.Printf("Server running on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}