package handlers

import (
	"errors"
	"github.com/coke-day/functions/reservations/model"
	httpHelper "github.com/coke-day/pkg/http"
	"github.com/coke-day/pkg/utils"
	"github.com/coke-day/pkg/validators"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
)

var errInvalidRoom = errors.New("invalid Room")
var errForbidden = errors.New("forbidden")

// Handler -
type Handler struct {
	repository ItemRepository
	validator  validator.Validate
}

func NewHandler(repository ItemRepository, validator validator.Validate) Handler {
	return Handler{
		repository: repository,
		validator:  validator,
	}
}

// ItemRepository -
type ItemRepository interface {
	Delete(item model.Reservation) error
	Store(item model.Reservation) error
	Search(room, time, userDomain string) ([]model.Reservation, error)
	Get(item model.Reservation) (*model.Reservation, error)
}

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

	var reservationsItem []Reservation
	for _, item := range items {
		reservationsItem = append(reservationsItem, mapFromModel(item))
	}

	return httpHelper.Response(map[string]interface{}{
		"items": reservationsItem,
	}, http.StatusOK)
}

// Store a new item
func (h *Handler) Store(request httpHelper.Req) (httpHelper.Res, error) {
	var reservation Reservation
	userEmail := request.RequestContext.Authorizer["email"].(string)
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
	domain, err := utils.GetDomain(userEmail)
	if err != nil {
		return httpHelper.ErrResponse(errInvalidRoom, http.StatusInternalServerError)
	}
	if !(strings.Contains(domain, validators.CokeDomain) && strings.HasPrefix(reservation.RoomName, "C") ||
		strings.Contains(domain, validators.PepsiDomain) && strings.HasPrefix(reservation.RoomName, "P")) {
		return httpHelper.ErrResponse(errInvalidRoom, http.StatusForbidden)
	}

	// Validate that time is empty or the same user
	storedReservation, err := h.repository.Get(buildModel(reservation, userEmail))
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}
	if storedReservation != nil && storedReservation.UserEmail != "" && storedReservation.UserEmail != userEmail {
		return httpHelper.ErrResponse(errForbidden, http.StatusForbidden)
	}

	// Storing
	if err := h.repository.Store(buildModel(reservation, userEmail)); err != nil {
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
	reservation := Reservation{
		RoomName: room,
		Time:     time,
	}
	userEmail := request.RequestContext.Authorizer["email"].(string)

	// Validating body
	err := h.validator.Struct(reservation)
	if err != nil {
		return httpHelper.ErrResponse(validators.ParseValidationErrors(err), http.StatusBadRequest)
	}

	// Validate user is the owner
	storedReservation, err := h.repository.Get(buildModel(reservation, userEmail))
	if err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}
	if storedReservation == nil || storedReservation.UserEmail != userEmail {
		return httpHelper.ErrResponse(errForbidden, http.StatusForbidden)
	}

	// Deleting
	if err := h.repository.Delete(buildModel(reservation, userEmail)); err != nil {
		return httpHelper.ErrResponse(err, http.StatusInternalServerError)
	}

	return httpHelper.Response(map[string]bool{
		"delete": true,
	}, http.StatusOK)
}
