package exec

import (
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/pkg/errors"
)

type Environment struct {
	Filters           FilterSet
	ControlStructures ControlStructureSet
	Tests             TestSet
	Context           *Context
	Methods           Methods
}

type FilterSet map[string]FilterFunction

// Exists returns true if the given filter is already registered
func (fs FilterSet) Exists(name string) bool {
	_, existing := fs[name]
	return existing
}

// Register registers a new filter. If there's already a filter with the same
// name, Register will panic. You usually want to call this
// function in the filter's init() function:
// http://golang.org/doc/effective_go.html#init
func (fs *FilterSet) Register(name string, fn FilterFunction) error {
	if fs.Exists(name) {
		return errors.Errorf("filter with name '%s' is already registered", name)
	}
	(*fs)[name] = fn
	return nil
}

// Replace replaces an already registered filter with a new implementation. Use this
// function with caution since it allows you to change existing filter behaviour.
func (fs *FilterSet) Replace(name string, fn FilterFunction) error {
	if !fs.Exists(name) {
		return errors.Errorf("filter with name '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*fs)[name] = fn
	return nil
}

func (fs *FilterSet) Update(other FilterSet) FilterSet {
	for name, filter := range other {
		(*fs)[name] = filter
	}
	return *fs
}

type ControlStructureSet map[string]parser.ControlStructureParser

// Exists returns true if the given test is already registered
func (ss ControlStructureSet) Exists(name string) bool {
	_, existing := ss[name]
	return existing
}

// Registers a new tag. You usually want to call this
// function in the tag's init() function:
// http://golang.org/doc/effective_go.html#init
func (ss *ControlStructureSet) Register(name string, parser parser.ControlStructureParser) error {
	if ss.Exists(name) {
		return errors.Errorf("ControlStructure '%s' is already registered", name)
	}
	(*ss)[name] = parser
	// &controlStructure{
	// 	name:   name,
	// 	parser: parserFn,
	// }
	return nil
}

// Replaces an already registered tag with a new implementation. Use this
// function with caution since it allows you to change existing tag behaviour.
func (ss *ControlStructureSet) Replace(name string, parser parser.ControlStructureParser) error {
	if !ss.Exists(name) {
		return errors.Errorf("ControlStructure '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ss)[name] = parser
	// controlStructures[name] = &controlStructure{
	// 	name:   name,
	// 	parser: parserFn,
	// }
	return nil
}

func (ss *ControlStructureSet) Update(other ControlStructureSet) ControlStructureSet {
	for name, parser := range other {
		(*ss)[name] = parser
	}
	return *ss
}

// TestSet maps test names to their TestFunction handler
type TestSet map[string]TestFunction

// Exists returns true if the given test is already registered
func (ts TestSet) Exists(name string) bool {
	_, existing := ts[name]
	return existing
}

// Register registers a new test. If there's already a test with the same
// name, RegisterTest will panic. You usually want to call this
// function in the test's init() function:
// http://golang.org/doc/effective_go.html#init
func (ts *TestSet) Register(name string, fn TestFunction) error {
	if ts.Exists(name) {
		return errors.Errorf("test with name '%s' is already registered", name)
	}
	(*ts)[name] = fn
	return nil
}

// Replace replaces an already registered test with a new implementation. Use this
// function with caution since it allows you to change existing test behaviour.
func (ts *TestSet) Replace(name string, fn TestFunction) error {
	if !ts.Exists(name) {
		return errors.Errorf("test with name '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ts)[name] = fn
	return nil
}

func (ts *TestSet) Update(other TestSet) TestSet {
	for name, test := range other {
		(*ts)[name] = test
	}
	return *ts
}

type Method[I interface{}] func(self I, selfValue *Value, arguments *VarArgs) (interface{}, error)

type Methods struct {
	Bool  MethodSet[bool]
	Int   MethodSet[int]
	Float MethodSet[float64]
	Str   MethodSet[string]
	Dict  MethodSet[map[string]interface{}]
	List  MethodSet[[]interface{}]
}

type MethodSet[I interface{}] map[string]Method[I]

func (m MethodSet[I]) Exists(name string) bool {
	_, existing := m[name]
	return existing
}
