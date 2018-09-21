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

	labelsQuery, ok := r.URL.Query()["labels"]
	if !ok || len(labelsQuery[0]) < 1 {
		log.Println(r.URL)
		log.Println("Url Parameter 'labels' is missing")
		return
	} else {
		labels := labelsQuery[0]
		labelsTable := strings.Split(labels, ",")
		for _, label := range labelsTable {
			query += ` label:"` + strings.Trim(label, " ") + `"`
		}
	}

	client := github.NewClient(nil)
	issuesPayload, err := FetchIssues(client, query)

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

func FetchIssues(client *github.Client, query string) (*github.IssuesSearchResult, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	res, _, err := client.Search.Issues(context.Background(), query, opts)
	return res, err
}
