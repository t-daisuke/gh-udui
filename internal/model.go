package internal

type Repository struct {
	Name string `json:"nameWithOwner"`
}

type PullRequest struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	UpdatedAt  string     `json:"updatedAt"`
	Repository Repository `json:"repository"`
}
