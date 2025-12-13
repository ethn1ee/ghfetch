# ghfetch

Your GitHub profile in a ~nut~shell.

![demo](./assets/demo.png)

## Installation

```sh
go install github.com/ethn1ee/ghfetch
```

This program uses the GitHub GraphQL API, which requires a GitHub access token. See the next session for how to create

## GitHub Token

To create a token from `Settings > Developer settings > Personal Access Tokens`. The minimal permission required for this program is `read:profile`.

After creating the token, set it as an environment variable `GHFETCH_TOKEN`.
