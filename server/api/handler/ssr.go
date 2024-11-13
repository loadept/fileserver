package handler

import (
	"html/template"
	"net/http"

	"github.com/jsusmachaca/go-router/pkg/response"
)

type Index struct{}
type Gallery struct{}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/index.html")
	if err != nil {
		response.JsonErrorFromString(w, "Error to render template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
