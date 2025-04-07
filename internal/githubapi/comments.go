package githubapi

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/t-daisuke/gh-udui/internal"
)

// "owner/repo"形式の文字列からownerとrepoを分割する
func SplitOwnerRepo(nameWithOwner string) (owner, repo string, err error) {
	parts := strings.Split(nameWithOwner, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", nameWithOwner)
	}
	return parts[0], parts[1], nil
}

func (c *RealGitHubClient) FetchIssueComments(owner, repo string, number int) ([]internal.IssueComment, error) {
	cmd := exec.Command("gh", "api", "--paginate",
		fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, number))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	var comments []internal.IssueComment
	if err := json.Unmarshal(out, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *RealGitHubClient) FetchPullRequestReviews(owner, repo string, number int) ([]internal.Review, error) {
	cmd := exec.Command("gh", "api", "--paginate",
		fmt.Sprintf("repos/%s/%s/pulls/%d/reviews", owner, repo, number))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull request reviews: %w", err)
	}

	var reviews []internal.Review
	if err := json.Unmarshal(out, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (c *RealGitHubClient) FetchPullRequestReviewComments(owner, repo string, number int) ([]internal.ReviewComment, error) {
	cmd := exec.Command("gh", "api", "--paginate",
		fmt.Sprintf("repos/%s/%s/pulls/%d/comments", owner, repo, number))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review comments: %w", err)
	}

	var comments []internal.ReviewComment
	if err := json.Unmarshal(out, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}
