package data

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DBRequest struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	DisableSSL bool
}

func NewDBRequest(host string, port int, user, password, dbname string, disableSSL bool) (DBRequest, error) {
	return DBRequest{
		Host:       host,
		Port:       port,
		User:       user,
		Password:   password,
		DBName:     dbname,
		DisableSSL: disableSSL,
	}, nil
}

func NewDBConn(dbr DBRequest) (*sql.DB, error) {
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
