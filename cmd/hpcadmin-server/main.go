package main

import (
	"flag"
	"log"

	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/db"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	connStr := "postgresql://postgres:postgres@localhost/hpcadmin"
	dbConn, err := db.GetDBConnection(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	api.Run(dbConn)
}
