package authentication_service

import (
	"fmt"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/database_service"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strconv"
)

type authenticationService struct {
	dbService database_service.DatabaseServiceInterface
}

func newAuthenticationService(dbService *database_service.DatabaseServiceInterface) AuthenticationServiceInterface {
	return &authenticationService{
		dbService: *dbService,
	}
}

func (a *authenticationService) signUp(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("the email address is invalid")
	}
	entropyBits, err := internal.InitializeAndGetEnv("MINENTROPYBITS")
	if err != nil {
		panic("Problem getting the min entropy bits")
	}
	minEntropyBits, err := strconv.ParseFloat(entropyBits, 64)
	if err != nil {
		panic("Problem converting the min entropy bits to the float64")
	}
	err = passwordValidator.Validate(password, minEntropyBits)
	if err != nil {
		return err
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("the hashing process of password went wrong")
	}
}
