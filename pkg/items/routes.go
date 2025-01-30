package items

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func AllocateItemRoute(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := AllocateMoneyForItem(id, ammountToAllocate.AmmountToAllocate, db); err != nil {
			slog.Error("Seems to be an error", "error", err)
			http.Error(w, "Sorry there was a problem", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewItemRoute(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newItem struct {
			Name      string          `json:"name"`
			Price     string          `json:"price"`
			Saved     string          `json:"saved"`
			Recurring Recurring `json:"recurring"`
			SectionID string          `json:"section_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			slog.Error("Error getting body of ammount to alocate.", "error", err)
			return
		}

		price, _ := strconv.ParseFloat(newItem.Price, 64)
		saved, _ := strconv.ParseFloat(newItem.Saved, 64)
		sectionID, _ := strconv.Atoi(newItem.SectionID)

		item := NewItem(newItem.Name, price, saved, newItem.Recurring, sectionID)
		if err := SaveItem(db, item); err != nil {
			slog.Error("Couldn't save new item", "error", err, "new item", item)
			http.Error(w, "Sorry there was a problem", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
