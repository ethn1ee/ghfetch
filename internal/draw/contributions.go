package draw

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
)

const (
	GraphHorizontal = iota
	GraphVertical
)

var (
	charSet  = []rune(" ·+=#")
	colorSet = []func(...any) string{
		fmt.Sprint,
		color.New(color.BgBlack).Sprint,
		color.New(color.BgGreen).Sprint,
		color.New(color.BgBlue).Sprint,
		color.New(color.BgWhite).Sprint,
	}
)

func FormatContributions(contributions [][]int, dir int) string {
	switch dir {
	case GraphHorizontal:
		return contributionsHorizontal(contributions)
	case GraphVertical:
		return contributionsVertical(contributions)
	default:
		return ""
	}
}

func contributionsHorizontal(contributions [][]int) string {
	buf := new(bytes.Buffer)
	for _, r := range contributions {
		for _, c := range r {
			buf.WriteString(colorSet[c](" "))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func contributionsVertical(contributions [][]int) string {
	buf := new(bytes.Buffer)
	for i := range len(contributions[0]) {
		for j := range 7 {
			c := contributions[j][i]
			buf.WriteString(colorSet[c](string(charSet[c])))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
