package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func CreateToken(user User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour*1).Unix()
	claims := jwt.MapClaims{
		"user": user.Username,
		"expires": validUntil,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	//TODO: add this to secret manager!
	secret := "secret"
	tokenString,err := token.SignedString([]byte(secret))
	if err!=nil {
		return ""
	}
	return tokenString
}