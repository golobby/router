package router

import (
	"fmt"
	"regexp"
	"strings"
)

func Match(path, uri string) bool {
	mc := regexp.MustCompile(`{([^}]+)}`)
	parameters := mc.FindAllString(path, -1)

	for _, p := range parameters {
		pattern := "(?P<" + p[1:len(p)-1] + ">[^/]+?)"
		path = strings.Replace(path, p, pattern, 1)
	}

	r := regexp.MustCompile("^" + path + "$")

	args := map[string]string{}

	if r.MatchString(uri) {
		for i, v := range r.FindAllStringSubmatch(uri, -1)[0] {
			if i > 0 {
				args[parameters[i-1]] = v
			}
		}

		fmt.Println(args)
		return true
	}

	return false
}
