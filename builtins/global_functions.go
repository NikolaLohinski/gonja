package builtins

import (
	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/utils"
)

var GlobalFunctions = exec.NewContext(map[string]interface{}{
	"cycler":    cyclerFunction,
	"dict":      dictFunction,
	"joiner":    joinerFunction,
	"lipsum":    lipsumFunction,
	"namespace": namespaceFunction,
	"range":     rangeFunction,
})

func rangeFunction(_ *exec.Evaluator, arguments *exec.VarArgs) (<-chan int, error) {
	var (
		start = 0
		stop  = -1
		step  = 1
	)
	switch len(arguments.Args) {
	case 1:
		stop = arguments.Args[0].Integer()
	case 2:
		start = arguments.Args[0].Integer()
		stop = arguments.Args[1].Integer()
	case 3:
		start = arguments.Args[0].Integer()
		stop = arguments.Args[1].Integer()
		step = arguments.Args[2].Integer()
	default:
		return nil, errors.New("range expect signature range([start, ]stop[, step])")
	}
	channel := make(chan int)
	go func() {
		for i := start; i < stop; i += step {
			channel <- i
		}
		close(channel)
	}()
	return channel, nil
}

func dictFunction(_ *exec.Evaluator, arguments *exec.VarArgs) *exec.Value {
	dict := exec.NewDict()
	for key, value := range arguments.KwArgs {
		dict.Pairs = append(dict.Pairs, &exec.Pair{
			Key:   exec.AsValue(key),
			Value: value,
		})
	}
	return exec.AsValue(dict)
}

type cycler struct {
	values  []string
	idx     int
	getters map[string]interface{}
}

func (c *cycler) Reset() {
	c.idx = 0
	c.getters["current"] = c.values[c.idx]
}

func (c *cycler) Next() string {
	c.idx++
	value := c.getters["current"].(string)
	if c.idx >= len(c.values) {
		c.idx = 0
	}
	c.getters["current"] = c.values[c.idx]
	return value
}

func cyclerFunction(_ *exec.Evaluator, arguments *exec.VarArgs) *exec.Value {
	c := &cycler{}
	for _, arg := range arguments.Args {
		c.values = append(c.values, arg.String())
	}
	c.getters = map[string]interface{}{
		"next":  c.Next,
		"reset": c.Reset,
	}
	c.Reset()
	return exec.AsValue(c.getters)
}

type joiner struct {
	sep   string
	first bool
}

func (j *joiner) String() string {
	if !j.first {
		j.first = true
		return ""
	}
	return j.sep
}

func joinerFunction(_ *exec.Evaluator, arguments *exec.VarArgs) *exec.Value {
	p := arguments.ExpectKwArgs([]*exec.KwArg{{Name: "sep", Default: ","}})
	if p.IsError() {
		return exec.AsValue(errors.Wrapf(p, `wrong signature for 'joiner'`))
	}
	sep := p.KwArgs["sep"].String()
	j := &joiner{sep: sep}
	return exec.AsValue(j.String)
}

func namespaceFunction(_ *exec.Evaluator, arguments *exec.VarArgs) map[string]interface{} {
	ns := map[string]interface{}{}
	for key, value := range arguments.KwArgs {
		ns[key] = value
	}
	return ns
}

func lipsumFunction(_ *exec.Evaluator, arguments *exec.VarArgs) *exec.Value {
	p := arguments.ExpectKwArgs([]*exec.KwArg{
		{Name: "n", Default: 5},
		{Name: "html", Default: true},
		{Name: "min", Default: 20},
		{Name: "max", Default: 100},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrapf(p, `wrong signature for 'lipsum'`))
	}
	n := p.KwArgs["n"].Integer()
	html := p.KwArgs["html"].Bool()
	min := p.KwArgs["min"].Integer()
	max := p.KwArgs["max"].Integer()
	return exec.AsSafeValue(utils.Lipsum(n, html, min, max))
}
