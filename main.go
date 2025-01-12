package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./gobutar.db")
	if err != nil {
		slog.Error("Error connecting to DB", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	databaseCreateLogger := slog.New(slog.NewLogLogger(slog.Default().Handler().WithGroup("database creation"), slog.LevelDebug)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS budget (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)
	`)

	if err != nil {
	}
}
