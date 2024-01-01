package builtins

import (
	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/utils"
)

var GlobalFunctions = exec.NewContext(map[string]interface{}{
	"cycler":    Cycler,
	"dict":      Dict,
	"joiner":    Joiner,
	"lipsum":    Lipsum,
	"namespace": Namespace,
	"range":     Range,
})

func Range(va *exec.VarArgs) <-chan int {
	var (
		start = 0
		stop  = -1
		step  = 1
	)
	switch len(va.Args) {
	case 1:
		stop = va.Args[0].Integer()
	case 2:
		start = va.Args[0].Integer()
		stop = va.Args[1].Integer()
	case 3:
		start = va.Args[0].Integer()
		stop = va.Args[1].Integer()
		step = va.Args[2].Integer()
		// default:
		// 	return nil, errors.New("range expect signature range([start, ]stop[, step])")
	}
	chnl := make(chan int)
	go func() {
		for i := start; i < stop; i += step {
			chnl <- i
		}

		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func Dict(va *exec.VarArgs) *exec.Value {
	dict := exec.NewDict()
	for key, value := range va.KwArgs {
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

func Cycler(va *exec.VarArgs) *exec.Value {
	c := &cycler{}
	for _, arg := range va.Args {
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

func Joiner(va *exec.VarArgs) *exec.Value {
	p := va.ExpectKwArgs([]*exec.KwArg{{"sep", ","}})
	if p.IsError() {
		return exec.AsValue(errors.Wrapf(p, `wrong signature for 'joiner'`))
	}
	sep := p.KwArgs["sep"].String()
	j := &joiner{sep: sep}
	return exec.AsValue(j.String)
}

// type namespace map[string]interface{}

func Namespace(va *exec.VarArgs) map[string]interface{} {
	ns := map[string]interface{}{}
	for key, value := range va.KwArgs {
		ns[key] = value
	}
	return ns
}

func Lipsum(va *exec.VarArgs) *exec.Value {
	p := va.ExpectKwArgs([]*exec.KwArg{
		{"n", 5},
		{"html", true},
		{"min", 20},
		{"max", 100},
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
