package controller

import (
	"html/template"
	"net/http"
)

var (
	homeController   home
	searchController search
)

// StartUp is the main controller of the app.
// Startup is handling static assets and setting up subcontrollers.
func StartUp(templatesMap map[string]*template.Template) {
	homeController.homeTemplate = templatesMap["home.html"]
	searchController.searchTemplate = templatesMap["search.html"]
	homeController.registerRoutes()
	searchController.registerRoutes()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
