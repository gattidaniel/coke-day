package http

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// ResponseError -
type ResponseError map[string]error

// Req is an alias for an api gateway request
type Req events.APIGatewayProxyRequest

// Res is an alis for an api gateway response
type Res events.APIGatewayProxyResponse

// Response is a wrapper around the api gateway proxy response, which takes
// a interface argument to be marshalled to json and returned, and an error code
func Response(data interface{}, code int) (Res, error) {
	body, _ := json.Marshal(data)
	return Res{
		Body:       string(body),
		StatusCode: code,
	}, nil
}

// ErrResponse returns an error in a specified format
func ErrResponse(err error, code int) (Res, error) {
	data := map[string]string{
		"err": err.Error(),
	}
	body, _ := json.Marshal(data)

	return Res{
		Body:       string(body),
		StatusCode: code,
	}, nil // I return nil so the client can receive the custom error
}

// ParseBody takes the body from the request, parses the json to a given struct pointer
func ParseBody(request Req, castTo interface{}) error {
	return json.Unmarshal([]byte(request.Body), &castTo)
}

// RequestHandleFunc is an alias for an api gateway request signature
type RequestHandleFunc func(request Req) (Res, error)
