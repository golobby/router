package router

import (
	"regexp"
	"strings"
)

// node holds a radix tree node including content (route path part), node children, and related route.
type node struct {
	content  string
	Route    *Route
	Children []*node
}

// newNode creates a new node instance.
func newNode(route *Route, token string) *node {
	return &node{content: token, Route: route}
}

// tree holds radix tree head node and route parameter patterns.
type tree struct {
	patterns map[string]string
	head     *node
}

// add appends a new route and inserts required nodes in the radix tree.
func (t *tree) add(route *Route) {
	parts := strings.Split(route.Method+route.Path, "/")
	t.insert(t.head, newNode(route, parts[len(parts)-1]), parts, 0)
}

// findByRequest searches for a route by request method and URI.
// It returns the route and its parameters.
func (t *tree) findByRequest(method, uri string) (*Route, map[string]string) {
	parts := strings.Split(method+uri, "/")
	parameters := map[int]string{}

	if node := t.searchByParts(t.head, parts, 0, parameters); node != nil {
		return node.Route, t.extractParameters(node.Route, parameters)
	}

	return nil, map[string]string{}
}

// findByName searches for a route by name.
func (t *tree) findByName(name string) *Route {
	if node := t.searchByName(t.head, name); node != nil {
		return node.Route
	}

	return nil
}

// extractParameters finds the route parameters (name-value pairs) by mapping the parameter position and route path.
func (t *tree) extractParameters(route *Route, parameterValues map[int]string) map[string]string {
	routeParts := strings.Split(route.Method+route.Path, "/")

	parameters := map[string]string{}
	for i, value := range parameterValues {
		parameters[routeParts[i][1:]] = value
	}

	return parameters
}

// searchByParts finds the node by parts.
func (t *tree) searchByParts(parent *node, parts []string, position int, parameters map[int]string) *node {
	isLeaf := position == len(parts)-1

	for _, child := range parent.Children {
		delete(parameters, position)
		if ok, key, value := t.match(child.content, parts[position]); ok {
			if key != "" {
				parameters[position] = value
			}

			if isLeaf {
				return child
			} else {
				return t.searchByParts(child, parts, position+1, parameters)
			}
		}
	}

	return nil
}

// searchByName finds the node by route name.
func (t *tree) searchByName(node *node, name string) *node {
	if node == nil {
		return nil
	}

	if node.Route != nil {
		if node.Route.Name == name {
			return node
		}
	}

	for _, child := range node.Children {
		if n := t.searchByName(child, name); n != nil {
			return n
		}
	}

	return nil
}

// insert adds a new route to the radix tree by recursive traversing.
func (t *tree) insert(parent, node *node, parts []string, index int) {
	isLeaf := index == len(parts)-1

	for i, child := range parent.Children {
		if child.content == parts[index] {
			if isLeaf {
				node.Children = parent.Children[i].Children
				parent.Children[i] = node
			} else {
				t.insert(child, node, parts, index+1)
			}
			return
		}
	}

	if isLeaf {
		parent.Children = append(parent.Children, node)
	} else {
		newNode := newNode(nil, parts[index])
		parent.Children = append(parent.Children, newNode)
		t.insert(newNode, node, parts, index+1)
	}
}

// match compares a request URI part with a route path part and returns the boolean result.
func (t *tree) match(routePart, UriPart string) (bool, string, string) {
	if strings.HasPrefix(routePart, ":") {
		name := routePart[1:]
		if pattern, exist := t.patterns[name]; exist {
			pathPattern := regexp.MustCompile("^" + pattern + "$")
			if pathPattern.MatchString(UriPart) {
				return true, name, UriPart
			}
		} else {
			return true, name, UriPart
		}
	}
	return routePart == UriPart, "", ""
}

// newTree creates a new tree instance.
func newTree() *tree {
	return &tree{head: newNode(nil, "*"), patterns: map[string]string{}}
}
