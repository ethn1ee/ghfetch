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
	years float32
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
		start := end.Add(-time.Hour * time.Duration(int(years*24*365)))

		if user.JoinedAt.After(start) {
			start = user.JoinedAt
		}

		contributions, _, err := client.GetContributions(ctx, start, end)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("failed to get contributions: %w", err))
		}

		draw.Print(user, contributions)
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
	rootCmd.Flags().Float32VarP(&years, "years", "y", 0.5, "length of contribution graph in years")
}
