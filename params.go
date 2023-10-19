package op

import "regexp"

var pathParamRegex = regexp.MustCompile(`{([^}]+)}`)

// parsePathParams gives the list of path parameters in a path.
// Example : /item/{user}/{id} -> [user, id]
func parsePathParams(path string) []string {
	matches := pathParamRegex.FindAllStringSubmatch(path, -1)

	params := make([]string, len(matches))
	for i, match := range matches {
		params[i] = match[1]
	}
	return params
}
