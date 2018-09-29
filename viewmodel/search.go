package viewmodel

// Content is a type that will be dispatched to home page template.
type Content struct {
	Label, Organization, Language string
	NextPage                      int
	PrevPage                      int
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
