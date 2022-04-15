package database_service

import "github.com/Hamifthi/authentication_microservice/entity"

type databaseService struct{}

func newDatabaseService() DatabaseServiceInterface {
	return &databaseService{}
}

func (d *databaseService) getUser(email string) (entity.User, error) {
	return entity.User{}, nil
}
