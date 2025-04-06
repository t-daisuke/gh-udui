package cmd

import (
	"fmt"
	"log"

	// "github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/t-daisuke/gh-udui/internal/aggregator"
	"github.com/t-daisuke/gh-udui/internal/githubapi"
)

var (
	limit int

	// 本番用の GitHubAPI 実装
	gitHubClient githubapi.GitHubAPI = githubapi.NewRealGitHubClient()
)

var rootCmd = &cobra.Command{
	Use:   "gh-udui",
	Short: "A GitHub CLI extension for viewing PR comments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fetching %d PRs...\n", limit)

		// 1. PRを取得
		prs, err := gitHubClient.FetchPullRequests(limit, "@me")
		if err != nil {
			log.Printf("Error fetching PRs: %v\n", err)
			return
		}

		// 2. PR一覧を表示
		// repoColor := color.New(color.FgHiBlue).SprintFunc()
		// prNumColor := color.New(color.FgHiGreen).SprintFunc()
		// titleColor := color.New(color.FgHiWhite).SprintFunc()
		// updatedColor := color.New(color.FgHiYellow).SprintFunc()

		/* for _, pr := range prs {
			// PRの基本情報を表示
			fmt.Printf("%s#%s %s (updated: %s)\n",
				repoColor(pr.Repository.Name),
				prNumColor(pr.Number),
				titleColor(pr.Title),
				updatedColor(pr.UpdatedAt))

			// 2-1. owner/repo を分割
			owner, repo, err := githubapi.SplitOwnerRepo(pr.Repository.Name)
			if err != nil {
				// リポジトリ名が想定外の形式だった場合
				log.Printf("Skipping invalid repo name: %s\n", pr.Repository.Name)
				continue
			}

			// 2-2. コメントを取得
			comments, err := githubapi.NewRealGitHubClient().FetchIssueComments(owner, repo, pr.Number)
			if err != nil {
				log.Printf("Error fetching comments for %s: %v\n", pr.Repository.Name, err)
				continue
			}

			// 2-3. コメント表示 (後でBot除外や時系列ソートを入れられる)
			for _, c := range comments {
				// ここでシンプルにBodyだけ出す
				// Bot判定: c.User.Login に "bot" が含まれてたらスキップするなどは後で実装
				fmt.Printf("  - %s (#%d) created at: %s by %s\n",
					c.Body, c.ID, c.CreatedAt, c.User.Login)
			}

			// PRの区切り
			fmt.Println("---")
		} */
		for _, pr := range prs {
			fmt.Printf("\n%s#%d %s\n", pr.Repository.Name, pr.Number, pr.Title)

			// (1) owner, repo 抽出
			owner, repo, _ := githubapi.SplitOwnerRepo(pr.Repository.Name)

			// (2) 3種類のコメントを取得
			issueCs, _ := gitHubClient.FetchIssueComments(owner, repo, pr.Number)
			reviews, _ := gitHubClient.FetchPullRequestReviews(owner, repo, pr.Number)
			reviewCs, _ := gitHubClient.FetchPullRequestReviewComments(owner, repo, pr.Number)

			// (3) aggregatorパッケージでUnifiedCommentに変換
			uIssueCs := aggregator.ConvertIssueComments(issueCs)
			uReviews := aggregator.ConvertPullRequestReviews(reviews)
			uReviewCs := aggregator.ConvertPullRequestReviewComments(reviewCs)

			// (4) 全部まとめて1つのスライスに
			allComments := append(uIssueCs, uReviews...)
			allComments = append(allComments, uReviewCs...)

			// (5) Bot除外 → 時系列ソート → 最新5件
			allComments = aggregator.FilterOutBots(allComments)
			aggregator.SortByCreatedAtDesc(allComments)
			top5 := aggregator.TopN(allComments, 5)

			// (6) 表示
			for _, c := range top5 {
				fmt.Printf("- %s (%s):\n%s\n---\n", c.User, c.CreatedAt, c.Body)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 5, "Number of PRs to fetch (default: 5)")
}

// 実行用
func Execute() error {
	return rootCmd.Execute()
}

// テストでモックを注入するための関数 (任意)
func SetGitHubClient(client githubapi.GitHubAPI) {
	gitHubClient = client
}

// テストでrootCmdを取得するための関数 (任意)
func GetRootCmd() *cobra.Command {
	return rootCmd
}
