package router

// state holds a group attributes.
type state struct {
	prefix      string
	middlewares []Middleware
}

// newState creates a new state instance.
func newState(prefix string, middlewares []Middleware) *state {
	return &state{prefix: prefix, middlewares: middlewares}
}

// stateStack holds the stack of states (group attributes).
type stateStack struct {
	states []*state
}

// prefix returns current state (group) prefix.
func (g *stateStack) prefix() string {
	if len(g.states) > 0 {
		return g.states[len(g.states)-1].prefix
	}
	return ""
}

// middlewares returns current state (group) middlewares.
func (g *stateStack) middlewares() []Middleware {
	if len(g.states) > 0 {
		return g.states[len(g.states)-1].middlewares
	}
	return []Middleware{}
}

// push adds a new state (group) to the stack.
func (g *stateStack) push(prefix string, middleware []Middleware) {
	g.states = append(g.states, newState(g.prefix()+prefix, append(g.middlewares(), middleware...)))
}

// pop removes (closes) the last state (group).
func (g *stateStack) pop() {
	if len(g.states) > 0 {
		g.states = g.states[:len(g.states)-1]
	}
}

// newStateStack creates a new stateStack instance.
func newStateStack() *stateStack {
	return &stateStack{}
}
