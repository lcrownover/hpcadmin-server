package data

import (
	"database/sql"
	"fmt"
	"time"
)

type APIKeyToken struct {
	Key       string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetAPIKeyToken(db *sql.DB, key string) (*APIKeyToken, error) {
	var k APIKeyToken
	err := db.QueryRow("SELECT key, role, created_at, updated_at FROM apikeys WHERE key = $1", key).Scan(&k.Key, &k.Role, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no APIKeyToken found with provided key")
		}
		return nil, err
	}
	return &k, err
}
