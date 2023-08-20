package data

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type DBRequest struct {
	Driver     string
	File       string
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	DisableSSL bool
}

func NewDBConn(dbr DBRequest) (*sql.DB, error) {
	switch dbr.Driver {
	case "postgres":
		db, err := newPostgresDB(dbr)
		if err != nil {
			err = fmt.Errorf("failed to create postgres database: %v", err.Error())
			return nil, err
		}
		return db, nil
	case "sqlite3":
		db, err := newSqliteDB(dbr)
		if err != nil {
			err = fmt.Errorf("failed to create sqlite database: %v", err.Error())
			return nil, err
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbr.Driver)
	}
}

func newPostgresDB(dbr DBRequest) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s", dbr.User, dbr.Password, dbr.Host, dbr.DBName)
	if dbr.DisableSSL {
		connStr = connStr + "?sslmode=disable"
	}
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err.Error())
	}
	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err.Error())
	}
	return dbConn, nil
}

func newSqliteDB(dbr DBRequest) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", dbr.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err.Error())
	}
	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err.Error())
	}
	return dbConn, nil
}

func WipeDB(db *sql.DB) error {
	tables := []string{"pirgs_users", "pirgs_groups", "pirgs_admins", "groups_users", "pirgs", "users"}
	for _, table := range tables {
		q := fmt.Sprintf("DELETE FROM %s", table)
		_, err := db.Exec(q)
		if err != nil {
			return fmt.Errorf("failed to delete %s: %v", table, err.Error())
		}
	}
	return nil
}
