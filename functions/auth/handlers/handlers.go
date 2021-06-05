package handlers

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/coke-day/pkg/jwt"
	"strings"
)

var errUnauthorized = errors.New("Unauthorized")

// Handler -
type Handler struct {
	jwtHandler jwt.JWT
}

func NewHandler(jwtHandler jwt.JWT) Handler {
	return Handler{
		jwtHandler: jwtHandler,
	}
}

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
