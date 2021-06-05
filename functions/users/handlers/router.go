package handlers

import (
	"errors"
	httpHelper "github.com/coke-day/pkg/http"
	"net/http"
	"strings"
)

// CreateUserRouting routes restful endpoints to the correct method
// POST with Path '/register' calls the Register method
// POST with Path '/login' calls the Login method
func (h Handler) CreateUserRouting() func(request httpHelper.Req) (httpHelper.Res, error) {
	router := func(request httpHelper.Req) (httpHelper.Res, error) {
		if request.HTTPMethod == "POST" && strings.ToLower(request.Path) == _registerPath {
			return h.Register(request)
		}
		if request.HTTPMethod == "POST" && strings.ToLower(request.Path) == _loginPath {
			return h.Login(request)
		}
		return httpHelper.ErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
	return router
}
