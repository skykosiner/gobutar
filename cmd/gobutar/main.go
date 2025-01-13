package main

import (
	"database/sql"
	"log/slog"
	"net/http"
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

	// _, err = db.Exec(`
	// CREATE TABLE IF NOT EXISTS budget (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT NOT NULL
	// )
	// `)
	//
	// if err != nil {
	// }

	fs := http.FileServer(http.Dir("./src"))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":42069", nil); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
