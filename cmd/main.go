package main

import (
	"context"
	"github.com/Hamifthi/authentication_microservice/internal"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication"
	"github.com/Hamifthi/authentication_microservice/pkg/authentication/adapters"
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
	// TODO postgres db initialization and implementation
	//// initialize, connecting and migrating user model to the database
	//db, err := internal.InitializeAndConnectDBAndMigrate(l)
	//if err != nil {
	//	l.Printf("[Error] got the %s database error", err)
	//  os.Exit(1)
	//}
	//// create the db service and auth service
	//dbService := database.New(db, l)

	// TODO mongodb initialization
	ctx := context.TODO()
	client, err := internal.ConnectMongoDB(ctx, l)
	if err != nil {
		l.Printf("[Error] got the %s mognodb database error", err)
		os.Exit(1)
	}
	dbName, _ := internal.GetEnv("MONGO_DATABASE")
	collName, _ := internal.GetEnv("MONGO_COLLECTION")
	collection := client.Database(dbName).Collection(collName)
	dbService := database.NewMongoSrv(collection, ctx, l)
	authService := authentication.New(dbService, l)

	// TODO this is the rest part of the code
	// create auth handler to use its functions in the router
	authHandler := adapters.NewHandler(authService, l)
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
	shutdownCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(shutdownCtx)

	// TODO enable this if you want to run the grpc part of the file
	//// create a new gRPC server, use WithInsecure to allow http connections
	//gs := grpc.NewServer()
	//
	//// create an instance of the Currency server
	//grpcAS := authentication.NewAuthServer(authService, l)
	//
	//// register the currency server
	//protos.RegisterAuthServiceServer(gs, grpcAS)
	//
	//// register the reflection service which allows clients to determine the methods
	//// for this gRPC service
	//reflection.Register(gs)
	//
	//// create a TCP socket for inbound server connections
	//listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	//if err != nil {
	//	l.Printf("[Error] Unable to create listener due to %s error", err)
	//	os.Exit(1)
	//}
	//
	//// listen for requests
	//gs.Serve(listener)
	// TODO graphql implementation
	//port, _ := internal.GetEnv("PORT")
	//if port == "" {
	//	port = "8000"
	//}
	//resolver := &authentication.Resolver{
	//	authService, l,
	//}
	//
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//l.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//l.Fatal(http.ListenAndServe(":"+port, nil))
}
