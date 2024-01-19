package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type APIKeyEntry struct {
	Key       string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetAPIKeyEntry looks for the provided key in the database
// and returns the APIKeyEntry if found, nil if not found, or an error
func GetAPIKeyEntry(db *sql.DB, key string) (*APIKeyEntry, error) {
	slog.Debug("querying database for api key", "package", "data", "method", "GetAPIKeyEntry")
	var k APIKeyEntry
	err := db.QueryRow("SELECT key, role, created_at, updated_at FROM api_keys WHERE key = $1", key).Scan(&k.Key, &k.Role, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("api key not found in database", "package", "data", "method", "GetAPIKeyEntry")
			return nil, fmt.Errorf("no APIKeyToken found with provided key")
		}
		slog.Debug("failed to look up key from database", "package", "data", "method", "GetAPIKeyEntry", "error", err)
		return nil, err
	}
	slog.Debug("found api key in database", "package", "data", "method", "GetAPIKeyEntry")
	return &k, err
}
