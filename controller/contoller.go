package controller

import (
	"html/template"
	"net/http"
)

var (
	homeController   home
	searchController search
)

func StartUp(templatesMap map[string]*template.Template) {
	homeController.homeTemplate = templatesMap["home.html"]
	searchController.searchTemplate = templatesMap["search.html"]
	homeController.registerRoutes()
	searchController.registerRoutes()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
