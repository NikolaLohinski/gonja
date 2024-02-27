package exec

import "sync"

type Context struct {
	data   map[string]interface{}
	parent *Context
	lock   sync.Mutex
}

func NewContext(data map[string]interface{}) *Context {
	return &Context{data: data}
}

func EmptyContext() *Context {
	return &Context{data: map[string]interface{}{}}
}

func (ctx *Context) Has(name string) bool {
	ctx.lock.Lock()
	_, exists := ctx.data[name]
	ctx.lock.Unlock()
	if !exists && ctx.parent != nil {
		return ctx.parent.Has(name)
	}
	return exists
}

func (ctx *Context) Get(name string) (interface{}, bool) {
	ctx.lock.Lock()
	value, exists := ctx.data[name]
	ctx.lock.Unlock()
	if exists {
		return value, true
	} else if ctx.parent != nil {
		return ctx.parent.Get(name)
	} else {
		return nil, false
	}
}

func (ctx *Context) Set(name string, value interface{}) {
	ctx.lock.Lock()
	ctx.data[name] = value
	ctx.lock.Unlock()
}

func (ctx *Context) Inherit() *Context {
	ctx.lock.Lock()
	inherited := &Context{
		data:   map[string]interface{}{},
		parent: ctx,
	}
	ctx.lock.Unlock()
	return inherited
}

// Update updates this context with the key/value pairs from a map.
func (ctx *Context) Update(other *Context) *Context {
	if other == nil {
		return ctx
	}
	ctx.lock.Lock()
	for k, v := range other.data {
		ctx.data[k] = v
	}
	ctx.lock.Unlock()
	return ctx
}
