package router

// state holds the collector common attributes.
type state struct {
	prefix     string
	middleware []Middleware
}

// collector holds the active and outer collector states.
// The active state holds attributes for the active collector.
// The outer state holds attributes for the outer (previous) collector.
type collector struct {
	active *state
	outer  *state
}

// update sets a new state as the active state to affect the forthcoming routes.
// It stores the previous active state as the outer state to roll back when the collector is closed.
func (g *collector) update(prefix string, middleware []Middleware) {
	g.outer.prefix = g.active.prefix
	g.outer.middleware = g.active.middleware
	g.active.prefix += prefix
	g.active.middleware = append(g.active.middleware, middleware...)
}

// rollback closes the active collector and rollbacks the active state to the outer (previous) state.
func (g *collector) rollback() {
	g.active.prefix = g.outer.prefix
	g.active.middleware = g.outer.middleware
}

//
func newCollector() *collector {
	return &collector{
		active: &state{
			prefix: "",
		},
		outer: &state{
			prefix: "",
		},
	}
}
