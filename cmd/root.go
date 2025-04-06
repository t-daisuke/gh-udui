package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/t-daisuke/gh-udui/internal/githubapi"
)

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

		prs, err := githubapi.FetchPullRequests(limit, "@me")
		if err != nil {
			log.Printf("Error fetching PRs: %v\n", err)
			return
		}

		repoColor := color.New(color.FgHiBlue).SprintFunc()
		prNumColor := color.New(color.FgHiGreen).SprintFunc()
		titleColor := color.New(color.FgHiWhite).SprintFunc()
		updatedColor := color.New(color.FgHiYellow).SprintFunc()

		for _, pr := range prs {
			fmt.Printf("%s#%s %s (updated: %s)\n",
				repoColor(pr.Repository.Name),
				prNumColor(pr.Number),
				titleColor(pr.Title),
				updatedColor(pr.UpdatedAt))
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 5, "Number of PRs to fetch (default: 5)")
}

func Execute() error {
	return rootCmd.Execute()
}
