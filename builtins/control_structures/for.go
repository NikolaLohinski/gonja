package controlStructures

import (
	"bytes"
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type ForControlStructure struct {
	Key             string
	Value           string // only for maps: for key, value in map
	ObjectEvaluator nodes.Expression
	IfCondition     nodes.Expression
	Recursive       bool

	BodyWrapper  *nodes.Wrapper
	EmptyWrapper *nodes.Wrapper
}

func (fcs *ForControlStructure) Position() *tokens.Token {
	return fcs.BodyWrapper.Position()
}
func (fcs *ForControlStructure) String() string {
	t := fcs.Position()
	return fmt.Sprintf("ForControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

// LoopInfos exposes Jinja2's `loop` variable inside a {% for %} body.
type LoopInfos struct {
	index     int
	index0    int
	length    int
	revindex  int
	revindex0 int
	first     bool
	last      bool
	depth     int
	depth0    int
	PrevItem  *exec.Value
	NextItem  *exec.Value
	lastValue *exec.Value

	// Recursive support: when the for-loop is declared `recursive`,
	// recurseFn is set to a function that renders the loop body for a
	// new iterable, allowing templates to call `{{ loop(children) }}`.
	recurseFn func(items *exec.Value) *exec.Value
}

// GetAttribute implements exec.AttributeGetter so templates can access
// loop.* fields/methods with Jinja2-style lowercase names.
func (li *LoopInfos) GetAttribute(name string) (*exec.Value, bool) {
	switch name {
	case "index":
		return exec.AsValue(li.index), true
	case "index0":
		return exec.AsValue(li.index0), true
	case "revindex":
		return exec.AsValue(li.revindex), true
	case "revindex0":
		return exec.AsValue(li.revindex0), true
	case "first":
		return exec.AsValue(li.first), true
	case "last":
		return exec.AsValue(li.last), true
	case "length":
		return exec.AsValue(li.length), true
	case "depth":
		return exec.AsValue(li.depth), true
	case "depth0":
		return exec.AsValue(li.depth0), true
	case "previtem":
		return li.PrevItem, true
	case "nextitem":
		return li.NextItem, true
	case "cycle":
		fn := func(va *exec.VarArgs) *exec.Value {
			if len(va.Args) == 0 {
				return exec.AsValue("")
			}
			idx := li.index0 % len(va.Args)
			if idx < 0 {
				idx += len(va.Args)
			}
			return va.Args[idx]
		}
		return exec.AsValue(fn), true
	case "changed":
		fn := func(va *exec.VarArgs) *exec.Value {
			var current *exec.Value
			if len(va.Args) == 1 {
				current = va.Args[0]
			} else {
				items := make([]any, 0, len(va.Args))
				for _, a := range va.Args {
					items = append(items, a.Interface())
				}
				current = exec.AsValue(items)
			}
			same := li.lastValue != nil && current.EqualValueTo(li.lastValue)
			li.lastValue = current
			return exec.AsValue(!same)
		}
		return exec.AsValue(fn), true
	}
	return exec.AsValue(nil), false
}

// loopCallable wraps LoopInfos so that `{{ loop(children) }}` works
// inside a recursive {% for %} loop. It implements both attribute access
// (delegated to the inner LoopInfos) and callability (which triggers
// another render pass over the supplied iterable).
type loopCallable struct {
	infos *LoopInfos
}

func (lc *loopCallable) GetAttribute(name string) (*exec.Value, bool) {
	return lc.infos.GetAttribute(name)
}

// Call lets the template invoke `loop(items)` to recurse.
func (lc *loopCallable) Call(va *exec.VarArgs) *exec.Value {
	if lc.infos.recurseFn == nil {
		return exec.AsValue("")
	}
	if len(va.Args) == 0 {
		return exec.AsValue("")
	}
	return lc.infos.recurseFn(va.Args[0])
}

func (fcs *ForControlStructure) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) (forError error) {
	obj := r.Eval(fcs.ObjectEvaluator)
	if obj.IsError() {
		return obj
	}
	return fcs.renderIterable(r, obj)
}

// renderIterable contains the main render logic, factored out so the
// recursive callable can invoke it again with a new iterable while
// writing to the same outer renderer.
func (fcs *ForControlStructure) renderIterable(r *exec.Renderer, obj *exec.Value) error {
	// First pass: materialise items, applying any {% if cond %} filter.
	items := exec.NewDict()
	obj.Iterate(func(idx, count int, key, value *exec.Value) bool {
		sub := r.Inherit()
		ctx := sub.Environment.Context
		pair := &exec.Pair{}

		if fcs.Value != "" && !key.IsString() && key.Len() == 2 {
			key.Iterate(func(idx, count int, key, value *exec.Value) bool {
				switch idx {
				case 0:
					ctx.Set(fcs.Key, key)
					pair.Key = key
				case 1:
					ctx.Set(fcs.Value, key)
					pair.Value = key
				}
				return true
			}, func() {})
		} else {
			ctx.Set(fcs.Key, key)
			pair.Key = key
			if value != nil {
				ctx.Set(fcs.Value, value)
				pair.Value = value
			}
		}

		if fcs.IfCondition != nil {
			if !sub.Eval(fcs.IfCondition).IsTrue() {
				return true
			}
		}
		items.Pairs = append(items.Pairs, pair)
		return true
	}, func() {})

	// Empty case.
	length := len(items.Pairs)
	if length == 0 {
		if fcs.EmptyWrapper != nil {
			return r.Inherit().ExecuteWrapper(fcs.EmptyWrapper)
		}
		return nil
	}

	// Nested-loop depth.
	depth := 1
	if parentLoop, ok := r.Environment.Context.Get("loop"); ok {
		switch p := parentLoop.(type) {
		case *LoopInfos:
			depth = p.depth + 1
		case *loopCallable:
			depth = p.infos.depth + 1
		}
	}

	loop := &LoopInfos{
		first:  true,
		index0: -1,
		length: length,
		depth:  depth,
		depth0: depth - 1,
	}

	// Wire recursion if requested.
	var loopContextValue any = loop
	if fcs.Recursive {
		callable := &loopCallable{infos: loop}
		callable.infos.recurseFn = func(childItems *exec.Value) *exec.Value {
			// Render the body for the child iterable into a buffer and
			// return that as a safe string for inline interpolation.
			buf := new(bytes.Buffer)
			sub := r.Inherit()
			originalOutput := sub.Output
			sub.Output = buf
			if err := fcs.renderIterable(sub, childItems); err != nil {
				sub.Output = originalOutput
				return exec.AsValue(err)
			}
			sub.Output = originalOutput
			return exec.AsSafeValue(buf.String())
		}
		loopContextValue = callable
	}

	for idx, pair := range items.Pairs {
		sub := r.Inherit()
		ctx := sub.Environment.Context

		ctx.Set(fcs.Key, pair.Key)
		if pair.Value != nil {
			ctx.Set(fcs.Value, pair.Value)
		}

		ctx.Set("loop", loopContextValue)
		loop.index0 = idx
		loop.index = loop.index0 + 1
		if idx == 1 {
			loop.first = false
		}
		if idx+1 == length {
			loop.last = true
		}
		loop.revindex = length - idx
		loop.revindex0 = length - (idx + 1)

		if idx == 0 {
			loop.PrevItem = exec.AsValue(nil)
		} else {
			pp := items.Pairs[idx-1]
			if pp.Value != nil {
				loop.PrevItem = exec.AsValue([2]*exec.Value{pp.Key, pp.Value})
			} else {
				loop.PrevItem = pp.Key
			}
		}

		if idx == length-1 {
			loop.NextItem = exec.AsValue(nil)
		} else {
			np := items.Pairs[idx+1]
			if np.Value != nil {
				loop.NextItem = exec.AsValue([2]*exec.Value{np.Key, np.Value})
			} else {
				loop.NextItem = np.Key
			}
		}

		err := sub.ExecuteWrapper(fcs.BodyWrapper)
		if err != nil {
			if _, ok := err.(*LoopBreakError); ok {
				break
			}
			if _, ok := err.(*LoopContinueError); ok {
				continue
			}
			return err
		}
	}

	return nil
}

func forParser(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	cs := &ForControlStructure{}

	// Arguments parsing
	var valueToken *tokens.Token
	keyToken := args.Match(tokens.Name)
	if keyToken == nil {
		return nil, args.Error("Expected an key identifier as first argument for 'for'-tag", nil)
	}

	if args.Match(tokens.Comma) != nil {
		valueToken = args.Match(tokens.Name)
		if valueToken == nil {
			return nil, args.Error("Value name must be an identifier.", nil)
		}
	}

	if args.Match(tokens.In) == nil {
		return nil, args.Error("Expected keyword 'in'.", nil)
	}

	objectEvaluator, err := args.ParseExpression()
	if err != nil {
		return nil, err
	}
	cs.ObjectEvaluator = objectEvaluator
	cs.Key = keyToken.Val
	if valueToken != nil {
		cs.Value = valueToken.Val
	}

	if args.MatchName("if") != nil {
		ifCondition, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		cs.IfCondition = ifCondition
	}

	// Jinja2 `recursive` keyword.
	if args.MatchName("recursive") != nil {
		cs.Recursive = true
	}

	if !args.End() {
		return nil, args.Error("Malformed for-loop args.", nil)
	}

	// Body wrapping
	wrapper, endargs, err := p.WrapUntil("else", "endfor")
	if err != nil {
		return nil, err
	}
	cs.BodyWrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.EndTag == "else" {
		wrapper, endargs, err = p.WrapUntil("endfor")
		if err != nil {
			return nil, err
		}
		cs.EmptyWrapper = wrapper

		if !endargs.End() {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return cs, nil
}
