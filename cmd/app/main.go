package main

import (
	"log"
	"os"

	"trip_app/api"
	"trip_app/internal/handler"
	"trip_app/internal/infrastructure/email"
	"trip_app/internal/middleware"
	"trip_app/internal/repository"
	"trip_app/internal/security"
	"trip_app/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	// connect to the database
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// get jwt secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// initialize repositories
	userRepo := repository.NewUserRepository(db)
	tripRepo := repository.NewTripRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	shareTokenRepo := repository.NewShareTokenRepository(db)
	publicTripRepo := repository.NewPublicTripRepository(db)

	// initialize services
	passwordGenerator := security.NewPasswordGenerator()
	tokenGenerator := security.NewTokenGenerator()
	authTokenGenerator := security.NewAuthTokenGenerator(jwtSecret)
	emailSender, err := email.NewEmailSender(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("EMAIL_FROM"),
	)
	if err != nil {
		log.Fatalf("failed to create email sender: %v", err)
	}

	// initialize validators
	userHandlerValidator := handler.NewUserHandlerValidator()
	userUsecaseValidator := usecase.NewUserUsecaseValidator()
	scheduleHandlerValidator := handler.NewScheduleHandlerValidator()
	scheduleUsecaseValidator := usecase.NewScheduleUsecaseValidator()

	// initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, userUsecaseValidator, passwordGenerator, tokenGenerator, authTokenGenerator, emailSender)
	tripUsecase := usecase.NewTripUsecase(tripRepo, tokenGenerator)
	scheduleUsecase := usecase.NewScheduleUsecase(scheduleRepo, scheduleUsecaseValidator)
	shareTokenUsecase := usecase.NewShareTokenUsecase(shareTokenRepo, tokenGenerator)
	publicTripUsecase := usecase.NewPublicTripUsecase(publicTripRepo, tokenGenerator)

	// initialize the composite handler
	h := handler.NewHandler(userUsecase, tripUsecase, scheduleUsecase, shareTokenUsecase, publicTripUsecase, userHandlerValidator, scheduleHandlerValidator)

	// initialize middlewares
	tripOwnershipMiddleware := middleware.TripOwnershipMiddleware(tripUsecase)
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	shareTokenOwnershipMiddleware := middleware.ShareTokenOwnershipMiddleware(publicTripUsecase)

	// start Echo server
	e := echo.New()

	// Create a wrapper for manual route registration
	wrapper := &api.ServerInterfaceWrapper{Handler: h}

	// Public routes (no authentication)
	e.POST("/login", wrapper.LoginUser)
	e.POST("/signup", wrapper.CreateUser)
	e.POST("/users/verify/:verificationToken", wrapper.VerifyUser)

	// Public trip routes (with share token validation)
	publicTripGroup := e.Group("/public/trips/:shareToken")
	publicTripGroup.Use(shareTokenOwnershipMiddleware)
	publicTripGroup.GET("", wrapper.GetPublicTripByShareToken)
	publicTripGroup.PUT("", wrapper.UpdatePublicTripByShareToken)
	publicTripGroup.GET("/details", wrapper.GetTripDetailsForPublicTrip)
	publicTripGroup.GET("/schedules", wrapper.GetSchedulesForPublicTrip)
	publicTripGroup.POST("/schedules", wrapper.AddScheduleToPublicTrip)
	publicTripGroup.GET("/schedules/:scheduleId", wrapper.GetScheduleForPublicTrip)
	publicTripGroup.PATCH("/schedules/:scheduleId", wrapper.UpdateScheduleForPublicTrip)
	publicTripGroup.DELETE("/schedules/:scheduleId", wrapper.DeleteScheduleForPublicTrip)

	// Auth-required routes
	authRequired := e.Group("")
	authRequired.Use(authMiddleware)
	authRequired.POST("/logout", wrapper.LogoutUser)
	authRequired.GET("/me", wrapper.GetMe)
	authRequired.PUT("/me/password", wrapper.ChangePassword)
	authRequired.GET("/trips", wrapper.GetUserTrips)
	authRequired.POST("/trips", wrapper.CreateUserTrip)

	// Trip ownership-required routes
	tripOwnerGroup := authRequired.Group("/trips/:tripId")
	tripOwnerGroup.Use(tripOwnershipMiddleware)
	tripOwnerGroup.GET("", wrapper.GetUserTrip)
	tripOwnerGroup.PUT("", wrapper.UpdateUserTrip)
	tripOwnerGroup.DELETE("", wrapper.DeleteUserTrip)
	tripOwnerGroup.GET("/details", wrapper.GetTripDetails)
	tripOwnerGroup.GET("/schedules", wrapper.GetSchedulesForTrip)
	tripOwnerGroup.POST("/schedules", wrapper.AddScheduleToTrip)
	tripOwnerGroup.GET("/schedules/:scheduleId", wrapper.GetScheduleForTrip)
	tripOwnerGroup.PATCH("/schedules/:scheduleId", wrapper.UpdateScheduleForTrip)
	tripOwnerGroup.DELETE("/schedules/:scheduleId", wrapper.DeleteScheduleForTrip)
	tripOwnerGroup.POST("/share", wrapper.CreateShareLinkForTrip)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
