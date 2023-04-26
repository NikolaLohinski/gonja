package exec

import (
	"strings"

	"github.com/nikolalohinski/gonja/nodes"
	"github.com/pkg/errors"
	// "github.com/nikolalohinski/gonja/nodes"
)

// FilterFunction is the type filter functions must fulfil
type Macro func(params *VarArgs) *Value

type MacroSet map[string]Macro

// Exists returns true if the given filter is already registered
func (ms MacroSet) Exists(name string) bool {
	_, existing := ms[name]
	return existing
}

// Register registers a new filter. If there's already a filter with the same
// name, Register will panic. You usually want to call this
// function in the filter's init() function:
// http://golang.org/doc/effective_go.html#init
//
// See http://www.john-doe.de/post/gonja/ for more about
// writing filters and tags.
func (ms *MacroSet) Register(name string, fn Macro) error {
	if ms.Exists(name) {
		return errors.Errorf("filter with name '%s' is already registered", name)
	}
	(*ms)[name] = fn
	return nil
}

// Replace replaces an already registered filter with a new implementation. Use this
// function with caution since it allows you to change existing filter behaviour.
func (ms *MacroSet) Replace(name string, fn Macro) error {
	if !ms.Exists(name) {
		return errors.Errorf("filter with name '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ms)[name] = fn
	return nil
}

func MacroNodeToFunc(node *nodes.Macro, r *Renderer) (Macro, error) {
	// Compute default values once
	defaultKwargs := []*KwArg{}
	var err error
	for _, pair := range node.Kwargs {
		key := r.Eval(pair.Key).String()
		value := r.Eval(pair.Value)
		if value.IsError() {
			return nil, errors.Wrapf(value, `Unable to evaluate parameter %s=%s`, key, pair.Value)
		}
		defaultKwargs = append(defaultKwargs, &KwArg{key, value.Interface()})
	}

	return func(params *VarArgs) *Value {
		var out strings.Builder
		sub := r.Inherit()
		sub.Out = &out

		if err != nil {
			return AsValue(err)
		}
		p := params.Expect(len(node.Args), defaultKwargs)
		if p.IsError() {
			return AsValue(errors.Wrapf(p, `Wrong '%s' macro signature`, node.Name))
		}
		for idx, arg := range p.Args {
			sub.Ctx.Set(node.Args[idx], arg)
		}
		for key, value := range p.KwArgs {
			sub.Ctx.Set(key, value)
		}
		err := sub.ExecuteWrapper(node.Wrapper)
		if err != nil {
			return AsValue(errors.Wrapf(err, `Unable to execute macro '%s`, node.Name))
		}
		return AsSafeValue(out.String())
	}, err
}
