package internal

import (
	"fmt"
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
	timezone, err := GetEnv("POSTGRES_TIMEZONE")
	if err != nil {
		logger.Println("[Error] reading database timezone environment variable")
		return "", errors.Wrap(err, "Error reading database timezone environment variable")
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, db, port, ssl, timezone), nil
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

func AutoMigrate(db *gorm.DB, model *interface{}) error {
	err := db.AutoMigrate(model)
	return err
}
