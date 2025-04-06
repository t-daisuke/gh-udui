package githubapi

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/t-daisuke/gh-udui/internal"
)

func FetchPullRequests(limit int, author string) ([]internal.PullRequest, error) {
	// 1. gh CLI コマンドを組み立て
	ghCmd := exec.Command("gh",
		"search", "prs",
		"--author", author,
		"--limit", fmt.Sprintf("%d", limit),
		"--json", "number,title,updatedAt,repository",
	)

	// 2. 標準出力を取得
	output, err := ghCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to call gh CLI: %w", err)
	}

	// 3. JSON をパース
	var prs []internal.PullRequest
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return prs, nil
}
