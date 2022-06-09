package database

import "github.com/Hamifthi/authentication_microservice/entity"

type DatabaseServiceMock struct {
	MockedGetUser    func(email string) (entity.User, error)
	MockedCreateUser func(email, hashPass, tokenHash string) error
}

func (dsm *DatabaseServiceMock) GetUser(email string) (entity.User, error) {
	return dsm.MockedGetUser(email)
}

func (dsm *DatabaseServiceMock) CreateUser(email, hashPass, tokenHash string) error {
	return dsm.MockedCreateUser(email, hashPass, tokenHash)
}
