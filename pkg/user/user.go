package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

func hashPassword(password string) string {
	bcrypt.GenerateFromPassword()
}

func NewUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUserReq struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(newUserReq); err != nil {
			slog.Error("Error getting new user json request.", "error", err, "r", r)
			http.Error(w, "Error cerating user.", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (name, password) VALUES (?,?)")
	}
}
