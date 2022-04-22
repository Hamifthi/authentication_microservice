package databaseService

import "github.com/Hamifthi/authentication_microservice/entity"

type DatabaseInterface interface {
	GetUser(email string) (entity.User, error)
	CreateUser(email, hashedPass, tokenHash string) error
}
