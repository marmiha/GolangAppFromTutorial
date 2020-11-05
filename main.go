package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"todo/handlers"
	"todo/postgres"
)

func main() {

	// Don't load environment variables if in production.
	appEnvironment := os.Getenv("APP_ENV")

	if appEnvironment != "production" {
		// Load the environment variables from .env file.
		godotenv.Load()
	}

	// Required environment variables.
	postgresUser := os.Getenv("PSQL_USER")
	postgresPassword := os.Getenv("PSQL_PASSWORD")
	databaseName := os.Getenv("PSQL_DB_NAME")
	port := os.Getenv("APP_PORT")

	if postgresUser == "" || postgresPassword == "" || databaseName == "" || port == "" {
		log.Fatalf("Required environment variables not set.")
	}

	// Database connection pointer.
	DB := postgres.New(&pg.Options{
		User:     postgresUser,
		Password: postgresPassword,
		Database: databaseName,
	})

	// After the main function returns, this will be called.
	defer DB.Close()

	// Setting up the router via handlers package.
	router := handlers.SetupRouter()

	// Start our Http server and register our router.
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		// Log if it http does not start listening and serving.
		log.Fatalf("cannot start server %v", err)
	}
}
