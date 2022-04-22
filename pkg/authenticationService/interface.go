package authenticationService

import "github.com/Hamifthi/authentication_microservice/entity"

type AuthenticationInterface interface {
	SignUp(email, password string) error
	SignIn(email, password string) (entity.Tokens, error)
}
