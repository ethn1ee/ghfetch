package draw

import (
	"os"

	"github.com/ethn1ee/ghfetch/internal/github"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func Print(user *github.User, contributions [][]int) {
	c := FormatContributions(contributions, GraphHorizontal)
	u := FormatUser(user)

	PrintTable([][]string{{c, u}})
}

func PrintTable(data [][]string) {
	symbols := tw.NewSymbolCustom("minimal").
		WithRow("").
		WithColumn("")

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(tw.Rendition{
				Symbols: symbols,
			}),
		),
	)

	_ = table.Bulk(data)
	_ = table.Render()
}
