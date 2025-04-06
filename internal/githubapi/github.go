package githubapi

import "github.com/t-daisuke/gh-udui/internal"

// GitHubAPI インターフェース: PR検索やコメント取得などをまとめる
type GitHubAPI interface {
	FetchPullRequests(limit int, author string) ([]internal.PullRequest, error)
	// 今後コメントやレビュー取得関数も追加する:
	// FetchIssueComments(owner, repo string, number int) ([]internal.IssueComment, error)
}
