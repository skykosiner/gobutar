package sections

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

func SectionNewColor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := EditSectionColor(db, id, newColor.NewColor); err != nil {
			http.Error(w, "Please provide a new name.", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
