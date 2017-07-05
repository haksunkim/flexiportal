package controller

import (
	"html/template"
	"net/http"
	"fmt"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("resources/templates/layout/*")
	t.ParseFiles("resources/templates/home/index.html")
	err := t.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Printf(err.Error(), w)
	}
}
