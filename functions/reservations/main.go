package main

import (
	"errors"
	"github.com/coke-day/functions/reservations/model"
	"github.com/coke-day/pkg/utils"
	"github.com/coke-day/pkg/validators"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/coke-day/pkg/datastore"
	httpHelper "github.com/coke-day/pkg/http"
)

// ItemRepository -
type ItemRepository interface {
	Delete(item *model.Reservation) error
	Store(item *model.Reservation) error
	Search(room, time, userDomain string) ([]model.Reservation, error)
	Get(item *model.Reservation) (*model.ReservationPersistence, error)
}

// Handler -
type Handler struct {
	repository ItemRepository
	validator  validator.Validate
}

var errInvalidRoom = errors.New("invalid Room")
var errForbidden = errors.New("forbidden")

// Search all items
func (h *Handler) Search(request httpHelper.Req) (httpHelper.Res, error) {
	// Getting parameters
	room := request.QueryStringParameters["room"]
	time := request.QueryStringParameters["time"]
	mail := request.RequestContext.Authorizer["email"].(string)

	// Searching
	items, err := h.repository.Search(room, time, mail)
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusNotFound)
	}

	return httpHelper.Response(map[string]interface{}{
		"items": items,
	}, http.StatusOK)
}

// Store a new item
func (h *Handler) Store(request httpHelper.Req) (httpHelper.Res, error) {
	var reservation model.Reservation
	reservation.UserEmail = request.RequestContext.Authorizer["email"].(string)
	if err := httpHelper.ParseBody(request, &reservation); err != nil {
		return httpHelper.ErrResponse(err, http.StatusBadRequest)
	}

	// Validating body
	err := h.validator.Struct(reservation)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Validating email domain and room name
	reservation.RoomName = strings.ToUpper(reservation.RoomName)
	domain, err := utils.GetDomain(reservation.UserEmail)
	if err != nil {
		return httpHelper.ErrResponse(errInvalidRoom, http.StatusInternalServerError)
	}
	if !(strings.Contains(domain, validators.CokeDomain) && strings.HasPrefix(reservation.RoomName, "C") ||
		strings.Contains(domain, validators.PepsiDomain) && strings.HasPrefix(reservation.RoomName, "P")) {
		return httpHelper.ErrResponse(errInvalidRoom, http.StatusForbidden)
	}

	// Validate that time is empty or the same user
	storedReservation, err := h.repository.Get(&reservation)
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}
	if storedReservation != nil && storedReservation.UserEmail != "" && storedReservation.UserEmail != reservation.UserEmail {
		return httpHelper.ErrResponse(errForbidden, http.StatusForbidden)
	}

	// Storing
	if err := h.repository.Store(&reservation); err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}

	return httpHelper.Response(map[string]bool{
		"created": true,
	}, http.StatusCreated)
}

// Delete a new item
func (h *Handler) Delete(request httpHelper.Req) (httpHelper.Res, error) {
	room := request.PathParameters["room"]
	time := request.PathParameters["time"]
	reservation := model.Reservation{
		RoomName:  room,
		Time:      time,
		UserEmail: request.RequestContext.Authorizer["email"].(string),
	}

	// Validating body
	err := h.validator.Struct(reservation)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Validate user is the owner
	storedReservation, err := h.repository.Get(&reservation)
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}
	if storedReservation == nil || storedReservation.UserEmail != reservation.UserEmail {
		return httpHelper.ErrResponse(errForbidden, http.StatusForbidden)
	}

	// Deleting
	if err := h.repository.Delete(&reservation); err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}

	return httpHelper.Response(map[string]bool{
		"delete": true,
	}, http.StatusOK)
}

func main() {
	conn, err := datastore.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}

	ddb := datastore.NewDynamoDB(conn, os.Getenv("DB_TABLE"))
	repository := model.NewItemRepository(ddb)

	// Get Validator
	v := validators.CreateValidator()

	handler := &Handler{repository, v}
	router := createRoomRouting(handler)
	lambda.Start(router)
}

// createRoomRouting routes restful endpoints to the correct method
// GET calls the Search method (may receive by parameter room or time and filter by this)
// POST calls the Schedule method.
// DELETE calls the Delete method.
func createRoomRouting(handler *Handler) func(request httpHelper.Req) (httpHelper.Res, error) {
	router := func(request httpHelper.Req) (httpHelper.Res, error) {
		switch request.HTTPMethod {
		case "GET":
			return handler.Search(request)
		case "DELETE":
			return handler.Delete(request)
		case "POST":
			return handler.Store(request)
		default:
			return httpHelper.ErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}
	}
	return router
}
