package items

import "fmt"

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
	Recurring Recurring `json:"recurring" sql:"recurring"`
	SectionID int       `json:"section_id" sql:"section_id"`
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
