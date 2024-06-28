package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Take payload and do something with it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	
	return fmt.Sprintf("Succesfully called by - %s", event.Username), nil
}

func main() {
	myApp:= app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)

		case "/login":
			return myApp.ApiHandler.LoginUSer(request)
		default:
			return events.APIGatewayProxyResponse{
				Body: "Not found",
				StatusCode: http.StatusNotFound,
			},nil
		}
	})
	fmt.Println("----App End----")
}