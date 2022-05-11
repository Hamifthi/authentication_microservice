package main

import (
	"context"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication"
	"github.com/Hamifthi/authentication_microservice/pkg/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	// create the db service and auth service
	dbService := database.New(db, l)
	authService := authentication.New(dbService, l)
	// create auth handler to use its functions in the router
	authHandler := authentication.NewHandler(authService, l)
	// create the router
	sm := mux.NewRouter()
	SignUpRouter := sm.Methods(http.MethodPost).Subrouter()
	SignUpRouter.HandleFunc("/signup", authHandler.UserSignUp)
	SignUpRouter.Use(authHandler.MiddlewareValidateUser)

	LoginRouter := sm.Methods(http.MethodPost).Subrouter()
	LoginRouter.HandleFunc("/login", authHandler.UserLogin)
	LoginRouter.Use(authHandler.MiddlewareValidateUser)

	// create a new server
	bindAddress, err := internal.GetEnv("BINDADDRESS")
	if err != nil {
		l.Println("[Error] getting the bindAddress")
		bindAddress = ":8000"
	}
	server := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 8000")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
