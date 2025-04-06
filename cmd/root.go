package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// Repository は repository フィールドの情報を格納
type Repository struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// PullRequest は gh search で返却される JSON に対応
type PullRequest struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	UpdatedAt  string     `json:"updatedAt"`
	Repository Repository `json:"repository"`
}

var (
	limit int
)

var rootCmd = &cobra.Command{
	Use:   "gh-udui",
	Short: "A GitHub CLI extension for viewing PR comments",
	Long: `gh-udui is a GitHub CLI extension that helps you view PR comments, reviews, and discussions
in a user-friendly format. It fetches PR information and displays the latest comments
from non-bot users.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fetching %d PRs...\n", limit)

		// 1. gh CLI コマンドを組み立て
		ghCmd := exec.Command("gh",
			"search", "prs",
			"--author", "@me", //将来optionにする
			"--limit", fmt.Sprintf("%d", limit), //将来optionにする
			"--json", "number,title,updatedAt,repository",
		)

		// 2. 標準出力を取得
		output, err := ghCmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("failed to call gh CLI: %v\n%s\n", err, exitErr.Stderr)
			} else {
				fmt.Printf("failed to call gh CLI: %v\n", err)
			}
			return
		}

		// 3. JSON をパース
		var prs []PullRequest
		if err := json.Unmarshal(output, &prs); err != nil {
			fmt.Printf("failed to unmarshal JSON: %v\n", err)
			return
		}

		// 4. 取得した PR を表示
		for i, pr := range prs {
			fmt.Printf("%d: [%d] %s (repo: %s/%s, updated: %s)\n",
				i+1, pr.Number, pr.Title, pr.Repository.Owner.Login, pr.Repository.Name, pr.UpdatedAt)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 5, "Number of PRs to fetch (default: 5)")
}

func Execute() error {
	return rootCmd.Execute()
}
