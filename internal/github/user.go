package github

import (
	"context"
	"fmt"
	"time"
)

func (gh *GitHub) GetUser(ctx context.Context, username string) (*User, error) {
	var q userQuery
	if err := gh.client.Query(ctx, &q, map[string]any{
		"username": username,
	}); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	joinedAt, err := time.Parse(time.RFC3339, q.User.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start datetime: %w", err)
	}

	return &User{
		Name:      q.User.Name,
		Username:  username,
		Bio:       q.User.Bio,
		JoinedAt:  joinedAt.Local(),
		Followers: q.User.Followers.TotalCount,
		Following: q.User.Following.TotalCount,
	}, nil
}

type userQuery struct {
	User struct {
		Name      string
		Login     string
		Bio       string
		CreatedAt string
		Followers struct {
			TotalCount int
		}
		Following struct {
			TotalCount int
		}
	} `graphql:"user(login: $username)"`
}

type User struct {
	Name      string
	Username  string
	Bio       string
	JoinedAt  time.Time
	Followers int
	Following int
}
