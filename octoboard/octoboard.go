package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseFiles(filepath.Join("octoboard", "tmpl", "index.html")))

type Person struct {
	Name string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok || len(names[0]) < 1 {
		log.Println("Url Parameter 'name' is missing")
		return
	}
	name := names[0]
	err := templates.ExecuteTemplate(w, "index.html", Person{name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
