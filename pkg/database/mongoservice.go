package database

import (
	"context"
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoDBService struct {
	collection *mongo.Collection
	ctx        context.Context
	logger     *log.Logger
}

func NewMongoSrv(collection *mongo.Collection, ctx context.Context, logger *log.Logger) *MongoDBService {
	return &MongoDBService{collection: collection, ctx: ctx, logger: logger}
}

func (d *MongoDBService) GetUser(email string) (entity.User, error) {
	var user entity.User
	err := d.collection.FindOne(d.ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		d.logger.Println("[Error] occurred while fetching the user from mongodb")
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, errors.New("User not found in mongodb")
		} else {
			return user, fmt.Errorf("Error fetching user with %s email from mongodb", email)
		}
	}
	return user, nil
}

func (d *MongoDBService) CreateUser(email, hashPass, tokenHash string) error {
	user := entity.User{Email: email, HashedPassword: hashPass, TokenHash: tokenHash}
	_, err := d.collection.InsertOne(d.ctx, &user, options.InsertOne())
	if err != nil {
		d.logger.Println("[Error] occurred while creating user in mongodb")
		return errors.Wrap(err, "Error occurred while creating user in mongodb")
	}
	return nil
}
