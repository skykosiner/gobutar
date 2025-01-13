package spend

import (
	"database/sql"
	"fmt"

	"github.com/skykosiner/gobutar/pkg/items"
)

type Spend struct {
	Item items.Item
	Date string
}

func GetDateSpent(day, month, year int, db *sql.DB) ([]Spend, error) {
	var spent []Spend
	rows, err := db.Query(fmt.Sprintf("GET * FROM spent WHERE date is = '%d/%d/%d'", day, month, year))
	if err != nil {
		return spent, err
	}

	defer rows.Close()

	for rows.Next() {
		var item Spend
		if err := rows.Scan(&item.Item.ID, &item.Date, &item.Item.Type, &item.Item.Name, &item.Item.Price, &item.Item.RecurringDate); err != nil {
			return spent, err
		}

		spent = append(spent, item)
	}

	return spent, nil
}
