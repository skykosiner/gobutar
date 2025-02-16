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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func validPassword(inputPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil

}

func NewUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUserReq struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newUserReq); err != nil {
			slog.Error("Error getting new user json request.", "error", err, "r", r)
			http.Error(w, "Error cerating user.", http.StatusBadRequest)
			return
		}

		hashedPass, err := hashPassword(newUserReq.Password)
		if err != nil {
			slog.Error("Error hashing user password.", "error", err, "new user", newUserReq)
			http.Error(w, "Error hashing your password. Please try again.", http.StatusBadRequest)
			return
		}

		if _, err = db.Exec("INSERT INTO users (email, password) VALUES (?,?)", newUserReq.Email, hashedPass); err != nil {
			slog.Error("Error creating new user.", "error", err, "new user", newUserReq)
			http.Error(w, "Sorry there was an error please try again.", http.StatusBadRequest)
			return
		}
	}
}
