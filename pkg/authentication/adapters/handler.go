package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication"
	"log"
	"net/http"
	"strings"
)

type AuthenticationHandler struct {
	authService *authentication.AuthenticationService
	l           *log.Logger
}

func NewHandler(authService *authentication.AuthenticationService, l *log.Logger) *AuthenticationHandler {
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

func (ah *AuthenticationHandler) MiddlewareValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authContent := strings.Split(authHeader, " ")
		token := authContent[1]
		if len(authContent) != 2 {
			ah.l.Println("[ERROR] Authorization token not provided or malformed")
			http.Error(
				rw,
				fmt.Sprint("Error authorization token not provided or malformed"),
				http.StatusUnauthorized,
			)
			return
		}
		user, err := ah.authService.ValidateRefreshToken(token)
		if user == (entity.User{}) || err != nil {
			ah.l.Println("[ERROR] Refresh token isn't valid")
			http.Error(
				rw,
				fmt.Sprintf("Error refresh token isn't valid due to %s", err),
				http.StatusUnauthorized,
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
		ah.l.Printf("[ERROR] login user has %s error", err)
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
