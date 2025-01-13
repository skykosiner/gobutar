package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skykosiner/gobutar/pkg/budget"
	"github.com/skykosiner/gobutar/pkg/sections"
)

type Page struct {
	Budget   budget.Budget
	Sections []sections.Section
}

var (
	templates = template.Must(template.ParseGlob("src/*.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		slog.Error("Error rendering template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
		saved REAL NOT NULL DEFAULT 0.00,
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sectionsSlice, err := sections.GetSections(db)
		if err != nil {
			slog.Error("Error getting sections", "error", err)
			return
		}

		budget, err := budget.NewBudget(db)
		if err != nil {
			slog.Error("Error getting budget", "error", err)
			return
		}

		renderTemplate(w, "index", Page{
			Budget:   budget,
			Sections: sectionsSlice,
		})
	})

	http.HandleFunc("/api/section/new-name", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Please provide an ID.", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Please provide a valid ID.", http.StatusBadRequest)
			return
		}

		newName := r.FormValue("newName")
		fmt.Println(newName)
		if len(newName) == 0 {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

		if err := sections.EditSectionName(db, idInt, newName); err != nil {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/section/new-color", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Please provide an ID.", http.StatusBadRequest)
			return
		}

		var newColor struct {
			NewColor string `json:"newColor"`
		}
		if err := json.NewDecoder(r.Body).Decode(&newColor); err != nil {
			slog.Error("Error getting body of new color", "error", err)
			return
		}

		if err := sections.EditSectionColor(db, id, newColor.NewColor); err != nil {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src"))))

	if err := http.ListenAndServe(":42069", nil); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
