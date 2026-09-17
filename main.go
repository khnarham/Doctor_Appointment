package main

import (
	"log"
	appointmentDatabase "my-appointment-app/Appointment_System/Internals/Infrastructure/Database"
	appointmentController "my-appointment-app/Appointment_System/Internals/Interface_Adapter/Controller"
	appointmentUsecase "my-appointment-app/Appointment_System/Internals/Usecase"
	Infra "my-appointment-app/Authentication_System/Internals/Infrastructure/Crypto"
	database "my-appointment-app/Authentication_System/Internals/Infrastructure/Database"
	jwt "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"
	"my-appointment-app/Authentication_System/Internals/Infrastructure/Mailer"
	"my-appointment-app/Authentication_System/Internals/Interfaces_Adapter/controller"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	http "my-appointment-app/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)
func main() {
	// DB Initilized //
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	// env Initlizzation //
    env_err := godotenv.Load()
    if env_err != nil {
        log.Println("Warning: .env file not found")
    }
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
    log.Fatal("BREVO_API_KEY is missing or empty!")
}
	senderEmail := "khnarham5678@gmail.com"

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
	specializationRepo := database.NewSpecializationRepository(db)
	receptionstInviteRepo := database.NewPostgresReceptionistInvitationRepository(db)
	receptionistAcceptRepo := database.NewPostgresUserRepository(db)
	receptionstUserRepo:= database.NewPostgreReceptionstRepo(db)
	getDoctorRepo  := database.NewGetDoctorRepository(db)
	specializationOperationRepo := database.NewSpecializationOperationsRepository(db)
	getDoctorsRepo := database.NewGetDoctorsRepository(db)
	getPatientAppointmentRepo := appointmentDatabase.NewPatientAppointmentRepository(db)
	reviewAppointmentRepo := appointmentDatabase.NewReviewRepository(db)
	getDoctorAppointmentRepo := appointmentDatabase.NewDoctorAppointmentRepository(db)
    getAnalyticRepo := appointmentDatabase.NewPostgresAnalyticsRepository(db)
	getDoctorAnalyticRepo := appointmentDatabase.NewDoctorAnalyticsRepository(db)
    // Infrastrcuture
	hasher := Infra.NewBcryptHasher()
	tokenGen := Infra.NewCryptoCodeGenerator()
	JWTService := jwt.NewJWTService(accessSecret , refreshSecret , accessExpiry , refreshExpiry)
	emailSender:= mailer.NewBrevoMailer(apiKey ,senderEmail)
	authMiddleware := http.NewAuthMiddleware(JWTService)
	doctoremailSender := mailer.NewBrevoDoctorVerification(apiKey , senderEmail)
	receptionstemailSender := mailer.NewBrevoReceptionistMailer(apiKey , senderEmail , "")
     
timeoutContext := 5 * time.Second

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
	doctoremailSender,
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
	getMeUsecase := usecase.NewGetMeUsecase(userRepo)
	listSpecializationUsecase := usecase.NewListSpecializationsUsecase(specializationRepo)

	receptionstInviteUsecase := usecase.NewInviteReceptionistUsecase(receptionstInviteRepo , receptionstUserRepo , receptionstemailSender , tokenGen)
	acceptingInviteUsecase := usecase.NewAcceptInvitationUsecase(receptionstInviteRepo , receptionistAcceptRepo , hasher)
	getDoctorUsecase := usecase.NewDoctorUsecase(getDoctorRepo)
	specializationOperationUsecase := usecase.NewSpecializationUsecase(specializationOperationRepo)
	getDoctorSchedulesUsecase := appointmentUsecase.NewGetDoctorSchedulesUseCase(scheduleRepo , appontmentRepo)
	getDoctorsUsecase := usecase.NewGetDoctorUsecase(getDoctorsRepo)
	getPatientUsecase := usecase.NewGetPatientProfileUseCase(patientRepo)
	getPatientAppointmentUsecase := appointmentUsecase.NewGetPatientAppointmentsUseCase(getPatientAppointmentRepo)
	reviewAppointmentUsecase := appointmentUsecase.NewSubmitReviewUsecase(bookappointmentRepo ,reviewAppointmentRepo)
	getDoctorAppointmentUsecase := appointmentUsecase.NewGetDoctorAppointmentsUseCase(getDoctorAppointmentRepo)
	updatedAppointmentUsecase := appointmentUsecase.NewUpdateAppointmentStatusUseCase(bookappointmentRepo)
	registerWalkInUsecase := usecase.NewRegisterWalkInPatientUsecase(patientRepo)
	GetAnalyticUsecase := appointmentUsecase.NewAnalyticsUsecase(getAnalyticRepo)
	GetDoctorAnalyticsUsecase := appointmentUsecase.NewDoctorAnalyticsUsecase(getDoctorAnalyticRepo , timeoutContext)
	 
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
	getMeHandler := controller.NewMeHandler(getMeUsecase)
	logoutHandler := controller.NewLogoutHandler()
	listSpecializationHandler := controller.NewListSpecializationHandler(listSpecializationUsecase)
	receptionistInviteHandler := controller.NewReceptionistHandler(receptionstInviteUsecase)
	acceptingInviteHandler := controller.NewAcceptInvitationHandler(acceptingInviteUsecase)
	getDoctorHandler := controller.NewGetDoctorHandler(getDoctorUsecase)
	specializationOperationHandler := controller.NewSpecializationOperationHandler(specializationOperationUsecase)
	getDoctorScheduleHandler := appointmentController.NewGetDoctorSchedulesHandler(getDoctorSchedulesUsecase)
	getDoctorsHandler := controller.NewDoctorHandler(*getDoctorsUsecase)
	getPatientHandler := controller.NewGetPatientProfileHandler(getPatientUsecase)
	getPatientAppointmentHandler := appointmentController.NewGetPatientAppointmentsHandler(getPatientAppointmentUsecase)
	reviewAppointmentHandler := appointmentController.NewSubmitReviewHandler(reviewAppointmentUsecase)
	getDoctorAppointmentHandler := appointmentController.NewGetDoctorAppointmentsHandler(getDoctorAppointmentUsecase)
	updateAppointmentStatusHandler := appointmentController.NewUpdateAppointmentStatusHandler(updatedAppointmentUsecase)
	registerWalkInPatientHandler := controller.NewRegisterWalkInPatientHandler(registerWalkInUsecase)
    getAnalyticHandler := appointmentController.NewAnalyticsHandler(GetAnalyticUsecase)
	getDoctorAnalyticsHandler := appointmentController.NewDoctorAnalyticsHandler(GetDoctorAnalyticsUsecase)


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
		Me: getMeHandler,
		Logout: logoutHandler,
		ListSpecialization: listSpecializationHandler,
		ReceptionistAccepting: acceptingInviteHandler,
		ReceptionistInvitation: receptionistInviteHandler,
		GetDoctor: getDoctorHandler,
        SpecilizationOperation: specializationOperationHandler,
		GetDoctors: getDoctorsHandler,
		GetPatient: getPatientHandler,
		RegisterWalkinPatient: registerWalkInPatientHandler,
	}
    
	// Appointment Handler // 
	AppointmentHandler := http.ScheduleHandlers{
         Create: createScheduleHandler,
		 GetSlots: getAvailableSlotHandler,
		 Update: updateScheduleHandler,
		 GetDoctorSchedules: getDoctorScheduleHandler,	
	}

	BookingAppointmentHandler := http.AppointmentHandlers{
		Book: bookAppointmentHandler,
		GetPatientAppointment: getPatientAppointmentHandler,
		ReviewAppointment: reviewAppointmentHandler,
		GetDoctorAppointment: getDoctorAppointmentHandler,
		UpdateStatus: updateAppointmentStatusHandler,
		GetAnalytics: getAnalyticHandler,
		GetDoctorAnalytics: getDoctorAnalyticsHandler,
	}

    
	// Mian Handler //
	handlers := http.AppHandlers{
		Auth: Authenticationhandlers,
		Schedule: AppointmentHandler,
		Appointment: BookingAppointmentHandler,
	}
	r := http.SetupRouter(handlers)
    _ = r.SetTrustedProxies(nil)
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