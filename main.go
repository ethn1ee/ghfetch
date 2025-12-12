package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

const (
	endpoint  = "https://api.github.com/graphql"
	yearRange = 2
)

var charMap = map[string]rune{
	"NONE":            ' ',
	"FIRST_QUARTILE":  '·',
	"SECOND_QUARTILE": '+',
	"THIRD_QUARTILE":  '=',
	"FOURTH_QUARTILE": '#',
}

type DateTime string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("username must be provided")
	}
	username := os.Args[1]
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("environment variable GITHUB_TOKEN is required")
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, src)
	client := graphql.NewClient(endpoint, httpClient)

	var uq userQuery
	if err := client.Query(ctx, &uq, map[string]any{
		"userName": username,
	}); err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	userCreated, err := time.Parse(time.RFC3339, uq.User.CreatedAt)
	if err != nil {
		log.Fatalf("failed to parse start datetime: %v", err)
	}
	end := time.Now()
	start := end.Add(-time.Hour * 24 * 365 * yearRange)

	if userCreated.After(start) {
		start = userCreated
	}
	curr := start

	graph := make([]*bytes.Buffer, 7)
	for i := range graph {
		graph[i] = new(bytes.Buffer)
	}

	total := 0
	for curr.Before(end) {
		next := curr.Add(time.Hour * 24 * 365)
		if next.After(end) {
			next = end
		}
		rows, contributions, err := fetchContribution(ctx, client, username, curr, next)
		if err != nil {
			log.Fatalf("failed to fetch contributions: %v", err)
		}
		curr = next

		for i, r := range graph {
			_, err := rows[i].WriteTo(r)
			if err != nil {
				log.Fatalf("failed to write to graph buffer: %v", err)
			}
		}

		total += contributions
	}

	for _, r := range graph {
		fmt.Println(r.String())
	}

	fmt.Printf("Total %d contributions from %s to %s",
		total, start.Format(time.DateOnly), curr.Format(time.DateOnly),
	)
}

func fetchContribution(ctx context.Context, client *graphql.Client, username string, from time.Time, to time.Time) ([]*bytes.Buffer, int, error) {
	var q contributionQuery
	if err := client.Query(ctx, &q, map[string]any{
		"userName": username,
		"from":     DateTime(from.Format(time.RFC3339)),
		"to":       DateTime(to.Format(time.RFC3339)),
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to query: %w", err)
	}

	weeks := q.User.ContributionsCollection.ContributionCalendar.Weeks
	minDate, maxDate := time.Now(), time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().UTC().Local().Location())

	rows := make([]*bytes.Buffer, 7)
	for d := range 7 {
		buf := new(bytes.Buffer)
		for _, w := range weeks {
			if len(w.ContributionDays) <= d {
				buf.WriteRune(charMap["NONE"])
				continue
			}

			day := w.ContributionDays[d]
			buf.WriteRune(charMap[day.ContributionLevel])

			date, err := time.Parse(time.DateOnly, day.Date)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to parse date: %w", err)
			}

			if date.Before(minDate) {
				minDate = date
			} else if date.After(maxDate) {
				maxDate = date
			}
		}

		rows[d] = buf
	}

	return rows, q.User.ContributionsCollection.ContributionCalendar.TotalContributions, nil
}

type userQuery struct {
	User struct {
		CreatedAt string
	} `graphql:"user(login: $userName)"`
}

type contributionQuery struct {
	User struct {
		ContributionsCollection struct {
			ContributionCalendar struct {
				TotalContributions int
				Weeks              []struct {
					ContributionDays []struct {
						ContributionCount int
						ContributionLevel string
						Date              string
					}
				}
			}
		} `graphql:"contributionsCollection(from: $from, to: $to)"`
	} `graphql:"user(login: $userName)"`
}
