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

var templates = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

type Content struct {
	Issues []GithubIssue
}

type GithubIssue struct {
	Title          string
	URL            string
	Body           string
	CommentsNumber int
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	query := "is:open org:google language:go"

	buildQuery(&query, r, "label")

	client := github.NewClient(nil)
	issuesPayload, err := fetchIssues(client, query)

	var issues []GithubIssue
	for _, issue := range issuesPayload.Issues {
		issues = append(issues, GithubIssue{
			Title:          *issue.URL,
			URL:            *issue.URL,
			Body:           *issue.Body,
			CommentsNumber: *issue.Comments,
		})
	}

	err = templates.ExecuteTemplate(w, "index.html", Content{issues})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildQuery(query *string, r *http.Request, opt string) {
	labelsQuery, ok := r.URL.Query()[opt]
	if !ok || len(labelsQuery[0]) < 1 {
		log.Println(r.URL)
		log.Printf("Url Parameter %s is missing", opt)
		return
	}
	labels := labelsQuery[0]
	labelsTable := strings.Split(labels, ",")
	for _, label := range labelsTable {
		*query += " " + opt + `:"` + strings.Trim(label, " ") + `"`
	}
}

func fetchIssues(client *github.Client, query string) (*github.IssuesSearchResult, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	res, _, err := client.Search.Issues(context.Background(), query, opts)
	return res, err
}
