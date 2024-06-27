package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)
type ApiHandler struct{
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error{
	if event.Username == "" || event.Pasword == "" {
		return fmt.Errorf("request has empty parameters")
	}

	//does a user with this username already exist?

	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil{
		return fmt.Errorf("there an error checking if user exists %w", err)
	}
	if userExists{
		return fmt.Errorf("a user with that username already exists")
	}
	
	//we know a user does not exist

	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("error registering the user %w", err)
	}

	return nil
}