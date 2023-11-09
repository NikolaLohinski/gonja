package gonja

import (
	"io/ioutil"
	"sync"

	"github.com/goph/emperror"

	"github.com/nikolalohinski/gonja/builtins"
	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/loaders"
)

type Environment struct {
	*exec.EvalConfig
	Loader loaders.Loader

	Cache      map[string]*exec.Template
	CacheMutex sync.Mutex
}

func NewEnvironment(cfg *config.Config, loader loaders.Loader) *Environment {
	env := &Environment{
		EvalConfig: exec.NewEvalConfig(cfg),
		Loader:     loader,
		Cache:      map[string]*exec.Template{},
	}
	env.EvalConfig.Loader = env
	env.Filters.Update(builtins.Filters)
	env.Statements.Update(builtins.Statements)
	env.Tests.Update(builtins.Tests)
	env.Globals.Merge(builtins.Globals)
	env.Globals.Set("gonja", map[string]interface{}{
		"version": VERSION,
	})
	return env
}

// CleanCache cleans the template cache. If filenames is not empty,
// it will remove the template caches of those filenames.
// Or it will empty the whole template cache. It is thread-safe.
func (env *Environment) CleanCache(filenames ...string) {
	env.CacheMutex.Lock()
	defer env.CacheMutex.Unlock()

	if len(filenames) == 0 {
		env.Cache = map[string]*exec.Template{}
	}

	for _, filename := range filenames {
		delete(env.Cache, filename)
	}
}

// FromCache is a convenient method to cache templates. It is thread-safe
// and will only compile the template associated with a filename once.
// If Environment.Debug is true (for example during development phase),
// FromCache() will not cache the template and instead recompile it on any
// call (to make changes to a template live instantaneously).
func (env *Environment) FromCache(filename string) (*exec.Template, error) {
	if env.Config.Debug {
		// Recompile on any request
		return env.FromFile(filename)
	}

	env.CacheMutex.Lock()
	defer env.CacheMutex.Unlock()

	tpl, has := env.Cache[filename]

	// Cache miss
	if !has {
		tpl, err := env.FromFile(filename)
		if err != nil {
			return nil, err
		}
		env.Cache[filename] = tpl
		return tpl, nil
	}

	// Cache hit
	return tpl, nil
}

// FromString loads a template from string and returns a Template instance.
func (env *Environment) FromString(tpl string) (*exec.Template, error) {
	return exec.NewTemplate("string", tpl, env.EvalConfig)
}

// FromBytes loads a template from bytes and returns a Template instance.
func (env *Environment) FromBytes(tpl []byte) (*exec.Template, error) {
	return exec.NewTemplate("bytes", string(tpl), env.EvalConfig)
}

// FromFile loads a template from a filename and returns a Template instance.
func (env *Environment) FromFile(filename string) (*exec.Template, error) {
	fd, err := env.Loader.Read(filename)
	if err != nil {
		return nil, emperror.With(err, "filename", filename)
	}
	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, emperror.With(err, "filename", filename)
	}

	return exec.NewTemplate(filename, string(buf), env.EvalConfig)
}

func (env *Environment) GetTemplate(filename string) (*exec.Template, error) {
	return env.FromFile(filename)
}

func (env *Environment) Path(path string) (string, error) {
	return env.Loader.Resolve(path)
}
