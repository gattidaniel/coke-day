package handlers

import (
	"errors"
	httpHelper "github.com/coke-day/pkg/http"
	"net/http"
)

// CreateRoomRouting routes restful endpoints to the correct method
// GET calls the Search method (may receive by parameter room or time and filter by this)
// POST calls the Schedule method.
// DELETE calls the Delete method.
func (h *Handler) CreateRoomRouting() func(request httpHelper.Req) (httpHelper.Res, error) {
	router := func(request httpHelper.Req) (httpHelper.Res, error) {
		switch request.HTTPMethod {
		case "GET":
			return h.Search(request)
		case "DELETE":
			return h.Delete(request)
		case "POST":
			return h.Store(request)
		default:
			return httpHelper.ErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}
	}
	return router
}
