package db

import (
	"errors"
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"gorm.io/gorm"
	"log"
)

type databaseService struct {
	db     *gorm.DB
	logger *log.Logger
}

func New(db *gorm.DB, logger *log.Logger) *databaseService {
	return &databaseService{db, logger}
}

func (d *databaseService) GetUser(email string) (entity.User, error) {
	var user entity.User
	result := d.db.First(&user, "email = ?", email)
	if result.Error != nil {
		d.logger.Println("[Error] occurred while fetching the user")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, errors.New("User not found")
		} else {
			return user, fmt.Errorf("Error fetching user with %s email from database", email)
		}
	}
	return user, nil
}

func (d *databaseService) CreateUser(email, hashPass, tokenHash string) error {
	user := entity.User{Email: email, HashedPassword: hashPass, TokenHash: tokenHash}
	result := d.db.Create(&user)
	if result.Error != nil && result.RowsAffected != 1 {
		d.logger.Println("[Error] creating the user in the database")
		return result.Error
	}
	return nil
}
