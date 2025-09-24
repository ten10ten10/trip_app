package main

import (
	"log"
	"os"

	"trip_app/api"
	"trip_app/internal/handler"
	"trip_app/internal/infrastructure/email"
	"trip_app/internal/repository"
	"trip_app/internal/security"
	"trip_app/internal/usecase"
	"trip_app/internal/validator"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
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

	// initialize repositories, usecases, and handlers
	userRepo := repository.NewUserRepository(db)
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

	userUsecase := usecase.NewUserUsecase(userRepo, userValidator, passwordGenerator, tokenGenerator, authTokenGenerator, emailSender)
	userHandler := handler.NewUserHandler(userUsecase)

	// start Echo server and register handlers
	e := echo.New()
	api.RegisterHandlers(e, userHandler)

	log.Println("Server starting on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
