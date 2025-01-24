package items

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/skykosiner/gobutar/pkg/budget"
)

type Recurring string

const (
	No      Recurring = "no"
	Daily             = "daily"
	Weekly            = "weekly"
	Monthly           = "monthly"
	Yearly            = "yearly"
)

type Item struct {
	ID        string    `json:"id" sql:"id"`
	Name      string    `json:"name" sql:"name"`
	Price     float64   `json:"price" sql:"price"`
	Saved     float64   `json:"saved" sql:"saved"`
	Recurring Recurring `json:"recurring" sql:"recurring"`
	SectionID int       `json:"section_id" sql:"section_id"`
}

func NewItem(name string, price, saved float64, recurring Recurring, sectionID int) Item {
	return Item{
		ID:        uuid.New().String(),
		Name:      name,
		Price:     price,
		Saved:     saved,
		Recurring: recurring,
		SectionID: sectionID,
	}
}

func (i Item) String() string {
	return fmt.Sprintf("ID: %s\nName: %s\nPrice: %f\nRecuring: %s, SectionID: %d", i.ID, i.Name, i.Price, i.Recurring, i.SectionID)
}


func ParseRecurring(value string) (Recurring, error) {
	switch value {
	case string(No), string(Monthly), string(Weekly), string(Yearly), string(Daily):
		return Recurring(value), nil
	default:
		return "", fmt.Errorf("invalid recurring value: %s", value)
	}
}

func AllocateMoneyForItem(itemID string, ammountToAlocate float64, db *sql.DB) error {
	b, err := budget.NewBudget(db)
	if err != nil {
		return err
	}

	if ammountToAlocate > b.Unallocated {
		return fmt.Errorf("You can't allocate money you don't have")
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE items SET saved = saved + %.2f WHERE id = '%s'", ammountToAlocate, itemID))
	if err != nil {
		return err
	}

	return b.SetUnallocated(db, ammountToAlocate)
}

func SaveItem(db *sql.DB, item Item) error {
	query := "INSERT INTO items (id, name, price, saved, recurring, section_id) VALUES (?,?,?,?,?,?)"
	_, err := db.Exec(query, item.ID, item.Name, item.Price, item.Saved, item.Recurring, item.SectionID)
	return err
}
