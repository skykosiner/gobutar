package payee

import "database/sql"

type Payee struct {
	ID   int    `sql:"id"`
	Name string `sql:"name"`
}

func GetPayees(db *sql.DB) ([]Payee, error) {
	var payees []Payee
	rows, err := db.Query("SELECT * FROM payees")
	if err != nil {
		return payees, err
	}

	defer rows.Close()

	for rows.Next() {
		var payee Payee
		if err := rows.Scan(&payee.ID, &payee.Name); err != nil {
			return payees, err
		}

		payees = append(payees, payee)
	}

	return payees, nil
}

func ProcessPayee(db *sql.DB, name string) error {
	// Check if the payee exists in the db if it doesn't then create it in the
	// payee table
	_, err := db.Query("SELECT name FROM payees WHERE name = ?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		// It doesn't exist so we should create it in the db
		_, err := db.Exec("INSERT INTO payees (name) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	return nil
}
