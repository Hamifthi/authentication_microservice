package adapters

import (
	"github.com/Hamifthi/authentication_microservice/pkg/authentication"
	"log"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService *authentication.AuthenticationService
	Logger      *log.Logger
}
