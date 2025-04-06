package githubapi

import "github.com/t-daisuke/gh-udui/internal"

// GitHubAPI インターフェース: PR検索やコメント取得などをまとめる
type GitHubAPI interface {
	FetchPullRequests(limit int, author string, state string) ([]internal.PullRequest, error)
	FetchIssueComments(owner, repo string, number int) ([]internal.IssueComment, error)
	FetchPullRequestReviews(owner, repo string, number int) ([]internal.Review, error)
	FetchPullRequestReviewComments(owner, repo string, number int) ([]internal.ReviewComment, error)
}
