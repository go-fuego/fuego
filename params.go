package fuego

import (
	"regexp"
	"strings"
)

var pathParamRegex = regexp.MustCompile(`{(.+?)}`)

// parsePathParams gives the list of path parameters in a path.
// Example : /item/{user}/{id} -> [user, id]
func parsePathParams(path string) []string {
	matches := pathParamRegex.FindAllString(path, -1)
	for i, match := range matches {
		matches[i] = strings.Trim(match, "{}")
	}
	return matches
}
