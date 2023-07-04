package statements

import (
	"fmt"
	"math"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
)

type ForStmt struct {
	key             string
	value           string // only for maps: for key, value in map
	objectEvaluator nodes.Expression
	ifCondition     nodes.Expression

	bodyWrapper  *nodes.Wrapper
	emptyWrapper *nodes.Wrapper
}

func (stmt *ForStmt) Position() *tokens.Token { return stmt.bodyWrapper.Position() }
func (stmt *ForStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("ForStmt(Line=%d Col=%d)", t.Line, t.Col)
}

type LoopInfos struct {
	index      int
	index0     int
	revindex   int
	revindex0  int
	first      bool
	last       bool
	length     int
	depth      int
	depth0     int
	PrevItem   *exec.Value
	NextItem   *exec.Value
	_lastValue *exec.Value
}

func (li *LoopInfos) Cycle(va *exec.VarArgs) *exec.Value {
	return va.Args[int(math.Mod(float64(li.index0), float64(len(va.Args))))]
}

func (li *LoopInfos) Changed(value *exec.Value) bool {
	same := li._lastValue != nil && value.EqualValueTo(li._lastValue)
	li._lastValue = value
	return !same
}

func (node *ForStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) (forError error) {
	obj := r.Eval(node.objectEvaluator)
	if obj.IsError() {
		return obj
	}

	// Create loop struct
	items := exec.NewDict()

	// First iteration: filter values to ensure proper LoopInfos
	obj.Iterate(func(idx, count int, key, value *exec.Value) bool {
		sub := r.Inherit()
		ctx := sub.Ctx
		pair := &exec.Pair{}

		// There's something to iterate over (correct type and at least 1 item)
		// Update loop infos and public context
		if node.value != "" && !key.IsString() && key.Len() == 2 {
			key.Iterate(func(idx, count int, key, value *exec.Value) bool {
				switch idx {
				case 0:
					ctx.Set(node.key, key)
					pair.Key = key
				case 1:
					ctx.Set(node.value, key)
					pair.Value = key
				}
				return true
			}, func() {})
		} else {
			ctx.Set(node.key, key)
			pair.Key = key
			if value != nil {
				ctx.Set(node.value, value)
				pair.Value = value
			}
		}

		if node.ifCondition != nil {
			if !sub.Eval(node.ifCondition).IsTrue() {
				return true
			}
		}
		items.Pairs = append(items.Pairs, pair)
		return true
	}, func() {})

	// 2nd pass: all values are defined, render
	length := len(items.Pairs)
	loop := &LoopInfos{
		first:  true,
		index0: -1,
	}
	if len(items.Pairs) == 0 && node.emptyWrapper != nil {
		if err := r.Inherit().ExecuteWrapper(node.emptyWrapper); err != nil {
			return err
		}
	}
	for idx, pair := range items.Pairs {
		sub := r.Inherit()
		ctx := sub.Ctx

		ctx.Set(node.key, pair.Key)
		if pair.Value != nil {
			ctx.Set(node.value, pair.Value)
		}

		ctx.Set("loop", loop)
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

		// Render elements with updated context
		err := sub.ExecuteWrapper(node.bodyWrapper)
		if err != nil {
			return err
		}
	}

	return forError
}

func forParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &ForStmt{}

	// Arguments parsing
	var valueToken *tokens.Token
	keyToken := args.Match(tokens.Name)
	if keyToken == nil {
		return nil, args.Error("Expected an key identifier as first argument for 'for'-tag", nil)
	}

	if args.Match(tokens.Comma) != nil {
		// Value name is provided
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
	stmt.objectEvaluator = objectEvaluator
	stmt.key = keyToken.Val
	if valueToken != nil {
		stmt.value = valueToken.Val
	}

	if args.MatchName("if") != nil {
		ifCondition, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		stmt.ifCondition = ifCondition
	}

	if !args.End() {
		return nil, args.Error("Malformed for-loop args.", nil)
	}

	// Body wrapping
	wrapper, endargs, err := p.WrapUntil("else", "endfor")
	if err != nil {
		return nil, err
	}
	stmt.bodyWrapper = wrapper

	if !endargs.End() {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.EndTag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = p.WrapUntil("endfor")
		if err != nil {
			return nil, err
		}
		stmt.emptyWrapper = wrapper

		if !endargs.End() {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return stmt, nil
}

func init() {
	All.Register("for", forParser)
}
