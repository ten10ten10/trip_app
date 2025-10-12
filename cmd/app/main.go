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
	"trip_app/internal/validator"

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

	// initialize services
	userValidator := validator.NewUserValidator()
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

	// schedule validators
	scheduleHandlerValidator := handler.NewScheduleHandlerValidator()
	scheduleUsecaseValidator := usecase.NewScheduleUsecaseValidator()

	// initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, userValidator, passwordGenerator, tokenGenerator, authTokenGenerator, emailSender)
	tripUsecase := usecase.NewTripUsecase(tripRepo)
	scheduleUsecase := usecase.NewScheduleUsecase(scheduleRepo, scheduleUsecaseValidator)

	// initialize the composite handler
	h := handler.NewHandler(userUsecase, tripUsecase, scheduleUsecase, scheduleHandlerValidator)

	// initialize middlewares
	tripOwnershipMiddleware := middleware.TripOwnershipMiddleware(tripUsecase)
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// start Echo server
	e := echo.New()

	// Create a wrapper for manual route registration
	wrapper := &api.ServerInterfaceWrapper{Handler: h}

	// Public routes
	e.POST("/login", wrapper.LoginUser)
	e.POST("/signup", wrapper.CreateUser)
	e.POST("/users/verify/:verificationToken", wrapper.VerifyUser)
	e.GET("/public/trips/:shareToken", wrapper.GetPublicTripByShareToken)
	e.PUT("/public/trips/:shareToken", wrapper.UpdatePublicTripByShareToken)
	e.GET("/public/trips/:shareToken/details", wrapper.GetTripDetailsForPublicTrip)
	e.GET("/public/trips/:shareToken/schedules", wrapper.GetSchedulesForPublicTrip)
	e.POST("/public/trips/:shareToken/schedules", wrapper.AddScheduleToPublicTrip)
	e.GET("/public/trips/:shareToken/schedules/:scheduleId", wrapper.GetScheduleForPublicTrip)
	e.PATCH("/public/trips/:shareToken/schedules/:scheduleId", wrapper.UpdateScheduleForPublicTrip)
	e.DELETE("/public/trips/:shareToken/schedules/:scheduleId", wrapper.DeleteScheduleForPublicTrip)

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
	tripOwnerGroup.GET("/share", wrapper.GetShareLinkForTrip)
	tripOwnerGroup.POST("/share", wrapper.CreateShareLinkForTrip)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
