package main

import (
	"database/sql"
	"fmt"
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

func (b Budget) GetDateSpent(day, month, year int) ([]Spend, error) {
	var spent []Spend
	rows, err := b.DB.Query(fmt.Sprintf("GET * FROM spent WHERE date is = '%d/%d/%d'", day, month, year))
	if err != nil {
		return spent, err
	}

	defer rows.Close()

	for rows.Next() {
		var item Spend
		if err := rows.Scan(&item.Date, &item.item.Type, &item.item.Name, &item.item.Price, &item.item.RecurringDate); err != nil {
			return spent, err
		}

		spent = append(spent, item)
	}

	return spent, nil
}
