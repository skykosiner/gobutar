package main

type ItemType int

const (
	OneTime ItemType = iota
	Recurring
)

type Item struct {
	RecurringDate int
	Price         float64
	Type          ItemType
	Name          string
}
