package internal

import "time"

type Repository struct {
	Name string `json:"nameWithOwner"`
}

type PullRequest struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	UpdatedAt  string     `json:"updatedAt"`
	Repository Repository `json:"repository"`
}

// UnifiedComment : Issue Comments, Review Comments, Reviews をまとめた共通構造
type UnifiedComment struct {
	User      string    // user.login
	Body      string    // comment/review body
	CreatedAt time.Time // created_at or submitted_at
}
