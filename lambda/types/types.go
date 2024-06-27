package types

import (
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Pasword string `json:"password"`
}

type User struct{
	Username string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser)(*User,error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Pasword),10)
	if err !=nil{
		return nil, err
	}

	return &User{
		Username: registerUser.Username,
		PasswordHash: string(hashedPassword),
	},nil
}

func ValidatePassword(hashedPassword, plainTextPassword string) bool {
	err:=bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
	
}