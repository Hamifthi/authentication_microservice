package authentication_service

type AuthenticationServiceInterface interface {
	signUp(email, password string) error
}
