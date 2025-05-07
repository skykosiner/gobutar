package user

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/skykosiner/gobutar/pkg/components"
)

// If the user isn't logged in this will return the login page
func IsUserLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("JWT")
		if err != nil {
			components.Login().Render(r.Context(), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var isFirstTime bool
var checkedFirstTime bool

func CheckFirstTime(db *sql.DB) {
	if checkedFirstTime {
		return
	}
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user);")
	var exists int
	if err := row.Scan(&exists); err != nil {
		slog.Error("Error checking for first-time use", "error", err)
		isFirstTime = true
	} else {
		isFirstTime = exists == 0
	}
	checkedFirstTime = true
}

func FirstTime(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isFirstTime {
			components.Introduction().Render(r.Context(), w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
