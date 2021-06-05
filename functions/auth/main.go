package main

import (
	"github.com/coke-day/functions/auth/handlers"
	"github.com/coke-day/pkg/jwt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Get JWT Secret
	jwtSecret := os.Getenv("JWTSECRET")

	// Create JWT
	jwtHandler := jwt.NewJWT(jwtSecret, "u-should-hire-me", 3*time.Hour)

	// Create the Auth instance
	handler := handlers.NewHandler(jwtHandler)

	// Create the router
	router := handler.CreateAuthRouting()

	lambda.Start(router)
}
