package templates

import (
	"log/slog"
	"net/http"
	"text/template"

	"github.com/skykosiner/gobutar/pkg/utils"
)
var (
	Templates = template.Must(template.New("base").Funcs(template.FuncMap{
		"formatFloat":     utils.FormatFloat,
		"formatRecurring": utils.FormatRecurring,
	}).ParseGlob("src/*.html"))
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	if err := Templates.ExecuteTemplate(w, tmpl, data); err != nil {
		slog.Error("Error rendering template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
