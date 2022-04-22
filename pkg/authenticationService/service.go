package authenticationService

import (
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/databaseService"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/mail"
	"strconv"
	"time"
)

type authenticationService struct {
	dbService databaseService.DatabaseInterface
	logger    *log.Logger
}

func New(dbService *databaseService.DatabaseInterface, logger *log.Logger) *authenticationService {
	return &authenticationService{
		dbService: *dbService,
		logger:    logger,
	}
}

func (a *authenticationService) GenerateAccessToken(email string) (string, error) {
	jwtExpirationStr, err := internal.InitializeAndGetEnv("JwtExpiration")
	if err != nil {
		a.logger.Println("Error reading jwt expiration key")
		return "", errors.Wrap(err, "Error reading jwt expiration")
	}
	jwtExpiration, _ := strconv.Atoi(jwtExpirationStr)

	claims := jwt.MapClaims{
		"iss": "authService",
		"exp": time.Now().Add(time.Minute * time.Duration(jwtExpiration)).Unix(),
		"data": map[string]string{
			"userEmail": email,
			"tokenType": "access",
		},
	}
	accessTokenPrivateKey, err := internal.InitializeAndGetEnv("AccessTokenPrivateKey")
	if err != nil {
		a.logger.Println("Error reading access token private key")
		return "", errors.Wrap(err, "Error reading access token private key")
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(accessTokenPrivateKey))
	if err != nil {
		a.logger.Println("Unable to parse the access token private key")
		return "", errors.New("Unable to parse the access token private key")
	}
	// its better use environment variable here
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

func (a *authenticationService) GenerateRefreshToken(email, tokenHash string) (string, error) {
	customKey := internal.GenerateCustomKey(email, tokenHash)
	claims := jwt.MapClaims{
		"iss": "authService",
		"data": map[string]string{
			"userEmail": email,
			"customKey": customKey,
			"tokenType": "refresh",
		},
	}
	refreshTokenPrivateKey, err := internal.InitializeAndGetEnv("RefreshTokenPrivateKey")
	if err != nil {
		a.logger.Println("Error reading refresh token private key")
		return "", errors.Wrap(err, "Error reading refresh token private key")
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(refreshTokenPrivateKey))
	if err != nil {
		a.logger.Println("Unable to parse the refresh token private key")
		return "", errors.New("Unable to parse the refresh token private key")
	}
	// its better use environment variable here
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

func (a *authenticationService) SignUp(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.Wrap(err, "The email address is invalid")
	}
	user, err := a.dbService.GetUser(email)
	if user.Email != "" {
		return fmt.Errorf("the user with %s email is already exist", email)
	}
	entropyBits, err := internal.InitializeAndGetEnv("MinEntropyBits")
	if err != nil {
		return errors.Wrap(err, "Problem getting the min entropy bits from config file")
	}
	minEntropyBits, err := strconv.ParseFloat(entropyBits, 64)
	if err != nil {
		return errors.Wrap(err, "Problem converting the min entropy bits to the float64")
	}
	err = passwordValidator.Validate(password, minEntropyBits)
	if err != nil {
		return err
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "The hashing process of password went wrong")
	}
	// its better use environment variable here
	tokenHash := internal.RandString(15)
	err = a.dbService.CreateUser(email, string(hashedPass), tokenHash)
	if err != nil {
		return errors.Wrap(err, "The user can't be inserted to the database")
	}
	return nil
}

func (a *authenticationService) SignIn(email, password string) (entity.Tokens, error) {
	emptyTokens := entity.Tokens{AccessToken: "", RefreshToken: ""}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return emptyTokens, errors.Wrap(err, "The email address is invalid")
	}
	user, err := a.dbService.GetUser(email)
	if user.Email == "" {
		return emptyTokens, fmt.Errorf("the user with %s email doesn't exist", email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return emptyTokens, errors.Wrap(err, "The invalid credentials, please try again.")
	}
	accessToken, err := a.GenerateAccessToken(email)
	if err != nil {
		a.logger.Println("Unable to get access token")
		return emptyTokens, errors.New("Unable to get access token")
	}
	refreshToken, err := a.GenerateRefreshToken(email, user.TokenHash)
	if err != nil {
		a.logger.Println("Unable to get refresh token")
		return emptyTokens, errors.New("Unable to get refresh token")
	}
	return entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
