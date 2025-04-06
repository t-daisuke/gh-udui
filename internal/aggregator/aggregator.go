package aggregator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/t-daisuke/gh-udui/internal"
)

// ConvertIssueComments : IssueComments → []UnifiedComment
func ConvertIssueComments(issueComments []internal.IssueComment) []internal.UnifiedComment {
	var result []internal.UnifiedComment
	for _, ic := range issueComments {
		t, _ := time.Parse(time.RFC3339, ic.CreatedAt) // エラー処理省略
		result = append(result, internal.UnifiedComment{
			User:      ic.User.Login,
			Body:      ic.Body,
			CreatedAt: t,
		})
	}
	return result
}

// ConvertPullRequestReviews : Reviews → []UnifiedComment
func ConvertPullRequestReviews(reviews []internal.Review) []internal.UnifiedComment {
	var result []internal.UnifiedComment
	for _, r := range reviews {
		t, _ := time.Parse(time.RFC3339, r.SubmittedAt)
		body := r.Body
		if r.State != "" {
			body = fmt.Sprintf("[State: %s]\n%s", r.State, r.Body)
		}

		result = append(result, internal.UnifiedComment{
			User:      r.User.Login,
			Body:      body,
			CreatedAt: t,
		})
	}
	return result
}

// ConvertPullRequestReviewComments : ReviewComments → []UnifiedComment
func ConvertPullRequestReviewComments(reviewComments []internal.ReviewComment) []internal.UnifiedComment {
	var result []internal.UnifiedComment
	for _, rc := range reviewComments {
		t, _ := time.Parse(time.RFC3339, rc.CreatedAt)
		result = append(result, internal.UnifiedComment{
			User:      rc.User.Login,
			Body:      rc.Body,
			CreatedAt: t,
		})
	}
	return result
}

// FilterOutBots : Bot除外 (Userに "bot" が含まれるものは除外)
func FilterOutBots(comments []internal.UnifiedComment) []internal.UnifiedComment {
	var filtered []internal.UnifiedComment
	for _, c := range comments {
		if !strings.Contains(strings.ToLower(c.User), "bot") {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// SortByCreatedAtDesc : 日時が新しい順にソート
func SortByCreatedAtDesc(comments []internal.UnifiedComment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})
}

// TopN : 先頭N件を取り出す
func TopN(comments []internal.UnifiedComment, n int) []internal.UnifiedComment {
	if n >= len(comments) {
		return comments
	}
	return comments[:n]
}
