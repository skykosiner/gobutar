package budget

import (
	"database/sql"
)

type Budget struct {
	Unallocated    float64 `sql:"unallocated"`
	Allocated      float64 `sql:"allocated"`
	CurrentBalance float64 `sql:"current_balance"`
	AllTimeSpent   float64 `sql:"all_time_spent"`
	AllTimeSaved   float64 `sql:"all_time_saved"`
}

func NewBudget(db *sql.DB) (Budget, error) {
	rows, err := db.Query("SELECT * FROM budget")
	if err != nil {
		return Budget{}, err
	}

	defer rows.Close()

	var unallocated, allocated, currentBalance, allTimeSpent, allTimeSaved float64
	for rows.Next() {
		if err := rows.Scan(&currentBalance, &allTimeSpent, &allTimeSaved, &allocated, &unallocated); err != nil {
			return Budget{}, err
		}
	}

	return Budget{
		unallocated,
		allocated,
		currentBalance,
		allTimeSpent,
		allTimeSaved,
	}, nil
}

func (b *Budget) SetUnallocated(db *sql.DB, newAmount float64) error {
	query := "UPDATE budget SET allocated = allocated + ?, unallocated = unallocated - ?"
	_, err := db.Exec(query, newAmount, newAmount)
	return err
}
