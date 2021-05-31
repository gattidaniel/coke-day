package main

import (
	"github.com/coke-day/pkg/validators"
	"net/http"
	"testing"

	"github.com/coke-day/functions/reservations/model"
	httpdelivery "github.com/coke-day/pkg/http"
	"github.com/stretchr/testify/assert"
)

type MockItemRepository struct{}

func (r *MockItemRepository) Delete(item *model.Reservation) error {
	return nil
}

func (r *MockItemRepository) Store(item *model.Reservation) error {
	return nil
}

func (r *MockItemRepository) Search(room, time, userDomain string) ([]model.Reservation, error) {
	return []model.Reservation{}, nil
}

func (r *MockItemRepository) Get(item *model.Reservation) (*model.ReservationPersistence, error) {
	return nil, nil
}

func TestCanFetchClient(t *testing.T) {
	request := httpdelivery.Req{
		PathParameters: map[string]string{"id": "123"},
		HTTPMethod:     "GET",
	}
	request.RequestContext.Authorizer = map[string]interface{}{
		"mail": "Bob@coke.com.us",
	}
	h := &Handler{&MockItemRepository{}, validators.CreateValidator()}
	router := createRoomRouting(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestCanCreateClient(t *testing.T) {
	request := httpdelivery.Req{
		Body:       `{ "room_name": "C01", "time": "19"}`,
		HTTPMethod: "POST",
	}
	request.RequestContext.Authorizer = map[string]interface{}{
		"mail": "Bob@coke.com.us",
	}
	h := &Handler{&MockItemRepository{}, validators.CreateValidator()}
	router := createRoomRouting(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestCanFetchClients(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "GET",
	}
	request.RequestContext.Authorizer = map[string]interface{}{
		"mail": "Bob",
	}
	h := &Handler{&MockItemRepository{}, validators.CreateValidator()}
	router := createRoomRouting(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandlerInvalidJSON(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "POST",
		Body:       "",
	}
	request.RequestContext.Authorizer = map[string]interface{}{
		"mail": "Bob",
	}
	h := &Handler{&MockItemRepository{}, validators.CreateValidator()}
	router := createRoomRouting(h)
	response, _ := router(request)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
