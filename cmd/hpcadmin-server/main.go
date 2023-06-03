package main

import (
	"context"
	"flag"
	"log"

	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/db"
	"github.com/lcrownover/hpcadmin-server/internal/types"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	connStr := "postgresql://postgres:postgres@localhost/hpcadmin?sslmode=disable"
	dbConn, err := db.GetDBConnection(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, types.DBKey, dbConn)

	api.Run(ctx)
}
