package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethn1ee/ghfetch/internal/draw"
	"github.com/ethn1ee/ghfetch/internal/github"
	"github.com/spf13/cobra"
)

var (
	years   int
	charSet string
)

var rootCmd = &cobra.Command{
	Use:   "ghfetch",
	Short: "GitHub profile in a (nut)shell",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr("username must be provided")
		}

		username := os.Args[1]
		token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			cobra.CheckErr("environment variable GITHUB_TOKEN is required")
		}

		ctx := context.Background()
		client := github.NewClient(ctx, username, token)

		user, err := client.GetUser(ctx, username)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("failed to get user: %w", err))
		}

		end := time.Now()
		start := end.Add(-time.Hour * 24 * 365 * time.Duration(years))

		if user.JoinedAt.After(start) {
			start = user.JoinedAt
		}

		graph, _, err := client.GetContributions(ctx, start, end)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("failed to get contributions: %w", err))
		}

		draw.Contributions(graph, []rune(charSet))

		fmt.Println()
		fmt.Println("Name:", user.Name)
		fmt.Println("Username:", user.Username)
		fmt.Println("Bio:", user.Bio)
		fmt.Println("Followers:", user.Followers)
		fmt.Println("Following:", user.Following)
		fmt.Println("Joined at:", user.JoinedAt)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&years, "years", "y", 1, "length of contribution graph in years")
	rootCmd.Flags().StringVar(&charSet, "charSet", " ·+=#", "contribution graph character set; must be length of 5")
}
