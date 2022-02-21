package router

import (
	"regexp"
	"strings"
)

// node holds a radix tree node including content (URI part), node children, and related route.
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

// find searches for a route by traversing the radix tree.
// It returns the route and its parameters.
func (t *tree) find(method, uri string) (*Route, map[string]string) {
	parts := strings.Split(method+uri, "/")
	parameters := map[int]string{}

	if node := t.search(t.head, parts, 0, parameters); node != nil {
		return node.Route, t.extract(node.Route, parameters)
	}

	return nil, map[string]string{}
}

// extract finds the route parameters by mapping the parameter position and route path.
func (t *tree) extract(route *Route, parameters map[int]string) map[string]string {
	parts := strings.Split(route.Method+route.Path, "/")

	list := map[string]string{}
	for i, value := range parameters {
		list[parts[i][1:]] = value
	}

	return list
}

// search finds the route by recursive traversing the radix tree.
func (t *tree) search(parent *node, parts []string, position int, parameters map[int]string) *node {
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
				return t.search(child, parts, position+1, parameters)
			}
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

// match compares a URI part with a route path part and returns the boolean result.
func (t *tree) match(token, part string) (bool, string, string) {
	if strings.HasPrefix(token, ":") {
		name := token[1:]
		if pattern, exist := t.patterns[name]; exist {
			pathPattern := regexp.MustCompile("^" + pattern + "$")
			if pathPattern.MatchString(part) {
				return true, name, part
			}
		} else {
			return true, name, part
		}
	}
	return token == part, "", ""
}

// newTree creates a new tree instance.
func newTree() *tree {
	return &tree{head: newNode(nil, "*"), patterns: map[string]string{}}
}
