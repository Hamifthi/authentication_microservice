package databaseService

import "github.com/Hamifthi/authentication_microservice/entity"

type databaseService struct{}

func New() *databaseService {
	return &databaseService{}
}

func (d *databaseService) GetUser(email string) (entity.User, error) {
	return entity.User{}, nil
}

func (d *databaseService) CreateUser(email, hashPass, tokenHash string) error {
	return nil
}
