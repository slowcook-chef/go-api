package database

import (
	"fmt"
	//"lambda-func/ledger"
	"lambda-func/types"
	"lambda-func/types/schedule"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	USER_TABLE_NAME     = "userTable"
	SCHEDULE_TABLE_NAME = "gamesTable"
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
	//POOP
	// // Initialize the batch write structure
	// batchInput := &dynamodb.BatchWriteItemInput{
	// 	RequestItems: map[string][]*dynamodb.WriteRequest{},
	// }

	// // Accumulate games into the batch
	// for _, game := range schedule.GameData.Games {
	// 	home, away := game.GetTeams()
	// 	day, time := game.GetDateTime()
	// 	exists, _ := dbclient.DoesGameExist(&game)
	// 	if exists {
	// 		continue
	// 	}

	// 	// Create a PutRequest for each game
	// 	putRequest := &dynamodb.WriteRequest{
	// 		PutRequest: &dynamodb.PutRequest{
	// 			Item: map[string]*dynamodb.AttributeValue{
	// 				"gameId": {
	// 					S: aws.String(game.GetID()),
	// 				},
	// 				"awayTeam": {
	// 					S: aws.String(away),
	// 				},
	// 				"homeTeam": {
	// 					S: aws.String(home),
	// 				},
	// 				"day": {
	// 					S: aws.String(day),
	// 				},
	// 				"time": {
	// 					S: aws.String(time),
	// 				},
	// 				"status": {
	// 					S: aws.String(game.GetStatus()),
	// 				},
	// 			},
	// 		},
	// 	}

	// 	// Add the PutRequest to the batch
	// 	batchInput.RequestItems[SCHEDULE_TABLE_NAME] = append(batchInput.RequestItems[SCHEDULE_TABLE_NAME], putRequest)
	// }

	// // Execute the batch write
	// output, err := dbclient.databaseStore.BatchWriteItem(batchInput)
	// if err != nil {
	// 	return fmt.Errorf("error executing batch write: %v", err)
	// }
	// ledger.LogHandlerProcess("!SUCCESS SCHEDULED!" + output.String())
	// // Handle any unprocessed items...
	// // For simplicity, this part is omitted. In a real scenario, you should handle retries for unprocessed items.

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

func (client DynamoDBClient) GetMLBSchedule() (mlbSchedule *schedule.Schedule, e error) {
	//Query database
	return mlbSchedule, nil

	//POOP
	// ledger.LogHandlerProcess("READING datastore schedules")
	// dbResult, err := client.databaseStore.GetItem(&dynamodb.GetItemInput{
	// 	TableName: aws.String(SCHEDULE_TABLE_NAME),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"scheduleID": {
	// 			S: aws.String("mlb"),
	// 		},
	// 	},
	// })
	// //Validate query result
	// if err != nil {
	// 	ledger.LogError(&err)
	// 	return mlbSchedule, fmt.Errorf("error[%w]: failed getting schedule in database", err)
	// }
	// //Couldnt find entry?
	// if dbResult.Item == nil {
	// 	ledger.LogHandlerProcess("NO ITEM FOUND, calling enpoint")
	// 	return schedule.GetMLBEndpointSchedule()
	// }
	// //Translate Item
	// err = dynamodbattribute.UnmarshalMap(dbResult.Item, &mlbSchedule)
	// if err != nil {
	// 	ledger.LogError(&err)
	// 	return mlbSchedule, err
	// }

	// return mlbSchedule, err
}

//"Internal server error"
