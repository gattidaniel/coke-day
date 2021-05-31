package main

import (
	"bytes"
	"errors"
	"github.com/coke-day/pkg/criptography"
	"github.com/coke-day/pkg/validators"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coke-day/functions/users/model"
	"github.com/coke-day/pkg/datastore"
	httpHelper "github.com/coke-day/pkg/http"
)

// UserRepository -
type UserRepository interface {
	Login(user *model.UserLogin) (*model.UserPersistence, error)
	Register(user *model.UserRegistration) error
}

// Handler -
type Handler struct {
	repository UserRepository
	salt       string
	validator  validator.Validate
}

const (
	_registerPath = "/register"
	_loginPath    = "/login"
)

var errInvalidLogin = errors.New("invalid login")

// Register store a user
func (h *Handler) Register(request httpHelper.Req) (httpHelper.Res, error) {
	var user *model.UserRegistration

	// Parsing body
	if err := httpHelper.ParseBody(request, &user); err != nil {
		return httpHelper.ErrResponse(err, http.StatusBadRequest)
	}

	// Validating body
	err := h.validator.Struct(user)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Hashing password
	user.HashPassword = criptography.HashPassword(user.Password, h.salt)

	// Register user
	if err := h.repository.Register(user); err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}

	// Here we should send mail to valid email account. Of course, this isn't possible..

	return httpHelper.Response(map[string]bool{
		"success": true,
	}, http.StatusCreated)
}

// Login get the user
func (h *Handler) Login(request httpHelper.Req) (httpHelper.Res, error) {
	var userLogin *model.UserLogin

	// Parsing body
	if err := httpHelper.ParseBody(request, &userLogin); err != nil {
		return httpHelper.ErrResponse(err, http.StatusBadRequest)
	}

	// Validating body
	err := h.validator.Struct(userLogin)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Hashing password
	hashPassword := criptography.HashPassword(userLogin.Password, h.salt)

	// Register user
	user, err := h.repository.Login(userLogin)
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusNotFound)
	}

	// Validating user is found and password are equal.
	// If any of them is false, we return the same error for security reasons
	if user == nil || !bytes.Equal(user.HashPassword, hashPassword) {
		return httpHelper.ErrResponse(errInvalidLogin, http.StatusNotFound)
	}

	// TODO: Return Token ID with time expiration
	return httpHelper.Response(map[string]interface{}{
		"success": true,
	}, http.StatusOK)
}

func main() {
	// Create a connection to the datastore, in this case, DynamoDB
	conn, err := datastore.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}

	// Create a new Dynamodb Table instance
	ddb := datastore.NewDynamoDB(conn, os.Getenv("DB_TABLE"))

	// Create a repository
	repository := model.NewClientRepository(ddb)

	// Get Salt
	salt := os.Getenv("DB_TABLE")

	// Get Validator
	v := validators.CreateValidator()

	// Create the handler instance, with the repository
	handler := &Handler{repository, salt, v}

	router := createUserRouting(handler)

	// Start the Lambda process
	lambda.Start(router)
}

// createUserRouting routes restful endpoints to the correct method
// POST with Path '/register' calls the Register method
// POST with Path '/login' calls the Login method
func createUserRouting(handler *Handler) func(request httpHelper.Req) (httpHelper.Res, error) {
	router := func(request httpHelper.Req) (httpHelper.Res, error) {
		if request.HTTPMethod == "POST" && strings.ToLower(request.Path) == _registerPath {
			return handler.Register(request)
		}
		if request.HTTPMethod == "POST" && strings.ToLower(request.Path) == _loginPath {
			return handler.Login(request)
		}
		return httpHelper.ErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
	return router
}
