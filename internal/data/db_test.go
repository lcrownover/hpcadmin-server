package data

import (
	"database/sql"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=hpcadmin sslmode=disable")
	WipeDB(db)

	code := m.Run()
	db.Close()

	os.Exit(code)
}

func TestNewDBConnPostgres(t *testing.T) {
	dbr := DBRequest{
		Driver:     "postgres",
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

func TestNewDBConnSQLite(t *testing.T) {
	dbr := DBRequest{
		Driver:     "sqlite3",
        File:       "hpcadmin.db",
	}
	db, err := NewDBConn(dbr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}
