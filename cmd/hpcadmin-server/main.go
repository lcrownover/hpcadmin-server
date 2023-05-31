package main

import (
	"flag"
	"log"

	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	db.MakeMigrations()

	connStr := "postgresql://postgres:postgres@localhost/hpcadmin?sslmode=disable"
	dbConn, err := db.GetDBConnection(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	api.Run(dbConn)
}
