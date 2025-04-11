package githubapi

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/google/go-github/v57/github"
	"github.com/t-daisuke/gh-utui/internal"
	"golang.org/x/oauth2"
)

type RealGitHubClient struct {
	Client *github.Client
}

func NewRealGitHubClient() *RealGitHubClient {
	// 環境変数GH_TOKENからトークンを取得する
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGitHubToken()},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return &RealGitHubClient{
		Client: client,
	}
}

func (c *RealGitHubClient) FetchPullRequests(limit int, author string, state string, reviewer string) ([]internal.PullRequest, error) {
	// 現時点ではGitHub API v4 (GraphQL)では検索機能が制限されているため、
	// 一時的にgh CLIを使用します。将来的にこれをGraphQLに変更することを検討してください。
	ghCmd := exec.Command("gh",
		"search", "prs",
		"--limit", fmt.Sprintf("%d", limit),
		"--json", "number,title,updatedAt,repository",
	)

	if state != "" {
		ghCmd.Args = append(ghCmd.Args, "--state", state)
	}

	if reviewer != "" {
		ghCmd.Args = append(ghCmd.Args, "--review-requested", reviewer)
	}

	if author != "" {
		ghCmd.Args = append(ghCmd.Args, "--author", author)
	}

	output, err := ghCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to call gh CLI: %w", err)
	}

	var prs []internal.PullRequest
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return prs, nil
}
