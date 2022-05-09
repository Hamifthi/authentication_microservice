package internal

import (
	"fmt"
	"github.com/Hamifthi/authentication_microservice/entity"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitializeDBCredentials(logger *log.Logger) (string, error) {
	user, err := GetEnv("POSTGRES_USER")
	if err != nil {
		logger.Println("[Error] reading database user environment variable")
		return "", errors.Wrap(err, "Error reading database user environment variable")
	}
	password, err := GetEnv("POSTGRES_PASSWORD")
	if err != nil {
		logger.Println("[Error] reading database password environment variable")
		return "", errors.Wrap(err, "Error reading database password environment variable")
	}
	db, err := GetEnv("POSTGRES_DB")
	if err != nil {
		logger.Println("[Error] reading database db environment variable")
		return "", errors.Wrap(err, "Error reading db environment variable")
	}
	host, err := GetEnv("POSTGRES_HOST")
	if err != nil {
		logger.Println("[Error] reading database host environment variable")
		return "", errors.Wrap(err, "Error reading database host environment variable")
	}
	port, err := GetEnv("POSTGRES_PORT")
	if err != nil {
		logger.Println("[Error] reading database port environment variable")
		return "", errors.Wrap(err, "Error reading database port environment variable")
	}
	ssl, err := GetEnv("POSTGRES_SSL")
	if err != nil {
		logger.Println("[Error] reading database ssl environment variable")
		return "", errors.Wrap(err, "Error reading database ssl environment variable")
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, db, port, ssl), nil
}

func CreateDBConnection(DSN string, logger *log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		logger.Println("[Error] occurred while creating database connection")
		return nil, err
	}

	// Create the connection pool
	sqlDB, err := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, err
}

func GetDatabaseConnection(dbConn *gorm.DB, logger *log.Logger) (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Println("[Error] occurred while connecting with the database")
		return dbConn, err
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Println("[Error] occurred while ping the database")
		return dbConn, err
	}
	return dbConn, nil
}

func AutoMigrate(db *gorm.DB, model interface{}) error {
	err := db.AutoMigrate(model)
	return err
}

func InitializeAndConnectDBAndMigrate(l *log.Logger) (*gorm.DB, error) {
	DSN, err := InitializeDBCredentials(l)
	if err != nil {
		l.Println("[Error] initializing db credentials")
		return nil, errors.Wrap(err, "Error initializing db credentials")
	}
	dbConn, err := CreateDBConnection(DSN, l)
	if err != nil {
		l.Println("[Error] creating database connections")
		return nil, errors.Wrap(err, "Error creating database connections")
	}
	db, err := GetDatabaseConnection(dbConn, l)
	if err != nil {
		l.Println("[Error] cannot get the database connection")
		return nil, errors.Wrap(err, "Error cannot get the database connection")
	}
	err = AutoMigrate(db, entity.User{})
	if err != nil {
		l.Println("[Error] cannot auto migrate the user to the database")
		return nil, errors.Wrap(err, "Error cannot auto migrate the user to the database")
	}
	return db, nil
}
