<div align="center">
<img src="./docs/logo.svg" width="200"/>
<h1><code>gonja</code></h1>
</div>

`gonja` is a pure `go` implementation of the [Jinja template engine](https://jinja.palletsprojects.com/). It aims to be _mostly_ compatible with the original `python` implementation but also provides additional features to compensate the lack of `python` scripting capabilities.

## Usage

### As a library

Install/update using `go get`:
```
go get github.com/nikolalohinski/gonja/v2
```

### As a `terraform` provider

This `gonja` library has been packaged as a `terraform` provider. For more information, please refer to the [dedicated documentation](https://registry.terraform.io/providers/NikolaLohinski/jinja/latest/docs).

## Example

```golang
package main

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func main() {
	loader, err := loaders.NewMemoryLoader(map[string]string{"/_": "Hello {{ name | capitalize }}!"})
	if err != nil {
		panic(err)
	}

	template, err := exec.NewTemplate("/_", gonja.DefaultConfig, loader, gonja.DefaultEnvironment)
	if err != nil {
		panic(err)
	}
	
	out, err := template.Execute(exec.NewContext(map[string]interface{}{"name": "bob"}))
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // Prints: Hello Bob!
}
```

## Documentation

* For details on how the **Jinja** template language works, please refer to [the Jinja documentation](https://jinja.palletsprojects.com) ;
* **gonja** API documentation is available on [godoc](https://godoc.org/github.com/nikolalohinski/gonja/v2) ;
* **filters**: please refer to [`docs/filters.md`](docs/filters.md) ;
* **control structures**: please take a look at [`docs/control_structures.md`](docs/control_structures.md) ;
* **tests**: please see [`docs/tests.md`](docs/tests.md) ;
* **global functions**: please browse through [`docs/global_functions.md`](docs/global_functions.md).


## Migrating from `v1` to `v2`

As this project now aims to reproduce the behavior of the `python` Jinja engine as closely as possible, some backwards incompatible changes have been made from the initial draft and need to be taken into account when upgrading from `v1.X.X`. Therefore, no `v1.X.X` versions will be maintained.

The following steps can be used as general guidelines to migrate from `v1` to `v2`:

* All references to `gonja` need to be changed from `"github.com/nikolalohinski/gonja"` to `"github.com/nikolalohinski/gonja/v2"`
* The following top level global variables/functions have been removed/updated and need to be adjusted accordingly:
	* `DefaultEnv` function is now called `DefaultEnvironment` and its properties have changed. See [gonja.go](./gonja.go) and [exec/environment.go](./exec/environment.go) for details.
	* `FromCache` function has been removed as caching logic was removed. If required, it can be done by implementing a custom `Loader` (see [`loaders/loader.go`](./loaders/loader.go)).
	* `Globals` is now referred to as `DefaultContext`
* What was called a `Statement` is now referred to as `ControlStructure` to be closer to `python`'s Jinja glossary and may require changes
* What was called `Globals` is now called `GlobalFunctions` to be closer to `python`'s Jinja glossary and may require changes


## Limitations 

* **format**: `format` does **not** take Python's string format syntax as a parameter, instead it takes Go's. Essentially `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`.
* **escape** / **force_escape**: Unlike Jinja's behavior, the `escape`-filter is applied immediately. Therefore there is no need for a `force_escape`-filter yet.

## Tribute

A massive thank you to the original author [@noirbizarre](https://github.com/noirbizarre) for doing the initial work in https://github.com/noirbizarre/gonja which this project was forked from.
