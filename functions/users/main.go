package main

import (
	"github.com/coke-day/functions/users/handlers"
	"github.com/coke-day/functions/users/repository"
	"github.com/coke-day/pkg/jwt"
	"github.com/coke-day/pkg/validators"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coke-day/pkg/datastore"
)

func main() {
	// Create a connection to the datastore, in this case, DynamoDB
	conn, err := datastore.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}

	// Create a new Dynamodb Table instance
	ddb := datastore.NewDynamoDB(conn, os.Getenv("DB_TABLE"))

	// Create a repository
	repository := repository.NewUserRepository(ddb)

	// Get Salt
	salt := os.Getenv("DB_TABLE")

	// Get Validator
	validator := validators.CreateValidator()

	// Get JWT Secret
	jwtSecret := os.Getenv("JWTSECRET")

	// Create JWT
	jwtHandler := jwt.NewJWT(jwtSecret, "u-should-hire-me", 3*time.Hour)

	// Create the handler instance, with the repository
	handler := handlers.CreateUserHandler(repository, salt, validator, jwtHandler)

	// Create Routing
	router := handler.CreateUserRouting()

	// Start the Lambda process
	lambda.Start(router)
}
