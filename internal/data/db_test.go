package data

import (
	"os"
	"testing"
    "log"
)

func TestMain(m *testing.M) {
    dbr, _ := NewDBRequest("localhost", 5432, "postgres", "postgres", "hpcadmin_test", true)
	db, err := NewDBConn(dbr)
	if err != nil {
		log.Fatalln(err)
	}
	WipeDB(db)

	code := m.Run()
	db.Close()

	os.Exit(code)
}

func TestNewDBConn(t *testing.T) {
	dbr := DBRequest{
		Host:       "localhost",
		Port:       5432,
		User:       "postgres",
		Password:   "postgres",
		DisableSSL: true,
	}
	db, err := NewDBConn(dbr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

