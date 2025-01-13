package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skykosiner/gobutar/pkg/sections"
)

func main() {
	db, err := sql.Open("sqlite3", "./gobutar.db")
	if err != nil {
		slog.Error("Error connecting to DB", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS items (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		recurring TEXT NOT NULL CHECK (recurring IN ('no', 'monthly', 'weekly', 'yearly', 'daily')),
		section_id INTEGER NOT NULL,
		FOREIGN KEY (section_id) REFERENCES sections (id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS budget (
		current_balance REAL NOT NULL DEFAULT 0.00,
		all_time_spent REAL NOT NULL DEFAULT 0.00,
		all_time_saved REAL NOT NULL DEFAULT 0.00
	);

	CREATE TABLE IF NOT EXISTS spent (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		purchase_date TEXT NOT NULL,
		item_id TEXT NOT NULL,
		FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		color TEXT NOT NULL
	);
	`)

	if err != nil {
		slog.Error("Error creating database tables", "error", err)
		return
	}

	sections, _ := sections.GetSections(db)

	for _, section := range sections {
		fmt.Println(section)
	}

	fs := http.FileServer(http.Dir("./src"))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":42069", nil); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
