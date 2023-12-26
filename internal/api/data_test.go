package api

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/lcrownover/hpcadmin-server/internal/data"
)

type testDataHandler struct {
	DB *sql.DB
}

func NewTestDataHandler() *testDataHandler {
	host, found := os.LookupEnv("HPCADMIN_TEST_DATABASE_HOST")
	if !found {
		panic("HPCADMIN_TEST_DATABASE_HOST not set")
	}
	portStr, found := os.LookupEnv("HPCADMIN_TEST_DATABASE_PORT")
	if !found {
		panic("HPCADMIN_TEST_DATABASE_PORT not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("HPCADMIN_TEST_DATABASE_PORT not an integer")
	}
	user, found := os.LookupEnv("HPCADMIN_TEST_DATABASE_USERNAME")
	if !found {
		panic("HPCADMIN_TEST_DATABASE_USERNAME not set")
	}
	password, found := os.LookupEnv("HPCADMIN_TEST_DATABASE_PASSWORD")
	if !found {
		panic("HPCADMIN_TEST_DATABASE_PASSWORD not set")
	}
	dbname, found := os.LookupEnv("HPCADMIN_TEST_DATABASE_NAME")
	if !found {
		panic("HPCADMIN_TEST_DATABASE_NAME not set")
	}
	dbr := data.DBRequest{
		Host:       host,
		Port:       port,
		User:       user,
		Password:   password,
		DBName:     dbname,
		DisableSSL: true,
	}
	db, err := data.NewDBConn(dbr)
	if err != nil {
		log.Fatal(err)
	}
	return &testDataHandler{
		DB: db,
	}
}
