package sections

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/skykosiner/gobutar/pkg/budget"
	"github.com/skykosiner/gobutar/pkg/items"
	"github.com/skykosiner/gobutar/pkg/templates"
)

type Section struct {
	ID    int          `json:"id" sql:"id"`
	Name  string       `json:"name" sql:"name"`
	Color string       `json:"color" sql:"color"`
	Items []items.Item `json:"item"`
}

func (s Section) String() string {
	return fmt.Sprintf("ID: %d\nName: %s\nColor: %s\nItems: %v", s.ID, s.Name, s.Color, s.Items)
}

func GetSections(db *sql.DB) ([]Section, error) {
	rows, err := db.Query(`
	SELECT
    s.id AS section_id,
    s.name AS section_name,
    s.color AS section_color,
    i.id AS item_id,
    i.name AS item_name,
    i.price AS item_price,
    i.saved AS item_saved,
    i.recurring AS item_recurring,
    i.section_id AS item_section_id
FROM
    sections s
LEFT JOIN
    items i ON s.id = i.section_id;
`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sectionMap := make(map[int]*Section)
	for rows.Next() {
		var sectionID int
		var sectionName, sectionColor string
		var itemID, itemName, itemRecurring sql.NullString
		var itemPrice, itemSaved sql.NullFloat64
		var itemSectionID sql.NullInt64

		err := rows.Scan(
			&sectionID, &sectionName, &sectionColor,
			&itemID, &itemName, &itemPrice, &itemSaved, &itemRecurring, &itemSectionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := sectionMap[sectionID]; !exists {
			sectionMap[sectionID] = &Section{
				ID:    sectionID,
				Name:  sectionName,
				Color: sectionColor,
				Items: []items.Item{},
			}
		}

		if itemID.Valid {
			recuring, err := items.ParseRecurring(itemRecurring.String)
			if err != nil {
				return nil, fmt.Errorf("invalid recurring value for item: %w", err)
			}

			sectionMap[sectionID].Items = append(sectionMap[sectionID].Items, items.Item{
				ID:        itemID.String,
				Name:      itemName.String,
				Price:     itemPrice.Float64,
				Saved:     itemSaved.Float64,
				Recurring: recuring,
				SectionID: int(itemSectionID.Int64),
			})
		}
	}

	sections := make([]Section, 0, len(sectionMap))
	for _, section := range sectionMap {
		sections = append(sections, *section)
	}

	return sections, nil
}

func EditSectionName(db *sql.DB, sectionID int, name string) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE sections SET name = '%s' WHERE id = %d", name, sectionID))
	return err
}

func EditSectionColor(db *sql.DB, sectionID, color string) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE sections SET color = '%s' WHERE id = %s", color, sectionID))
	return err
}

func SendNewItemForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newFormInfo struct {
			Unallocated float64
			ItemsInfo   []struct {
				Name      string
				ID        string
				Allocated float64
			}
			SectionInfo []struct {
				Name string
				ID   int
			}
		}

		s, err := GetSections(db)
		if err != nil {
			http.Error(w, "Error getting your sections.", http.StatusInternalServerError)
			return
		}

		for _, section := range s {
			newFormInfo.SectionInfo = append(newFormInfo.SectionInfo, struct {
				Name string
				ID   int
			}{
				section.Name,
				section.ID,
			})
		}

		b, err := budget.NewBudget(db)
		if err != nil {
			http.Error(w, "Error getting your budget.", http.StatusInternalServerError)
			return
		}

		newFormInfo.Unallocated = b.Unallocated

		rows, err := db.Query("SELECT name, id, saved FROM items WHERE saved > 0.0")
		if err != nil {
			http.Error(w, "Error getting item info.", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var name, id string
			var allocated float64

			if err := rows.Scan(&name, &id, &allocated); err != nil {
				http.Error(w, "Error getting item info.", http.StatusInternalServerError)
				return
			}

			newFormInfo.ItemsInfo = append(newFormInfo.ItemsInfo, struct {
				Name      string
				ID        string
				Allocated float64
			}{
				name,
				id,
				allocated,
			})
		}

		templates.RenderTemplate(w, "new-item", newFormInfo)
	}
}
