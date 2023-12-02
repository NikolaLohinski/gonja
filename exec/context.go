package exec

type Context struct {
	data   map[string]interface{}
	parent *Context
}

func NewContext(data map[string]interface{}) *Context {
	return &Context{data: data}
}

func EmptyContext() *Context {
	return &Context{data: map[string]interface{}{}}
}

func (ctx *Context) Has(name string) bool {
	_, exists := ctx.data[name]
	if !exists && ctx.parent != nil {
		return ctx.parent.Has(name)
	}
	return exists
}

func (ctx *Context) Get(name string) (interface{}, bool) {
	value, exists := ctx.data[name]
	if exists {
		return value, true
	} else if ctx.parent != nil {
		return ctx.parent.Get(name)
	} else {
		return nil, false
	}
}

func (ctx *Context) Set(name string, value interface{}) {
	ctx.data[name] = value
}

func (ctx *Context) Inherit() *Context {
	return &Context{
		data:   map[string]interface{}{},
		parent: ctx,
	}
}

// Update updates this context with the key/value pairs from a map.
func (ctx *Context) Update(other *Context) *Context {
	if other == nil {
		return ctx
	}
	for k, v := range other.data {
		ctx.data[k] = v
	}
	return ctx
}
