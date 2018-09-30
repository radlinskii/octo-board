package controller

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/radlinskii/octo-board/viewmodel"

	"github.com/google/go-github/v18/github"
)

type search struct {
	searchTemplate *template.Template
}

func (s search) registerRoutes() {
	http.HandleFunc("/search", s.handleSearch)
}

func (s search) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := "type:issue state:open"

	if checkFilter(r, "uncommented") {
		query += " comments:0"
	}
	if checkFilter(r, "unassigned") {
		query += " no:assignee"
	}

	label := buildQuery(&query, r, "label")
	language := buildQuery(&query, r, "language")
	org := buildQuery(&query, r, "org")

	page := getPage(r)

	client := github.NewClient(nil)
	issuesPayload, err := fetchIssues(client, query, page)

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
	err = s.searchTemplate.Execute(w, viewmodel.Content{
		Issues:       issues,
		Label:        label,
		Organization: org,
		Language:     language,
		NextPage:     page + 1,
		PrevPage:     page - 1,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkFilter(r *http.Request, filter string) bool {
	fs := r.URL.Query().Get(filter)
	if len(fs) < 1 {
		return false
	}
	f, err := strconv.Atoi(fs)
	if err != nil {
		return false
	}
	if f != 1 {
		return false
	}
	return true
}

func getRepositoryFullName(url string) string {
	return strings.Replace(url, "https://api.github.com/repos/", "", 1)
}

func getPage(r *http.Request) int {
	ps := r.URL.Query().Get("page")
	if len(ps) < 1 {
		return 1
	}
	page, err := strconv.Atoi(ps)
	if err != nil {
		return 1
	}
	return page
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

func fetchIssues(client *github.Client, query string, page int) (*github.IssuesSearchResult, error) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc", ListOptions: github.ListOptions{PerPage: 20, Page: page}}
	res, _, err := client.Search.Issues(context.Background(), query, opts)
	return res, err
}
