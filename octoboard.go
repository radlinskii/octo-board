package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

type Person struct {
	Name string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok || len(names[0]) < 1 {
		log.Println(r.URL)
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
