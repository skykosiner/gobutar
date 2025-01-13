package budget

import (
	"database/sql"
)

type Budget struct {
	DB              *sql.DB
	AllTimeSpent    int
	AllTimeSaved    int
	CurrentBallance int
}

func NewBudget(db *sql.DB) (Budget, error) {
	var allTimeSpent, allTimeSaved, currentBallance int
	rows, err := db.Query("GET all_time_spent, all_time_saved, current_ballance FROM buget")
	if err != nil {
		return Budget{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&allTimeSpent, &allTimeSaved, &currentBallance); err != nil {
			return Budget{}, err
		}
	}

	return Budget{
		DB:              db,
		AllTimeSpent:    allTimeSpent,
		AllTimeSaved:    allTimeSaved,
		CurrentBallance: currentBallance,
	}, nil
}
