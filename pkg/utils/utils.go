package utils

import (
	"strconv"
	"strings"

	"github.com/skykosiner/gobutar/pkg/items"
)

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64) // 'f' for fixed-point, 2 for 2 decimal places
}

func FormatRecurring(recurring items.Recurring) string {
	if recurring != items.No {
		return strings.Title(string(recurring))
	}

	return "One Time"
}
