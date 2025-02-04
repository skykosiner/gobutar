package main

// TODO: Add in logging with a verbose option

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skykosiner/gobutar/pkg/budget"
	"github.com/skykosiner/gobutar/pkg/components"
	"github.com/skykosiner/gobutar/pkg/items"
	"github.com/skykosiner/gobutar/pkg/sections"
	"github.com/skykosiner/gobutar/pkg/transactions"
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
		saved REAL NOT NULL DEFAULT 0.00,
		recurring TEXT NOT NULL CHECK (recurring IN ('no', 'monthly', 'weekly', 'yearly', 'daily')),
		section_id INTEGER NOT NULL,
		FOREIGN KEY (section_id) REFERENCES sections (id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS budget (
		allocated REAL NOT NULL DEFAULT 0.00,
		unallocated REAL NOT NULL DEFAULT 0.00,
		current_balance REAL NOT NULL DEFAULT 0.00,
		all_time_spent REAL NOT NULL DEFAULT 0.00,
		all_time_saved REAL NOT NULL DEFAULT 0.00
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		purchase_date TEXT NOT NULL,
		payee TEXT NOT NULL,
		outflow REAL NOT NULL,
		inflow REAL NOT NULL,
		item_id TEXT NOT NULL,
		FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		color TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS payees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	`)

	if err != nil {
		slog.Error("Error creating database tables", "error", err)
		return
	}

	http.HandleFunc("/api/item/allocate", items.AllocateItemRoute(db))
	http.HandleFunc("/api/item/new", items.NewItemRoute(db))

	http.HandleFunc("/api/section/new-color", sections.SectionNewColor(db))

	http.HandleFunc("/api/transaction/new", transactions.NewTransaction(db))
	http.HandleFunc("/api/transaction/delete", transactions.DeleteTransaction(db))
	http.HandleFunc("/api/transaction/new-form", transactions.SendNewTransactionForm(db))

	http.Handle("/api/get-form-new-item", sections.SendNewItemForm(db))

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sectionsSlice, err := sections.GetSections(db)
		if err != nil {
			slog.Error("Error getting sections", "error", err)
			return
		}

		b, err := budget.NewBudget(db)
		if err != nil {
			slog.Error("Error getting budget", "error", err)
			return
		}

		home := components.Home(components.Page{
			Budget:   b,
			Sections: sectionsSlice,
		})
		home.Render(r.Context(), w)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		t, err := transactions.GetTransactions(db)
		if err != nil {
			slog.Error("Error getting spent items.", "error", err, "spent", t)
			return
		}

		b, err := budget.NewBudget(db)
		if err != nil {
			slog.Error("Error getting budget.", "error", err, "budget", b)
			return
		}
		spentComponent := components.Transactions(t, b.CurrentBalance)
		spentComponent.Render(r.Context(), w)
	})

	if err := http.ListenAndServe(":42069", nil); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
