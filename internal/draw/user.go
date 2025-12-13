package draw

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/ethn1ee/ghfetch/internal/github"
	"github.com/fatih/color"
)

func FormatUser(user *github.User) string {
	buf := new(bytes.Buffer)

	uc := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	buf.WriteString(uc(user.Username + "\n"))

	bc := color.New(color.FgHiBlue).SprintFunc()
	buf.WriteString(bc(strings.Repeat("-", len(user.Username)) + "\n"))

	buf.WriteString(userKV("Name", user.Name))
	buf.WriteString(userKV("Bio", user.Bio))
	buf.WriteString(userKV("Followers", strconv.Itoa(user.Followers)))
	buf.WriteString(userKV("Following", strconv.Itoa(user.Following)))
	buf.WriteString(userKV("Joined at", user.JoinedAt.Format(time.RFC822)))

	return buf.String()
}

func userKV(key, value string) string {
	kc := color.New(color.FgCyan, color.Bold).SprintFunc()
	vc := color.New(color.FgWhite).SprintFunc()

	return kc(key) + ": " + vc(value) + "\n"
}
