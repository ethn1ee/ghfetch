package github

import (
	"context"
	"fmt"
	"slices"
	"time"
)

func (gh *GitHub) GetContributions(ctx context.Context, from, to time.Time) ([][]int, int, error) {
	graph := make([][]int, 7)
	total := 0
	curr := from
	for curr.Before(to) {
		next := curr.Add(time.Hour * 24 * 365)
		if next.After(to) {
			next = to
		}

		rows, contributions, err := gh.getContributions(ctx, curr, next)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get contributions from %s to %s: %w",
				from.Format(time.DateOnly), to.Format(time.DateOnly), err,
			)
		}

		for i := range graph {
			graph[i] = slices.Concat(graph[i], rows[i])
		}

		total += contributions
		curr = next
	}

	return graph, total, nil
}

type DateTime string

var levelEnum = map[string]int{
	"NONE":            0,
	"FIRST_QUARTILE":  1,
	"SECOND_QUARTILE": 2,
	"THIRD_QUARTILE":  3,
	"FOURTH_QUARTILE": 4,
}

func (gh *GitHub) getContributions(ctx context.Context, from, to time.Time) ([][]int, int, error) {
	var q contributionQuery
	if err := gh.client.Query(ctx, &q, map[string]any{
		"userName": gh.username,
		"from":     DateTime(from.Format(time.RFC3339)),
		"to":       DateTime(to.Format(time.RFC3339)),
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to query: %w", err)
	}

	weeks := q.User.ContributionsCollection.ContributionCalendar.Weeks
	minDate, maxDate := time.Now(), time.Time{}

	rows := make([][]int, 7)
	for d := range 7 {
		for _, w := range weeks {
			if len(w.ContributionDays) <= d {
				rows[d] = append(rows[d], (levelEnum["NONE"]))
				continue
			}

			day := w.ContributionDays[d]
			rows[d] = append(rows[d], levelEnum[day.ContributionLevel])

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

	}

	return rows, q.User.ContributionsCollection.ContributionCalendar.TotalContributions, nil
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
