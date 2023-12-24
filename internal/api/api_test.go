package api

import (
	"database/sql"
	"log"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

type testDataHandler struct {
	db *sql.DB
}

func newTestDataHandler() *testDataHandler {
	dbr := data.DBRequest{
		Host:       "localhost",
		Port:       5432,
		User:       "hpcadmin",
		Password:   "superfancytestpasswordthatnobodyknows&",
		DBName:     "hpcadmin_test",
		DisableSSL: true,
	}
	db, err := data.NewDBConn(dbr)
	if err != nil {
		log.Fatal(err)
	}
	return &testDataHandler{
		db: db,
	}
}
