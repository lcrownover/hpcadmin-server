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
