package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func MakeMigrations() {
	connStr := "postgresql://postgres:postgres@localhost?sslmode=disable"
	dbConn, err := GetDBConnection(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	m.Up() // or m.Step(2)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err.Error())
	}
}
