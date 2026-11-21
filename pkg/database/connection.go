package database

import (
	"database/sql"
	"log/slog"
)

const FileName = "file:icas.db"

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", FileName)
	if err != nil {
		slog.Error("error opening database", "error", err)
	}
	return db, err
}
