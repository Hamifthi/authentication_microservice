package main

import (
	"github.com/Hamifthi/authentication_microservice/internal"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "Auth-Service ", log.LstdFlags)
	// initialize environment variable
	err := internal.InitializeEnv(".env")
	if err != nil {
		l.Println("[Error] reading database user environment variable")
		os.Exit(1)
	}
	// initialize, connecting and migrating user model to the database
	db, err := internal.InitializeAndConnectDBAndMigrate(l)
	if err != nil {
		l.Printf("[Error] got the %s database error", err)
	}
}
