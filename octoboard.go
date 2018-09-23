package main

import (
	"github.com/radlinskii/octo-board/controller"
	"github.com/radlinskii/octo-board/middleware"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Failed to load variable $PORT from environment")
	}
	templatesMap := populateTemplates()

	controller.StartUp(templatesMap)
	log.Fatal(http.ListenAndServe(":"+port, &middleware.TimeoutMiddleware{Next: new(middleware.GzipMiddleware)}))
}

func populateTemplates() map[string]*template.Template {
	result := make(map[string]*template.Template)
	const basePath = "templates"
	layout := template.Must(template.ParseFiles(filepath.Join(basePath, "_layout.html")))
	dir, err := os.Open(filepath.Join(basePath, "content"))
	if err != nil {
		panic("Failed to open template blocks directory: " + err.Error())
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		panic("Failed to read content directory: " + err.Error())
	}
	for _, fi := range fis {
		fname := fi.Name()
		f, err := os.Open(filepath.Join(basePath, "content", fname))
		if err != nil {
			panic("Failed to open template file: " + fname)
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			panic("Failed to read content from template file: " + fname)
		}
		f.Close()
		tmplt := template.Must(layout.Clone())
		_, err = tmplt.Parse(string(content))
		if err != nil {
			panic("Failed to parse template: " + fname)
		}
		result[fname] = tmplt
	}
	return result
}
