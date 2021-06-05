package handlers

import "github.com/aws/aws-lambda-go/events"

// CreateAuthRouting routes restful endpoints to the correct method
// POST with Path '/register' calls the Register method
// POST with Path '/login' calls the Login method
func (h *Handler) CreateAuthRouting() func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	router := func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return h.Auth(request)
	}
	return router
}
