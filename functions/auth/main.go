package main

import (
	"errors"
	"github.com/coke-day/pkg/jwt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler -
type Handler struct {
	jwtHandler jwt.JWT
}

var errUnauthorized = errors.New("Unauthorized")

func (h Handler) Auth(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	var bearerToken string
	if len(tokenSlice) > 1 {
		bearerToken = tokenSlice[len(tokenSlice)-1]
	}
	if bearerToken == "" {
		return events.APIGatewayCustomAuthorizerResponse{}, errUnauthorized
	}

	parseToken, err := h.jwtHandler.ParseToken(bearerToken)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errUnauthorized
	}

	return generatePolicy("user", "Allow", request.MethodArn, map[string]interface{}{"email": parseToken.Email}), nil
}

func main() {
	// Get JWT Secret
	jwtSecret := os.Getenv("JWTSECRET")

	// Create JWT
	jwtHandler := jwt.NewJWT(jwtSecret, "u-should-hire-me", 3*time.Hour)

	// Create the Auth instance
	handler := &Handler{jwtHandler}

	// Create the router
	router := createAuthRouting(handler)

	lambda.Start(router)
}

// createUserRouting routes restful endpoints to the correct method
// POST with Path '/register' calls the Register method
// POST with Path '/login' calls the Login method
func createAuthRouting(handler *Handler) func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	router := func(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		return handler.Auth(request)
	}
	return router
}

func generatePolicy(principalID, effect, resource string, context map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	authResponse.Context = context
	return authResponse
}
