//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	linkCount := 0
	var sb strings.Builder
	sb.Grow(len(format))
	wasOpen := false
	openIndex := -1

	for i, r := range format {
		if r == '{' {
			if wasOpen {
				panic(fmt.Sprint("incorrect format: ", format))
			}
			wasOpen = true
			openIndex = i
			continue
		}

		if r == '}' {
			if !wasOpen {
				panic(fmt.Sprint("incorrect format: ", format))
			}

			var targetIdx int

			if i == openIndex+1 {
				targetIdx = linkCount
			} else {
				n, err := strconv.Atoi(format[openIndex+1 : i])
				if err != nil {
					panic(fmt.Sprint("incorrect format: ", format))
				}
				targetIdx = n
			}

			if targetIdx < 0 || targetIdx >= len(args) {
				panic(fmt.Sprintf("argument index %d out of range (len = %d)", targetIdx, len(args)))
			}

			fmt.Fprint(&sb, args[targetIdx])
			linkCount++
			wasOpen = false
			openIndex = -1
		} else if wasOpen {
			if r < '0' || r > '9' {
				panic(fmt.Sprint("incorrect format: ", format))
			}
		} else {
			sb.WriteRune(r)
		}
	}

	if wasOpen {
		panic(fmt.Sprint("incorrect format: ", format))
	}

	return sb.String()
}
