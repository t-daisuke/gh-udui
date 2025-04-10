package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
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

// 各種コメント取得のためのレスポンスチャネル型
type commentResponse struct {
	issueComments  []internal.IssueComment
	reviews        []internal.Review
	reviewComments []internal.ReviewComment
	err            error
	responseType   string
}

// FetchAllCommentsParallel は3種類のコメントを並列に取得する
func (c *RealGitHubClient) FetchAllCommentsParallel(owner, repo string, number int) ([]internal.IssueComment, []internal.Review, []internal.ReviewComment, error) {
	respCh := make(chan commentResponse, 3)

	// Issue Commentsを非同期で取得
	go func() {
		comments, err := c.FetchIssueComments(owner, repo, number)
		respCh <- commentResponse{issueComments: comments, err: err, responseType: "issueComments"}
	}()

	// Reviewsを非同期で取得
	go func() {
		reviews, err := c.FetchPullRequestReviews(owner, repo, number)
		respCh <- commentResponse{reviews: reviews, err: err, responseType: "reviews"}
	}()

	// Review Commentsを非同期で取得
	go func() {
		comments, err := c.FetchPullRequestReviewComments(owner, repo, number)
		respCh <- commentResponse{reviewComments: comments, err: err, responseType: "reviewComments"}
	}()

	// 結果を受け取る
	var issueComments []internal.IssueComment
	var reviews []internal.Review
	var reviewComments []internal.ReviewComment
	var firstError error

	// 3つの結果を待つ
	for i := 0; i < 3; i++ {
		resp := <-respCh
		if resp.err != nil && firstError == nil {
			firstError = fmt.Errorf("error fetching %s: %w", resp.responseType, resp.err)
		}

		switch resp.responseType {
		case "issueComments":
			issueComments = resp.issueComments
		case "reviews":
			reviews = resp.reviews
		case "reviewComments":
			reviewComments = resp.reviewComments
		}
	}

	return issueComments, reviews, reviewComments, firstError
}

// GitHub APIを使用してIssue Commentsを取得
func (c *RealGitHubClient) FetchIssueComments(owner, repo string, number int) ([]internal.IssueComment, error) {
	ctx := context.Background()

	// APIからコメントを取得
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []internal.IssueComment

	for {
		comments, resp, err := c.Client.Issues.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			// フォールバック: gh CLIを使用
			return c.fetchIssueCommentsWithCLI(owner, repo, number)
		}

		// レスポンスを内部形式に変換
		for _, comment := range comments {
			allComments = append(allComments, internal.IssueComment{
				Body:      comment.GetBody(),
				CreatedAt: comment.GetCreatedAt().Format(time.RFC3339),
				User: struct {
					Login string `json:"login"`
				}{
					Login: comment.User.GetLogin(),
				},
			})
		}

		// ページネーション処理
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}

// GitHub CLIを使用してIssue Commentsを取得（フォールバック用）
func (c *RealGitHubClient) fetchIssueCommentsWithCLI(owner, repo string, number int) ([]internal.IssueComment, error) {
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

// GitHub APIを使用してPR Reviewsを取得
func (c *RealGitHubClient) FetchPullRequestReviews(owner, repo string, number int) ([]internal.Review, error) {
	ctx := context.Background()

	// APIからレビューを取得
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allReviews []internal.Review

	for {
		reviews, resp, err := c.Client.PullRequests.ListReviews(ctx, owner, repo, number, opts)
		if err != nil {
			// フォールバック: gh CLIを使用
			return c.fetchPullRequestReviewsWithCLI(owner, repo, number)
		}

		// レスポンスを内部形式に変換
		for _, review := range reviews {
			allReviews = append(allReviews, internal.Review{
				Body:        review.GetBody(),
				SubmittedAt: review.GetSubmittedAt().Format(time.RFC3339),
				State:       review.GetState(),
				User: struct {
					Login string `json:"login"`
				}{
					Login: review.User.GetLogin(),
				},
			})
		}

		// ページネーション処理
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReviews, nil
}

// GitHub CLIを使用してPR Reviewsを取得（フォールバック用）
func (c *RealGitHubClient) fetchPullRequestReviewsWithCLI(owner, repo string, number int) ([]internal.Review, error) {
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

// GitHub APIを使用してPR Review Commentsを取得
func (c *RealGitHubClient) FetchPullRequestReviewComments(owner, repo string, number int) ([]internal.ReviewComment, error) {
	ctx := context.Background()

	// APIからレビューコメントを取得
	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []internal.ReviewComment

	for {
		comments, resp, err := c.Client.PullRequests.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			// フォールバック: gh CLIを使用
			return c.fetchPullRequestReviewCommentsWithCLI(owner, repo, number)
		}

		// レスポンスを内部形式に変換
		for _, comment := range comments {
			allComments = append(allComments, internal.ReviewComment{
				Body:      comment.GetBody(),
				CreatedAt: comment.GetCreatedAt().Format(time.RFC3339),
				User: struct {
					Login string `json:"login"`
				}{
					Login: comment.User.GetLogin(),
				},
			})
		}

		// ページネーション処理
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}

// GitHub CLIを使用してPR Review Commentsを取得（フォールバック用）
func (c *RealGitHubClient) fetchPullRequestReviewCommentsWithCLI(owner, repo string, number int) ([]internal.ReviewComment, error) {
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
