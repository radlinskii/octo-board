package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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
	/*
		labels, ok := r.URL.Query()["labels"]
		if !ok || len(labels[0]) < 1 {
			log.Println(r.URL)
			log.Println("Url Parameter 'labels' is missing")
			return
		}
		label := labels[0]
	*/
	client := github.NewClient(nil)
	issuesPayload, err := FetchIssues(client)

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

func FetchIssues(client *github.Client) (*github.IssuesSearchResult, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "asc"}
	res, _, err := client.Search.Issues(context.Background(), "is:open org:google language:go label:\"good first issue\"", opts)
	return res, err
}
