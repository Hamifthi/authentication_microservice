package authentication

import (
	"errors"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/database"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func initializeAuthAndDBService() (*AuthenticationService, *database.DatabaseServiceMock) {
	dbService := database.DatabaseServiceMock{}
	logger := log.Logger{}
	authService := New(&dbService, &logger)
	return authService, &dbService
}

func TestSignUpWithInvalidEmail(t *testing.T) {
	authService, _ := initializeAuthAndDBService()
	email := "test.com"
	password := "123test123"
	err := authService.SignUp(email, password)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "The email address is invalid")
}

func TestSignUpExistedEmail(t *testing.T) {
	authService, dbService := initializeAuthAndDBService()
	email := "test@test.com"
	password := "123test123"
	user := entity.User{Email: email, Password: password}
	dbService.MockedGetUser = func(email string) (entity.User, error) {
		return user, nil
	}
	err := authService.SignUp(email, password)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "the user with %s email is already exist", email)
}

func TestSignUpInvalidPassword(t *testing.T) {
	authService, dbService := initializeAuthAndDBService()
	err := internal.InitializeEnv("../../.env")
	dbService.MockedGetUser = func(email string) (entity.User, error) {
		return entity.User{}, nil
	}
	email := "test@test.com"
	password := "123test123"
	err = authService.SignUp(email, password)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "insecure password, try including more special"+
		" characters, using uppercase letters or using a longer password")
}

func TestSignUpSuccessfully(t *testing.T) {
	authService, dbService := initializeAuthAndDBService()
	_ = internal.InitializeEnv("../../test.env")
	dbService.MockedGetUser = func(email string) (entity.User, error) {
		return entity.User{}, nil
	}
	dbService.MockedCreateUser = func(email, hashPass, tokenHash string) error {
		return nil
	}
	email := "test@test.com"
	password := "587@_Testing123"
	err := authService.SignUp(email, password)
	assert.Nil(t, err)
}

func TestSignInWithInvalidEmail(t *testing.T) {
	authService, _ := initializeAuthAndDBService()
	email := "test.com"
	password := "123test123"
	_, err := authService.SignIn(email, password)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "The email address is invalid")
}

func TestSignInUserNotExist(t *testing.T) {
	authService, dbService := initializeAuthAndDBService()
	dbService.MockedGetUser = func(email string) (entity.User, error) {
		return entity.User{}, errors.New("user doesn't exist")
	}
	email := "test@test.com"
	password := "123test123"
	_, err := authService.SignIn(email, password)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "the user with %s email doesn't exist", email)
}

func TestSignInUserSuccessfully(t *testing.T) {
	authService, dbService := initializeAuthAndDBService()
	_ = internal.InitializeEnv("../../test.env")
	email := "test@test.com"
	password := "587@_Testing123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := entity.User{Email: email, Password: password, HashedPassword: string(hashedPassword)}
	dbService.MockedGetUser = func(email string) (entity.User, error) {
		return user, nil
	}
	tokens, err := authService.SignIn(email, password)
	assert.Nil(t, err)
	assert.NotNil(t, tokens)
}
