package authentication

import (
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/database"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/mail"
	"strconv"
	"time"
)

type AuthenticationService struct {
	dbService database.DatabaseInterface
	logger    *log.Logger
}

func New(dbService database.DatabaseInterface, logger *log.Logger) *AuthenticationService {
	return &AuthenticationService{
		dbService: dbService,
		logger:    logger,
	}
}

func (a *AuthenticationService) readPrivateKey() ([]byte, error) {
	accessTokenPrivateKeyPath, err := internal.GetEnv("TokenPrivateKeyPath")
	if err != nil {
		a.logger.Println("[Error] reading access token private key path from environment")
		return nil, errors.Wrap(err, "Error reading access token private key path")
	}
	signBytes, err := ioutil.ReadFile(accessTokenPrivateKeyPath)
	if err != nil {
		a.logger.Println("[Error] reading access token private key")
		return nil, errors.Wrap(err, "Error reading access token private key")
	}
	return signBytes, nil
}

func (a *AuthenticationService) readPublicKey() ([]byte, error) {
	accessTokenPublicKeyPath, err := internal.GetEnv("TokenPublicKeyPath")
	if err != nil {
		a.logger.Println("[Error] reading access token public key path from environment")
		return nil, errors.Wrap(err, "Error reading access token public key path")
	}
	verifyBytes, err := ioutil.ReadFile(accessTokenPublicKeyPath)
	if err != nil {
		a.logger.Println("[Error] reading access token public key")
		return nil, errors.Wrap(err, "Error reading access token public key")
	}
	return verifyBytes, nil
}

func (a *AuthenticationService) generateAccessToken(email string) (string, error) {
	jwtExpirationStr, err := internal.GetEnv("JwtExpiration")
	if err != nil {
		a.logger.Println("[Error] reading jwt expiration key")
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
	signBytes, err := a.readPrivateKey()
	if err != nil {
		a.logger.Println("[Error] reading access token private key")
		return "", errors.Wrap(err, "Error reading access token private key")
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		a.logger.Println("[Error] Unable to parse the access token private key")
		return "", errors.New("Unable to parse the access token private key")
	}
	// its better use environment variable here
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

type RefreshTokenCustomClaims struct {
	Foo string `json:"foo"`
	jwt.MapClaims
}

func (a *AuthenticationService) generateRefreshToken(email, tokenHash string) (string, error) {
	customKey := internal.GenerateCustomKey(email, tokenHash)
	claims := RefreshTokenCustomClaims{
		"authServiceCustomClaims",
		jwt.MapClaims{
			"iss": "authService",
			"data": map[string]string{
				"userEmail": email,
				"customKey": customKey,
				"tokenType": "refresh",
			},
		},
	}
	signBytes, err := a.readPrivateKey()
	if err != nil {
		a.logger.Println("[Error] reading access token private key")
		return "", errors.Wrap(err, "Error reading access token private key")
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		a.logger.Println("Unable to parse the refresh token private key")
		return "", errors.New("Unable to parse the refresh token private key")
	}
	// its better use environment variable here
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

func (a *AuthenticationService) SignUp(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.Wrap(err, "The email address is invalid")
	}
	user, err := a.dbService.GetUser(email)
	if user.Email != "" {
		return errors.Errorf("the user with %s email is already exist", email)
	}
	entropyBits, err := internal.GetEnv("MinEntropyBits")
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

func (a *AuthenticationService) SignIn(email, password string) (entity.Tokens, error) {
	emptyTokens := entity.Tokens{AccessToken: "", RefreshToken: ""}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return emptyTokens, errors.Wrap(err, "The email address is invalid")
	}
	user, err := a.dbService.GetUser(email)
	if user.Email == "" {
		return emptyTokens, errors.Wrapf(err, "the user with %s email doesn't exist", email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return emptyTokens, errors.Wrap(err, "The invalid credentials, please try again.")
	}
	accessToken, err := a.generateAccessToken(email)
	if err != nil {
		a.logger.Println("Unable to get access token")
		return emptyTokens, errors.New("Unable to get access token")
	}
	refreshToken, err := a.generateRefreshToken(email, user.TokenHash)
	if err != nil {
		a.logger.Println("Unable to get refresh token")
		return emptyTokens, errors.New("Unable to get refresh token")
	}
	return entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *AuthenticationService) ValidateRefreshToken(refreshToken string) (entity.User, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			a.logger.Println("[Error] unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		verifyBytes, err := a.readPublicKey()
		if err != nil {
			a.logger.Println("[Error] reading access token public key")
			return "", errors.Wrap(err, "Error reading access token private key")
		}
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			a.logger.Println("[Error] unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})
	user := entity.User{}
	if err != nil {
		a.logger.Println("[Error] parsing the claims from refresh token")
		return user, errors.Wrap(err, "Error parsing the claims from refresh token")
	}
	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	data, ok := claims.MapClaims["data"].(map[string]string)
	if !ok || !token.Valid || data["userEmail"] == "" || data["tokenType"] != "refresh" {
		a.logger.Println("[Error] getting claims from token")
		return user, errors.New("Error getting claims from token")
	}
	user, err = a.dbService.GetUser(data["userEmail"])
	if err != nil {
		a.logger.Println("[Error] can't retrieve user from database")
		return user, errors.Wrap(err, "Error can't retrieve user from database")
	}
	generatedCustomKey := internal.GenerateCustomKey(user.Email, user.TokenHash)
	if data["customKey"] != generatedCustomKey {
		a.logger.Println("[Error] refresh token is malformed")
		return user, errors.New("Refresh token is malformed")
	}
	return user, nil
}

func (a *AuthenticationService) RefreshAccessToken(user entity.User) (string, error) {
	accessToken, err := a.generateAccessToken(user.Email)
	if err != nil {
		a.logger.Println("Unable to refresh access token")
		return "", errors.Wrap(err, "Unable to refresh access token")
	}
	return accessToken, nil
}
