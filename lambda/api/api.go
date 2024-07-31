package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/ledger"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dataStore database.DataStore
}

// Instantiate API handler
func NewApiHandler(dbStore database.DataStore) ApiHandler {
	return ApiHandler{
		dataStore: dbStore,
	}
}

// Resgister User
func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//Recieve "RegisterUser" passenger(data) plane(request)
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	//Plane lands safely?
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	//Passengers pass through immigration(validity) checks
	if registerUser.Username == "" || registerUser.Pasword == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request - fields empty",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	//does a passenger with this username already exist?
	userExists, err := api.dataStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error, does user exist",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	//Deport them
	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}
	//Let them through
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error, types new user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create new user %w", err)
	}

	//Passenger enters the country(database)
	err = api.dataStore.InsertUser(*user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error, inserting user",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	//Send thank you note
	return events.APIGatewayProxyResponse{
		Body:       "Successfully registered user",
		StatusCode: http.StatusOK,
	}, nil
}

// Login User
func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//Recieve login user INPUT, treat with increased security
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//Translate request bytes to known type
	var loginRequest LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	//Call data store to recieve data
	user, err := api.dataStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error, getting user",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	//Validate identity
	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid user credentials",
			StatusCode: http.StatusBadRequest,
		}, nil
	}
	//Give secure info
	accesToken := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token: "%s"}`, accesToken)
	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

// ////////////////// ////////////////// ////////////////// //
// ////////////////// Under Construction ////////////////// //
// ////////////////// ////////////////// ////////////////// // by Tomas
func (api ApiHandler) GetMLBSchedule(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//GET "GetMLBSchedule" (data) (request)
	ledger.LogHandlerStart("GetMLBSchedule() is invoked")
	//Call datastore, returns empty schedule
	mlbSchedule, err := api.dataStore.GetMLBSchedule()

	if err != nil {
		ledger.LogError(&err)
		return events.APIGatewayProxyResponse{
			Body:       "Get MLB Schedule error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// Translate
	scheduleJSON, err := json.Marshal(mlbSchedule)
	if err != nil {
		ledger.LogError(&err)
		return events.APIGatewayProxyResponse{
			Body:       "Error converting MLB schedule to JSON",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	ledger.LogHandlerEnd("GetMLBSchedule is succesful")
	//Send response OK
	return events.APIGatewayProxyResponse{
		Body:       string(scheduleJSON),
		StatusCode: http.StatusOK,
	}, nil

}
