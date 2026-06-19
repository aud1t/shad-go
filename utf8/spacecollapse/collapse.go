//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
)

func CollapseSpaces(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	wasSpace := false

	for _, r := range input {
		if !unicode.IsSpace(r) {
			sb.WriteRune(r)
			wasSpace = false
		} else if !wasSpace {
			sb.WriteRune(' ')
			wasSpace = true
		}
	}
	return sb.String()
}
