package draw

import (
	"bytes"
	"fmt"
)

func Contributions(contributions [][]int, charSet []rune) {
	for _, r := range contributions {
		buf := new(bytes.Buffer)
		for _, c := range r {
			buf.WriteRune(charSet[c])
		}
		fmt.Println(buf.String())
	}
}
