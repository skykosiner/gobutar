package items

type ItemType int

const (
	OneTime ItemType = iota
	Recurring
)

type Item struct {
	ID            string
	RecurringDate int
	Price         float64
	Type          ItemType
	Name          string
}
