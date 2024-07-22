package database

import (
	"fmt"
	"lambda-func/types"
	"lambda-func/types/schedule"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	USER_TABLE_NAME     = "userTable"
	SCHEDULE_TABLE_NAME = "scheduleTable"
)

type DataStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	InsertSchedule(schedule schedule.Schedule) error
	GetUser(username string) (types.User, error)
	GetMLBSchedule() (SCH *schedule.Schedule, e error)
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return true, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (dbClient DynamoDBClient) DoesGameExist(game *schedule.Game) (bool, error) {
	result, err := dbClient.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"game_id": {
				S: aws.String((*game).GetID()),
			},
		},
	})
	if err != nil {
		return true, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

// /////////////////////////////////////////////////////////////////////////////Inserters
func (u DynamoDBClient) InsertUser(user types.User) error {
	//assemble the item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := u.databaseStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}
func (dbclient DynamoDBClient) InsertSchedule(schedule schedule.Schedule) error {
	var item *dynamodb.PutItemInput
	//assemble the item
	for _, game := range schedule.GameData.Games {
		home, away := game.GetTeams()
		day, time := game.GetDateTime()
		exists, _ := dbclient.DoesGameExist(&game)
		if exists {
			continue
		}
		item = &dynamodb.PutItemInput{
			TableName: aws.String(SCHEDULE_TABLE_NAME),
			Item: map[string]*dynamodb.AttributeValue{
				"game_id": {
					S: aws.String(game.GetID()),
				},
				"away_team": {
					S: aws.String(away),
				},
				"home_team": {
					S: aws.String(home),
				},
				"day": {
					S: aws.String(day),
				},
				"time": {
					S: aws.String(time),
				},
				"status": {
					S: aws.String(game.GetStatus()),
				},
			},
		}
		_, err := dbclient.databaseStore.PutItem(item)
		if err != nil {
			return fmt.Errorf("error putting item in schedule table")
		}
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////Getters
func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}
	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (client DynamoDBClient) GetMLBSchedule() (SCH *schedule.Schedule, e error) {
	//Query database
	dbResult, err := client.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"scheduleID": {
				S: aws.String("mlb"),
			},
		},
	})
	//Validate query result
	if err != nil {
		return SCH, fmt.Errorf("error[%w]: failed getting schedule in database", err)
	}
	//Couldnt find entry?
	if dbResult.Item == nil {
		return schedule.GetMLBEndpointSchedule()
	}
	//Translate Item
	err = dynamodbattribute.UnmarshalMap(dbResult.Item, &SCH)
	if err != nil {
		return SCH, err
	}

	return SCH, err
}

//"Internal server error"
