package authentication_service

import "github.com/Hamifthi/authentication_microservice/pkg/database_service"

type authenticationService struct {
	dbService database_service.DatabaseServiceInterface
}

func newAuthenticationService(dbService *database_service.DatabaseServiceInterface) AuthenticationServiceInterface {
	return &authenticationService{
		dbService: *dbService,
	}
}

func (a *authenticationService) signUp(email, password string) error {
	return nil
}
