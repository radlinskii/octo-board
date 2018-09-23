package controller

import (
	"context"
	"github.com/radlinskii/octo-board/viewmodel"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

type search struct {
	searchTemplate *template.Template
}

func (s search) registerRoutes() {
	http.HandleFunc("/search", s.handleSearch)
}

func (s search) handleSearch(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "text/html")
	err = s.searchTemplate.Execute(w, viewmodel.Content{Issues: issues, Label: label, Organization: org, Language: language})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRepositoryFullName(url string) string {
	return strings.Replace(url, "https://api.github.com/repos/", "", 1)
}

func buildQuery(query *string, r *http.Request, opt string) string {
	parameters := r.URL.Query().Get(opt)
	if len(parameters) < 1 {
		return ""
	}
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
