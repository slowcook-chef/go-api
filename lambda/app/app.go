package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	//we actually initialize database store
	//gets passed UP to the he api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)
	return App{
		ApiHandler: apiHandler,
	}
}
