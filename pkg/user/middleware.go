package user

import (
	"database/sql"
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

func FirstTime(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users;")
		if err != nil || !rows.Next() {
			components.Introduction().Render(r.Context(), w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
