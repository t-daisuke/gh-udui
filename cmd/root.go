package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/t-daisuke/gh-udui/internal/aggregator"
	"github.com/t-daisuke/gh-udui/internal/githubapi"
)

var (
	limit int
	// デフォルトは本番用のRealGitHubClientを使う
	gitHubClient githubapi.GitHubAPI = githubapi.NewRealGitHubClient()
)

var rootCmd = &cobra.Command{
	Use:   "gh-udui",
	Short: "A GitHub CLI extension for viewing PR comments",
	Long: `gh-udui is a GitHub CLI extension that helps you view PR comments, reviews, and discussions
in a user-friendly format. It fetches PR information and displays the latest comments
from non-bot users.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fetching %d PRs...\n", limit)

		// 1. PR一覧を取得
		prs, err := gitHubClient.FetchPullRequests(limit, "@me")
		if err != nil {
			log.Printf("Error fetching PRs: %v\n", err)
			return
		}

		// 2. 色付き出力のためのSprintFuncを準備
		repoColor := color.New(color.FgHiBlue).SprintFunc()
		prNumColor := color.New(color.FgHiGreen).SprintFunc()
		titleColor := color.New(color.FgHiWhite).SprintFunc()
		updatedColor := color.New(color.FgHiYellow).SprintFunc()

		// 3. PRごとにコメントを取得し、表示
		for _, pr := range prs {
			// 3.1 PRの基本情報を色付きで表示
			fmt.Printf(
				"%s#%s %s (updated: %s)\n",
				repoColor(pr.Repository.Name),
				prNumColor(pr.Number),
				titleColor(pr.Title),
				updatedColor(pr.UpdatedAt),
			)

			// 3.2 owner/repo を分割
			owner, repo, err := githubapi.SplitOwnerRepo(pr.Repository.Name)
			if err != nil {
				log.Printf("Invalid repo name: %s\n", pr.Repository.Name)
				fmt.Println("---")
				continue
			}

			// 3.3 IssueComments・Reviews・ReviewComments を取得
			issueCs, err := gitHubClient.FetchIssueComments(owner, repo, pr.Number)
			if err != nil {
				log.Printf("Error fetching IssueComments: %v\n", err)
				fmt.Println("---")
				continue
			}
			reviews, err := gitHubClient.FetchPullRequestReviews(owner, repo, pr.Number)
			if err != nil {
				log.Printf("Error fetching Reviews: %v\n", err)
				fmt.Println("---")
				continue
			}
			reviewCs, err := gitHubClient.FetchPullRequestReviewComments(owner, repo, pr.Number)
			if err != nil {
				log.Printf("Error fetching ReviewComments: %v\n", err)
				fmt.Println("---")
				continue
			}

			// 3.4 3種類のコメントを UnifiedComment に変換
			uIssueCs := aggregator.ConvertIssueComments(issueCs)
			uReviews := aggregator.ConvertPullRequestReviews(reviews)
			uReviewCs := aggregator.ConvertPullRequestReviewComments(reviewCs)

			// 3.5 マージ
			allComments := append(uIssueCs, uReviews...)
			allComments = append(allComments, uReviewCs...)

			// 3.6 Bot除外 → 時系列ソート(最新→古い順) → 上位5件
			allComments = aggregator.FilterOutBots(allComments)
			aggregator.SortByCreatedAtDesc(allComments)
			top5 := aggregator.TopN(allComments, 5) // TODO ここもオプションで変えられるようにしたい。

			// 3.7 コメントを表示
			//   (コメントも色を付ける場合は適宜 SprintFunc を用意)
			userColor := color.New(color.FgCyan).SprintFunc()
			dateColor := color.New(color.FgHiYellow).SprintFunc()

			for _, c := range top5 {
				// CreatedAt をフォーマットして一行にまとめる例
				fmt.Printf("%s (%s):\n%s\n",
					userColor(c.User),
					dateColor(c.CreatedAt.Format("2006-01-02 15:04:05")),
					c.Body)
			}

			// 3.8 PRの区切り線
			fmt.Println("---")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 5, "Number of PRs to fetch (default: 5)")
}

func Execute() error {
	return rootCmd.Execute()
}

// テストでモックを注入したい場合に呼び出すための関数
func SetGitHubClient(client githubapi.GitHubAPI) {
	gitHubClient = client
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
