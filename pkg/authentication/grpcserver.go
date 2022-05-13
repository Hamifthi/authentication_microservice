package authentication

import (
	"context"
	"github.com/Hamifthi/authentication_microservice/entity"
	protos "github.com/Hamifthi/authentication_microservice/pkg/authentication/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type AuthServiceServer struct {
	authService *authenticationService
	l           *log.Logger
	protos.UnimplementedAuthServiceServer
}

func NewAuthServer(authService *authenticationService, l *log.Logger) *AuthServiceServer {
	return &AuthServiceServer{authService, l, protos.UnimplementedAuthServiceServer{}}
}

func (ass *AuthServiceServer) SignUp(ctx context.Context, req *protos.SignUpRequest) (*protos.SignUpResponse, error) {
	ass.l.Println("Handle Sign up of the User In Grpc Server")
	user := entity.User{Email: req.Email, Password: req.Password}
	err := user.Validate()
	if err != nil {
		grpcErr := status.Newf(
			codes.InvalidArgument,
			"Error invalid argument %s",
			err,
		)
		return nil, grpcErr.Err()
	}
	err = ass.authService.SignUp(user.Email, user.Password)
	if err != nil {
		grpcErr := status.Newf(
			codes.Internal,
			"Error get %s error when trying to sign up the user",
			err,
		)
		return nil, grpcErr.Err()
	}
	return &protos.SignUpResponse{Status: int64(codes.OK)}, nil
}
func (ass *AuthServiceServer) Login(ctx context.Context, req *protos.LoginRequest) (*protos.LoginResponse, error) {
	ass.l.Println("Handle Sign in of the User In Grpc Server")
	user := entity.User{Email: req.Email, Password: req.Password}
	err := user.Validate()
	if err != nil {
		grpcErr := status.Newf(
			codes.InvalidArgument,
			"Error invalid argument %s",
			err,
		)
		return nil, grpcErr.Err()
	}
	tokens, err := ass.authService.SignIn(user.Email, user.Password)
	if err != nil {
		grpcErr := status.Newf(
			codes.Internal,
			"Error get %s error when trying to sign in the user",
			err,
		)
		return nil, grpcErr.Err()
	}
	return &protos.LoginResponse{
		Status:       int64(codes.OK),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
