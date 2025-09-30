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
	"github.com/skykosiner/gobutar/pkg/user"
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

	CREATE TABLE IF NOT EXISTS user (
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		currency TEXT NOT NULL
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

	user.CheckFirstTime(db)

	privateMux := http.NewServeMux()
	publicMux := http.NewServeMux()

	publicMux.HandleFunc("/api/user/set-currency", user.SetCurrency())
	publicMux.HandleFunc("/api/user/sign-up", user.NewUser(db))
	publicMux.HandleFunc("/api/user/login", user.Login(db))

	privateMux.HandleFunc("/api/user/logout", user.Logout())

	privateMux.HandleFunc("/api/item/allocate", items.AllocateItemRoute(db))
	privateMux.HandleFunc("/api/item/new", items.NewItemRoute(db))

	privateMux.HandleFunc("/api/section/new-color", sections.SectionNewColor(db))

	privateMux.HandleFunc("/api/transaction/new", transactions.NewTransaction(db))
	privateMux.HandleFunc("/api/transaction/delete", transactions.DeleteTransaction(db))
	privateMux.HandleFunc("/api/transaction/new-form", transactions.SendNewTransactionForm(db))

	privateMux.HandleFunc("/api/get-form-new-item", sections.SendNewItemForm(db))

	publicMux.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src"))))

	privateMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	privateMux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
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

	finalMux := http.NewServeMux()

	finalMux.Handle("/api/user/set-currency", publicMux)
	finalMux.Handle("/api/user/sign-up", publicMux)
	finalMux.Handle("/api/user/login", publicMux)

	finalMux.Handle("/src/", publicMux)
	finalMux.Handle("/", user.FirstTime(db, user.IsUserLoggedIn(privateMux)))
	// finalMux.Handle("/", privateMux)

	if err := http.ListenAndServe(":42069", finalMux); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
