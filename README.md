# ghfetch

Your GitHub profile in a ~nut~shell.

The graphic on the left is an ASCII version of the [contribution graph](https://docs.github.com/en/account-and-profile/how-tos/contribution-settings/viewing-contributions-on-your-profile).

![demo](./assets/demo.png)

## Installation

```sh
go install github.com/ethn1ee/ghfetch
```

This program uses the GitHub GraphQL API, which requires a GitHub access token. See the next session for how to create one.

## GitHub Token

To create a token from `Settings > Developer settings > Personal Access Tokens`. The minimal permission required for this program is `read:profile`.

After creating the token, set it as an environment variable `GHFETCH_TOKEN`.

## Usage

### Contribution Graph

By default, the program shows the past 6 months' contribution. You can change this by providing the number of years as a flag. The value can be a float.

```sh
ghfetch <username> --years 2
```

Note that setting a wide range will result in a horizontally long graph.
