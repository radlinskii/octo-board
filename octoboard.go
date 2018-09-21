package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
)

var templates = template.Must(template.ParseFiles(filepath.Join("templates", "index.html"), filepath.Join("templates", "search.html")))

// Content is a type that will be dispatched to home page template.
type Content struct {
	Label, Organization, Language string
	Issues                        []GithubIssue
}

// GithubIssue is a type that holds needed values of Github issue.
type GithubIssue struct {
	Title          string
	Repo           string
	HTMLURL        string
	Number         int
	Body           string
	CommentsNumber int
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := "is:open"

	label := buildQuery(&query, r, "label")
	language := buildQuery(&query, r, "language")
	org := buildQuery(&query, r, "org")

	client := github.NewClient(nil)
	issuesPayload, err := fetchIssues(client, query)

	var issues []GithubIssue
	for _, issue := range issuesPayload.Issues {
		issues = append(issues, GithubIssue{
			Title:          issue.GetTitle(),
			Repo:           getRepositoryFullName(issue.GetRepositoryURL()),
			HTMLURL:        issue.GetURL(),
			Number:         issue.GetNumber(),
			Body:           issue.GetBody(),
			CommentsNumber: issue.GetComments(),
		})
	}

	err = templates.ExecuteTemplate(w, "search.html", Content{Issues: issues, Label: label, Organization: org, Language: language})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRepositoryFullName(url string) string {
	return strings.Replace(url, "https://api.github.com/repos/", "", 1)
}

func buildQuery(query *string, r *http.Request, opt string) string {
	queryParameters, ok := r.URL.Query()[opt]
	if !ok || len(queryParameters[0]) < 1 {
		log.Println(r.URL)
		log.Printf("Url Parameter %s is missing", opt)
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
