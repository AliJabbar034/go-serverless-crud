package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alijabbar034/serverless-crud/pkg/user"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMessage string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynoClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{
				ErrorMessage: fmt.Sprintf("invalid email %s", email),
			})
		}
		return apiResponse(http.StatusOK, result)
	}
	return nil, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	user := &user.User{}

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "proides all required fields",
		})
	}
	result, err := user.CreateUser(req, tableName, dynoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "cannot create",
		})
	}

	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	user := &user.User{}
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "invalid body",
		})
	}

	err := user.DeleteUser(tableName, dynoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "invalid body",
		})
	}

	return apiResponse(http.StatusOK, "successfully deleted")
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	use := &user.User{}

	if err := json.Unmarshal([]byte(req.Body), &use); err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "Bad request",
		})

	}
	fetched, _ := user.FetchUser(use.Email, tableName, dynoClient)
	if fetched == nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "not found",
		})
	}

	res, err := use.UpdateUser(tableName, dynoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMessage: "updation error",
		})
	}
	return apiResponse(http.StatusOK, res)

}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {

	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
