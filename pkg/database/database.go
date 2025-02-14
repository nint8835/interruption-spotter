package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/nint8835/interruption-spotter/pkg/config"
)

//go:generate go tool sqlc generate

//go:embed schema.sql
var schema string

func Connect(cfg *config.Config) (*Queries, error) {
	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return New(db), nil
}
