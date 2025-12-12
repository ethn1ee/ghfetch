package github

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client   *graphql.Client
	username string
}

const (
	endpoint = "https://api.github.com/graphql"
)

func NewClient(ctx context.Context, username string, token string) *GitHub {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	client := graphql.NewClient(endpoint, httpClient)

	return &GitHub{
		client:   client,
		username: username,
	}
}
