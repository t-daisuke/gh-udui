package githubapi

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/t-daisuke/gh-udui/internal"
)

type RealGitHubClient struct{}

func NewRealGitHubClient() *RealGitHubClient {
	return &RealGitHubClient{}
}

func (c *RealGitHubClient) FetchPullRequests(limit int, author string, state string, reviewer string) ([]internal.PullRequest, error) {
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
