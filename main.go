package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"todo/domain"
	"todo/handlers"
	"todo/postgres"
)

func main() {

	// Don't load environment variables if in production.
	appEnvironment := os.Getenv("APP_ENV")

	if appEnvironment != "production" {
		// Load the environment variables from .env file if this is not a production
		// application environment.
		godotenv.Load()
	}

	// Required environment variables.
	postgresUser := os.Getenv("PSQL_USER")
	postgresPassword := os.Getenv("PSQL_PASSWORD")
	postgresAddress := os.Getenv("PSQL_ADDRESS")
	databaseName := os.Getenv("PSQL_DB_NAME")
	port := os.Getenv("APP_PORT")

	if postgresUser == "" || postgresPassword == "" || databaseName == "" || port == "" {
		log.Fatalf("Required environment variables not set.")
	}

	// Database connection pointer.
	DB := postgres.New(&pg.Options{
		User:     postgresUser,
		Password: postgresPassword,
		Addr:     postgresAddress,
		Database: databaseName,
		PoolSize: 20,
	})

	// After the main function returns, this will be called.
	defer DB.Close()

	// This is how we inject the domain package requirements from postgres package.
	// We have to get the connection for the database, but we have to be careful of circular dependencies
	// for our domain and postgres package.
	// So we have to inject the Database in our domain package. Our injection will happen in domain.DB struct,
	// which contains Repository interfaces.
	// Repository interfaces (database gateways) are defined inside our domain package but are implemented
	// inside the postgres package which contains Repository struct (not an interface) containing DB pointer. Mind blown.
	// This way only postgres package operates the DB exposing the database gateways over the
	// defined interface of domain package.
	domainDB := domain.DB{
		UserRepository: postgres.NewUserRepository(DB),
		TodoRepository: postgres.NewTodoRepository(DB),
	}
	// Now the domain includes everything we need for our REST endpoints.
	domain := domain.Domain{
		DB: domainDB,
	}


	// Create our database schema with this options.
	options := orm.CreateTableOptions{
		IfNotExists: true,
	}
	err := postgres.CreateSchema(DB, &options)
	if err != nil {
		log.Fatalf("Unable to create database schema:\n%v", err)
	}

	// Setting up the router via handlers package whilst injecting it with our domain/business logic.
	router := handlers.SetupRouter(&domain)

	// Start our Http server and register our router.
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		// Log if it http does not start listening and serving.
		log.Fatalf("Can not start server %v", err)
	}
}
