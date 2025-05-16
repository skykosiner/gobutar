package utils

import (
	"fmt"
	"html"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/skykosiner/gobutar/pkg/items"
)

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}

func FormatRecurring(recurring items.Recurring) string {
	if recurring != items.No {
		return strings.Title(string(recurring))
	}

	return "One Time"
}

func SortItems(items []items.Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
}

func HTMXError(w http.ResponseWriter, error string, status int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprintf(w, `<p>%s</p>`, html.EscapeString(error))
}
