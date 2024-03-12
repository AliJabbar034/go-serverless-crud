package main

import (
	"fmt"

	"github.com/alijabbar034/serverless-crud/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynoClient dynamodbiface.DynamoDBAPI
)

func main() {

	fmt.Println("starting serverless......")

	awssession, err := session.NewSession()
	if err != nil {
		return
	}
	dynoClient = dynamodb.New(awssession)
	lambda.Start(handler)
}

const tableName string = "User"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynoClient)

	case "POST":
		return handlers.CreateUser(req, tableName, dynoClient)

	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynoClient)

	case "PUT":
		return handlers.UpdateUser(req, tableName, dynoClient)

	default:
		return handlers.UnhandledMethod()

	}
}
