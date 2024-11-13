package handler

import (
	"html/template"
	"net/http"
)

type Index struct{}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/index.html"))

	tmpl.Execute(w, nil)
}
