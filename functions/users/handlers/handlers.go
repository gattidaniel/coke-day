package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/coke-day/functions/users/model"
	"github.com/coke-day/pkg/criptography"
	httpHelper "github.com/coke-day/pkg/http"
	"github.com/coke-day/pkg/jwt"
	"github.com/coke-day/pkg/validators"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var errInvalidLogin = errors.New("invalid login")

const (
	_registerPath = "/register"
	_loginPath    = "/login"
)

// UserRepository -
type UserRepository interface {
	Login(user model.User) (*model.User, error)
	Register(user model.User) error
}

// Handler -
type Handler struct {
	repository UserRepository
	salt       string
	validator  validator.Validate
	jwtHandler jwt.JWT
}

func CreateUserHandler(repository UserRepository, salt string, validator validator.Validate, jwtHandler jwt.JWT) Handler {
	return Handler{
		repository: repository,
		salt:       salt,
		validator:  validator,
		jwtHandler: jwtHandler,
	}
}

// Register store a user
func (h *Handler) Register(request httpHelper.Req) (httpHelper.Res, error) {
	var user *UserRegistration

	// Parsing body
	if err := httpHelper.ParseBody(request, &user); err != nil {
		return httpHelper.ErrResponse(err, http.StatusBadRequest)
	}

	// Validating body
	err := h.validator.Struct(user)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Register user
	if err := h.repository.Register(mapRegisterToUser(*user, h.salt)); err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}

	// Here we should send mail to valid email account. Of course, this isn't possible..

	return httpHelper.Response(map[string]bool{
		"success": true,
	}, http.StatusCreated)
}

// Login get the user
func (h *Handler) Login(request httpHelper.Req) (httpHelper.Res, error) {
	var userRequest *UserLogin

	// Parsing body
	if err := httpHelper.ParseBody(request, &userRequest); err != nil {
		return httpHelper.ErrResponse(err, http.StatusBadRequest)
	}

	// Validating body
	err := h.validator.Struct(userRequest)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Login user
	user, err := h.repository.Login(mapLoginToUser(*userRequest, h.salt))
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusNotFound)
	}

	// Validating user is found and password are equal.
	// If any of them is false, we return the same error for security reasons
	if user == nil || !bytes.Equal(user.HashPassword, hashPassword(userRequest.Password, h.salt)) {
		return httpHelper.ErrResponse(errInvalidLogin, http.StatusNotFound)
	}

	// Create JWT Token
	token, err := h.jwtHandler.CreateJWT(user.Email)
	if err != nil {
		return httpHelper.Res{}, fmt.Errorf("fail to create token: %v", err)
	}

	return httpHelper.Response(map[string]interface{}{
		"bearer-token": token,
	}, http.StatusOK)
}

func hashPassword(password, salt string) []byte {
	return criptography.HashPassword(password, salt)
}
