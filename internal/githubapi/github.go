package githubapi

import (
	"os/exec"
	"strings"

	"github.com/t-daisuke/gh-utui/internal"
)

// GitHubAPI インターフェース: PR検索やコメント取得などをまとめる
type GitHubAPI interface {
	FetchPullRequests(limit int, author string, state string, reviewer string) ([]internal.PullRequest, error)
	FetchIssueComments(owner, repo string, number int) ([]internal.IssueComment, error)
	FetchPullRequestReviews(owner, repo string, number int) ([]internal.Review, error)
	FetchPullRequestReviewComments(owner, repo string, number int) ([]internal.ReviewComment, error)
	FetchAllCommentsParallel(owner, repo string, number int) ([]internal.IssueComment, []internal.Review, []internal.ReviewComment, error)
}

// getGitHubToken はGH_TOKENから認証トークンを取得する
func getGitHubToken() string {
	// gh authコマンドの出力を使用して認証情報を取得する
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
