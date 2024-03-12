package user

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FetchUser(email string, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	result, err := dynoClient.GetItem(input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, errors.New("user not found")
	}

	user := new(User)
	if err := dynamodbattribute.UnmarshalMap(result.Item, user); err != nil {
		return nil, errors.New("error unmarshalling user data")
	}

	return user, nil
}

func FetchUsers(tableName string, dynoClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynoClient.Scan(input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("no users found")
	}

	var users []User
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &users); err != nil {
		return nil, errors.New("error unmarshalling user data")
	}

	return &users, nil
}

func (user *User) CreateUser(req events.APIGatewayProxyRequest, tableName string, dynoClient dynamodbiface.DynamoDBAPI) (string, error) {

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return "", err
	}

	fetchedUser, _ := FetchUser(user.Email, tableName, dynoClient)
	if fetchedUser != nil && len(fetchedUser.Email) != 0 {
		return "", errors.New("User already exist")
	}
	input := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	_, err = dynoClient.PutItem(&input)
	if err != nil {
		return "", err
	}
	return "created succesfully", nil
}

func (user *User) UpdateUser(tableName string, dynoClient dynamodbiface.DynamoDBAPI) (string, error) {

	u, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return "", err
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      u,
	}

	_, err = dynoClient.PutItem(&input)
	if err != nil {
		return "", err
	}
	return "", nil
}

// expression.
func (user *User) UpdateMovie(tableName string, dynoClient dynamodbiface.DynamoDBAPI) (any, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	update := expression.Set(expression.Name("name"), expression.Value(user.Name))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		updateInput := dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableName),
			Key:                       user.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              aws.String(dynamodb.ReturnValueUpdatedNew),
		}
		response, err = dynoClient.UpdateItem(&updateInput)
		if err != nil {
			log.Printf("Couldn't update movie %v. Here's why: %v\n", user.Name, err)
		} else {
			err = dynamodbattribute.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}
	fmt.Println("updated")
	return attributeMap, err
}

func (user *User) DeleteUser(tableName string, dynoClient dynamodbiface.DynamoDBAPI) error {
	_, err := dynoClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName), Key: user.GetKey(),
	})
	if err != nil {
		log.Printf("Couldn't delete %v from the table. Here's why: %v\n", user.Name, err)
	}
	return err
}

func (user *User) GetKey() map[string]*dynamodb.AttributeValue {
	name, err := dynamodbattribute.Marshal(user.Name)
	if err != nil {
		panic(err)
	}
	email, err := dynamodbattribute.Marshal(user.Email)
	if err != nil {
		panic(err)
	}
	return map[string]*dynamodb.AttributeValue{"name": name, "email": email}
}

// String returns the title, year, rating, and plot of a movie, formatted for the example.
func (user *User) String() string {
	return fmt.Sprintf("%v\n\tName: %v\n\tEmail:",
		user.Name, user.Email)
}
