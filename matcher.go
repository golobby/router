package router

import (
	"regexp"
	"strings"
)

// matcher holds route parameter patterns and uses them to find appropriate routes for incoming requests.
type matcher struct {
	parameters map[string]string
}

// addParameter adds a new route parameter pattern to the list.
func (m *matcher) addParameter(name, pattern string) {
	m.parameters[name] = pattern
}

// match compares the route path with the request URI.
// It will return the boolean result and route arguments if the comparison is successful.
func (m matcher) match(path, uri string) (bool, map[string]string) {
	parameters := regexp.MustCompile(`{[^}]+}`).FindAllString(path, -1)
	for _, parameter := range parameters {
		name := parameter[1 : len(parameter)-1]

		pattern := "(?P<" + name + ">[^/]+?)"
		if definedPattern, exist := m.parameters[name]; exist {
			pattern = "(?P<" + name + ">" + definedPattern + ")"
		}

		path = strings.Replace(path, parameter, pattern, -1)
	}

	arguments := map[string]string{}
	pathPattern := regexp.MustCompile("^" + path + "$")
	if pathPattern.MatchString(uri) {
		for i, value := range pathPattern.FindAllStringSubmatch(uri, -1)[0] {
			if i > 0 {
				arguments[parameters[i-1]] = value
			}
		}

		return true, arguments
	}

	return false, nil
}

// newMatcher creates a new matcher.
func newMatcher() *matcher {
	return &matcher{map[string]string{}}
}
