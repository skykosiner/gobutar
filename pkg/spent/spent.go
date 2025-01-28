package spent

import (
	"database/sql"

	"github.com/skykosiner/gobutar/pkg/items"
)

type Spent struct {
	ID            int    `sql:"id"`
	PurchaseDate  string `sql:"purchase_date"`
	ItemID        string `sql:"item_id"`
	ItemName      string
	ItemPrice     float64
	ItemRecurring items.Recurring
}

func GetSpentItems(db *sql.DB) ([]Spent, error) {
	var items []Spent

	rows, err := db.Query(`
		SELECT
			s.id AS spent_id,
			s.purchase_date AS purchase_date,
			s.item_id AS spent_item_id,
			i.name AS item_name,
			i.price AS item_price,
			i.recurring AS item_recurring
		FROM
			spent s
		LEFT JOIN
			items i ON s.item_id = i.id;
		`)

	if err != nil {
		return items, err
	}

	defer rows.Close()

	for rows.Next() {
		var item Spent
		if err := rows.Scan(&item.ID, &item.PurchaseDate, &item.ItemID, &item.ItemName, &item.ItemPrice, &item.ItemRecurring); err != nil {
			return items, err
		}

		items = append(items, item)
	}

	return items, nil
}
