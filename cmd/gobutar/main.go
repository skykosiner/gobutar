package main

// TODO: Add in logging with a verbose option

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skykosiner/gobutar/pkg/budget"
	"github.com/skykosiner/gobutar/pkg/components"
	"github.com/skykosiner/gobutar/pkg/items"
	"github.com/skykosiner/gobutar/pkg/sections"
	"github.com/skykosiner/gobutar/pkg/templates"
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
	`)

	if err != nil {
		slog.Error("Error creating database tables", "error", err)
		return
	}

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
		if len(newName) == 0 {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

		if err := sections.EditSectionName(db, idInt, newName); err != nil {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

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

		templates.RenderTemplate(w, "index", components.Page{
			Budget:   budget,
			Sections: sectionsSlice,
		})
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

		templates.RenderTemplate(w, "index", components.Page{
			Budget:   budget,
			Sections: sectionsSlice,
		})
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/item/allocate", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Please provide an ID.", http.StatusBadRequest)
			return
		}

		var ammountToAllocate struct {
			AmmountToAllocate float64 `json:"ammountToAllocate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&ammountToAllocate); err != nil {
			slog.Error("Error getting body of ammount to alocate.", "error", err)
			return
		}

		if err := items.AllocateMoneyForItem(id, ammountToAllocate.AmmountToAllocate, db); err != nil {
			slog.Error("Seems to be an error", "error", err)
			http.Error(w, "Sorry there was a problem", http.StatusBadRequest)
			return
		}

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

		templates.RenderTemplate(w, "index", components.Page{
			Budget:   budget,
			Sections: sectionsSlice,
		})
	})

	http.HandleFunc("/api/transaction/new", transactions.NewTransaction(db))

	http.HandleFunc("/api/item/new", func(w http.ResponseWriter, r *http.Request) {
		var newItem struct {
			Name      string          `json:"name"`
			Price     string          `json:"price"`
			Saved     string          `json:"saved"`
			Recurring items.Recurring `json:"recurring"`
			SectionID string          `json:"section_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			slog.Error("Error getting body of ammount to alocate.", "error", err)
			return
		}

		fmt.Println(newItem)

		price, _ := strconv.ParseFloat(newItem.Price, 64)
		saved, _ := strconv.ParseFloat(newItem.Saved, 64)
		sectionID, _ := strconv.Atoi(newItem.SectionID)

		item := items.NewItem(newItem.Name, price, saved, newItem.Recurring, sectionID)
		if err := items.SaveItem(db, item); err != nil {
			slog.Error("Couldn't save new item", "error", err, "new item", item)
			http.Error(w, "Sorry there was a problem", http.StatusBadRequest)
			return
		}

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

		templates.RenderTemplate(w, "index", components.Page{
			Budget:   budget,
			Sections: sectionsSlice,
		})
	})

	http.Handle("/api/get-form-new-item", sections.SendNewItemForm(db))

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src"))))

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

	home := components.Home(components.Page{
		Budget:   budget,
		Sections: sectionsSlice,
	})

	http.Handle("/", templ.Handler(home))
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		t, err := transactions.GetTransactions(db)
		if err != nil {
			slog.Error("Error getting spent items.", "error", err, "spent", t)
			return
		}

		spentComponent := components.Transactions(t)
		spentComponent.Render(r.Context(), w)
	})

	if err := http.ListenAndServe(":42069", nil); err != nil {
		slog.Error("Error starting webserver", "error", err)
	}
}
