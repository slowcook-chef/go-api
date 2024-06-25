package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient{
	dbSession:=session.Must(session.NewSession())
	db:= dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (user DynamoDBClient)DoesUserExist(username string)(bool,error){
	result, err :=
}