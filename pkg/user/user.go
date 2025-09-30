package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/skykosiner/gobutar/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type Currency = string

const (
	GBP Currency = "GBP"
	USD Currency = "USD"
)

type User struct {
	Email    string `sql:"email"`
	Password string `sql:"password"`
	Currency string `sql:"currency"`
}

func (u *User) GetCurrencySymbol() string {
	switch u.Currency {
	case "GBP":
		return "£"
	case "USD":
		return "$"
	}

	u.GetCurrencySymbol()

	return "$"
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
			utils.HTMXError(w, "Error creating user.", http.StatusBadRequest)
			return
		}

		if err := newUserReq.hashPassword(); err != nil {
			slog.Error("Error hashing user password.", "error", err, "new user", newUserReq)
			utils.HTMXError(w, "Error hashing your password. Please try again.", http.StatusBadRequest)
			return
		}

		if _, err := db.Exec("INSERT INTO user (email, password) VALUES (?,?)", newUserReq.Email, newUserReq.Password); err != nil {
			slog.Error("Error creating new user.", "error", err, "new user", newUserReq)
			utils.HTMXError(w, "Sorry there was an error please try again.", http.StatusBadRequest)
			return
		}

		jwtToken, err := CreateJWT(newUserReq.Email)
		if err != nil {
			slog.Error("Erorr creating JWT token.", "error", err)
			utils.HTMXError(w, "Error logging you in. Please try again.", http.StatusInternalServerError)
			return
		}

		isFirstTime = false

		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    jwtToken,
			Expires:  time.Now().AddDate(0, 0, 90),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest User

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			slog.Error("Error getting JSON request to login.", "error", err, "r", r)
			utils.HTMXError(w, "Error decoding your login request. Please try again.", http.StatusBadRequest)

			return
		}

		// There is only one user so this is fine, this is just self hosted so
		// it doesn't need to support loads of users
		rows, err := db.Query("SELECT * FROM user;")
		if err != nil {
			slog.Error("Error finding user in the database.", "error", err, "user info", loginRequest)
			utils.HTMXError(w, "Error please try again.", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		if !rows.Next() {
			utils.HTMXError(w, "Please make sure that your user account exists.", http.StatusBadRequest)
			return
		}

		var dbUser User
		if err := rows.Scan(&dbUser.Email, &dbUser.Password); err != nil {
			slog.Error("Error getting your user from the database.", "error", err)
			utils.HTMXError(w, "Error getting your info from the database. Please try again.", http.StatusInternalServerError)
			return
		}

		if loginRequest.Email != dbUser.Email {
			utils.HTMXError(w, "The email you entered isn't correct.", http.StatusBadRequest)
			return
		}

		if !loginRequest.validPassword(dbUser.Password) {
			utils.HTMXError(w, "Your password is incorrect.", http.StatusUnauthorized)
			return
		}

		// Just do a cookie I guess, and add middlewear we can call to check if
		// the user is logged in
		jwtToken, err := CreateJWT(loginRequest.Email)
		if err != nil {
			slog.Error("Erorr creating JWT token.", "error", err)
			utils.HTMXError(w, "Error logging you in. Please try again.", http.StatusInternalServerError)
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

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func SetCurrency() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var currency Currency
		currency := r.FormValue("currency")

		fmt.Fprintln(w, "Successfully Set Currency")
		w.WriteHeader(http.StatusOK)
	}
}
