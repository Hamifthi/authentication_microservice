package database_service

import "github.com/Hamifthi/authentication_microservice/entity"

type DatabaseServiceInterface interface {
	getUser(email string) (entity.User, error)
}
