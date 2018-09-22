package main

import (
	"context"
	"github.com/radlinskii/octo-board/viewmodel"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
)

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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Failed to load variable $PORT from environment")
	}

	templatesMap := populateTemplates()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handleRoot(templatesMap["index.html"]))
	http.HandleFunc("/search", handleSearch(templatesMap["search.html"]))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func handleRoot(tmplt *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		err := tmplt.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleSearch(tmplt *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "is:open"

		label := buildQuery(&query, r, "label")
		language := buildQuery(&query, r, "language")
		org := buildQuery(&query, r, "org")

		client := github.NewClient(nil)
		issuesPayload, err := fetchIssues(client, query)

		var issues []viewmodel.GithubIssue
		for _, issue := range issuesPayload.Issues {
			issues = append(issues, viewmodel.GithubIssue{
				Title:          issue.GetTitle(),
				Repo:           getRepositoryFullName(issue.GetRepositoryURL()),
				HTMLURL:        issue.GetHTMLURL(),
				Number:         issue.GetNumber(),
				Body:           issue.GetBody(),
				CommentsNumber: issue.GetComments(),
			})
		}

		err = tmplt.Execute(w, viewmodel.Content{Issues: issues, Label: label, Organization: org, Language: language})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getRepositoryFullName(url string) string {
	return strings.Replace(url, "https://api.github.com/repos/", "", 1)
}

func buildQuery(query *string, r *http.Request, opt string) string {
	queryParameters, ok := r.URL.Query()[opt]
	if !ok || len(queryParameters[0]) < 1 {
		return ""
	}
	parameters := queryParameters[0]
	parameters = strings.Replace(parameters, `"`, "", -1)
	parametersTable := strings.Split(parameters, ",")
	for _, parameter := range parametersTable {
		*query += " " + opt + `:"` + strings.Trim(parameter, " ") + `"`
	}
	return parameters
}

func fetchIssues(client *github.Client, query string) (*github.IssuesSearchResult, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	res, _, err := client.Search.Issues(context.Background(), query, opts)
	return res, err
}
