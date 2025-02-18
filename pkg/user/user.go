package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u User) validPassword(correctPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(u.Password))
	return err == nil
}

func NewUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUserReq User

		if err := json.NewDecoder(r.Body).Decode(&newUserReq); err != nil {
			slog.Error("Error getting new user json request.", "error", err, "r", r)
			http.Error(w, "Error cerating user.", http.StatusBadRequest)
			return
		}

		if err := newUserReq.hashPassword(); err != nil {
			slog.Error("Error hashing user password.", "error", err, "new user", newUserReq)
			http.Error(w, "Error hashing your password. Please try again.", http.StatusBadRequest)
			return
		}

		if _, err := db.Exec("INSERT INTO users (email, password) VALUES (?,?)", newUserReq.Email, newUserReq.Password); err != nil {
			slog.Error("Error creating new user.", "error", err, "new user", newUserReq)
			http.Error(w, "Sorry there was an error please try again.", http.StatusBadRequest)
			return
		}
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest User

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			slog.Error("Error getting JSON request to login.", "error", err, "r", r)
			http.Error(w, "Error decoding your login request. Please try again.", http.StatusBadRequest)

			return
		}

		// There is only one user so this is fine, this is just self hosted so
		// it doesn't need to support loads of users
		rows, err := db.Query("SELECT * FROM user;")
		if err != nil {
			slog.Error("Error finding user in the database.", "error", err, "user info", loginRequest)
			http.Error(w, "Error please try again.", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		if !rows.Next() {
			http.Error(w, "Please make sure that your user account exists.", http.StatusBadRequest)
			return
		}

		var dbUser User
		if err := rows.Scan(&dbUser.Email, &dbUser.Password); err != nil {
			slog.Error("Error getting your user from the database.", "error", err)
			http.Error(w, "Error getting your info from the database. Please try again.", http.StatusInternalServerError)
			return
		}

		if loginRequest.Email != dbUser.Email {
			http.Error(w, "The email you entered isn't correct.", http.StatusBadRequest)
			return
		}

		if !loginRequest.validPassword(dbUser.Password) {
			http.Error(w, "Your password is incorrect.", http.StatusUnauthorized)
			return
		}

		// Just do a cookie I guess, and add middlewear we can call to check if
		// the user is logged in
		jwtToken, err := CreateJWT(loginRequest.Email)
		if err != nil {
			slog.Error("Erorr creating JWT token.", "error", err)
			http.Error(w, "Error logging you in. Please try again.", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    jwtToken,
			Expires:  time.Now().AddDate(0, 0, 90),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})
	}
}
