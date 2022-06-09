package adapters

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication/adapters/graph/generated"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication/adapters/graph/model"
)

func (r *mutationResolver) SignUp(ctx context.Context, input model.UserInput) (string, error) {
	r.Logger.Println("Handle sign up of the user in GraphQL server")
	user := entity.User{Email: input.Email, Password: input.Password}
	err := user.Validate()
	if err != nil {
		r.Logger.Printf("[ERROR] validating the user email and password has %s error", err)
		return "", err
	}
	err = r.AuthService.SignUp(user.Email, user.Password)
	if err != nil {
		r.Logger.Printf("[ERROR] signing up user has %s error", err)
		return "", err
	}
	return "User successfully signed up", nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.UserInput) (*model.Tokens, error) {
	r.Logger.Println("Handle login of the user in GraphQL server")
	user := entity.User{Email: input.Email, Password: input.Password}
	err := user.Validate()
	if err != nil {
		r.Logger.Printf("[ERROR] validating the user email and password has %s error", err)
		return nil, err
	}
	authsrcTokens, err := r.AuthService.SignIn(user.Email, user.Password)
	if err != nil {
		r.Logger.Printf("[ERROR] login user has %s error", err)
		return nil, err
	}
	tokens := &model.Tokens{
		Access:  authsrcTokens.AccessToken,
		Refresh: authsrcTokens.RefreshToken,
	}
	return tokens, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
