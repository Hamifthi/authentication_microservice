package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"log"
	"net/http"
)

type AuthenticationHandler struct {
	authService *authenticationService
	l           *log.Logger
}

func NewHandler(authService *authenticationService, l *log.Logger) *AuthenticationHandler {
	return &AuthenticationHandler{authService, l}
}

type keyUser struct{}

func (ah *AuthenticationHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := entity.User{}

		err := user.FromJson(r.Body)

		if err != nil {
			ah.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		err = user.Validate()
		if err != nil {
			ah.l.Println("[ERROR] validating user", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating user: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		ctx := context.WithValue(r.Context(), keyUser{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func (ah *AuthenticationHandler) UserSignUp(rw http.ResponseWriter, r *http.Request) {
	ah.l.Println("Handle Sign up of User")
	user := r.Context().Value(keyUser{}).(entity.User)
	err := ah.authService.SignUp(user.Email, user.Password)
	if err != nil {
		ah.l.Printf("[ERROR] signing up user has %s error", err)
		http.Error(rw, "Unable to signing up the user", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("User successfully signed up"))
}

func (ah *AuthenticationHandler) UserLogin(rw http.ResponseWriter, r *http.Request) {
	ah.l.Println("Handle User Login")
	user := r.Context().Value(keyUser{}).(entity.User)
	tokens, err := ah.authService.SignIn(user.Email, user.Password)
	if err != nil {
		ah.l.Printf("[ERROR] signing up user has %s error", err)
		http.Error(rw, "Unable to signing in the user", http.StatusBadRequest)
		return
	}
	jsonResponse, err := json.Marshal(tokens)
	if err != nil {
		ah.l.Printf("[ERROR] happened in JSON marshal. Err: %s", err)
		http.Error(rw, "Unable to signing in the user", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonResponse)
}
