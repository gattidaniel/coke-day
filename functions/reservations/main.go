package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coke-day/functions/reservations/handlers"
	"github.com/coke-day/functions/reservations/repository"
	"github.com/coke-day/pkg/validators"
	"log"
	"os"

	"github.com/coke-day/pkg/datastore"
)

func main() {
	conn, err := datastore.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}

	ddb := datastore.NewDynamoDB(conn, os.Getenv("DB_TABLE"))
	repository := repository.NewItemRepository(ddb)

	// Get Validator
	v := validators.CreateValidator()

	handler := handlers.NewHandler(repository, v)
	router := handler.CreateRoomRouting()
	lambda.Start(router)
}
