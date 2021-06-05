package handlers

import (
	"github.com/coke-day/functions/users/model"
	"github.com/coke-day/pkg/criptography"
	"github.com/coke-day/pkg/jwt"
	"github.com/coke-day/pkg/validators"
	"net/http"
	"testing"
	"time"

	httpdelivery "github.com/coke-day/pkg/http"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct{}

func (r *MockUserRepository) Login(user model.User) (*model.User, error) {
	return &model.User{Email: "email@coke.com.us", HashPassword: []byte(criptography.HashPassword("123123", "thisIsJustForThisTest"))}, nil
}

func (r *MockUserRepository) Register(user model.User) error {
	return nil
}

func TestCanLogin(t *testing.T) {
	request := httpdelivery.Req{
		Body:       `{ "email": "email@coke.com.us", "password": "123123" }`,
		HTTPMethod: "POST",
		Path:       _loginPath,
	}
	h := createTestHandler()
	router := h.CreateUserRouting()
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func createTestHandler() Handler {
	return CreateUserHandler(&MockUserRepository{}, "thisIsJustForThisTest", validators.CreateValidator(), jwt.NewJWT("123", "issuer", 10*time.Hour))
}

func TestCanRegisterClient(t *testing.T) {
	request := httpdelivery.Req{
		Body:       `{ "name": "John Deere", "email":  "email@coke.com.us", "password": "test" }`,
		HTTPMethod: "POST",
		Path:       _registerPath,
	}
	h := createTestHandler()
	router := h.CreateUserRouting()
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestHandlerInvalidMethod(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "GET",
	}
	h := createTestHandler()
	router := h.CreateUserRouting()
	response, _ := router(request)
	assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
}

func TestHandlerInvalidJSON(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "POST",
		Body:       "",
		Path:       _registerPath,
	}
	h := createTestHandler()
	router := h.CreateUserRouting()
	response, _ := router(request)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
