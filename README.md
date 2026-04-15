<div align="center">
<img src="docs/logo.png" />
<h1><code>gonja</code></h1>
</div>

`gonja` is a pure `go` implementation of the [Jinja template engine](https://jinja.palletsprojects.com/). It aims to be as compatible as possible with the original `python` implementation.

## Usage

### As a library

Install or update using `go get`:

```
go get github.com/ardanlabs/gonja
```

## Example

```golang
package main

import (
	"os"
	"fmt"

	"github.com/ardanlabs/gonja"
	"github.com/ardanlabs/gonja/exec"
)

func main() {
	template, err := gonja.FromString("Hello {{ name | capitalize }}!")
	if err != nil {
		panic(err)
	}

	data := exec.NewContext(map[string]interface{}{
		"name": "bob",
	})

	if err = template.Execute(os.Stdout, data); err != nil { // Prints: Hello Bob!
		panic(err)
	}
}
```

## Documentation

- For details on how the **Jinja** template language works, please refer to [the Jinja documentation](https://jinja.palletsprojects.com) ;
- **gonja** API documentation is available on [godoc](https://godoc.org/github.com/ardanlabs/gonja) ;
- **filters**: please refer to [`docs/filters.md`](docs/filters.md) ;
- **control structures**: please take a look at [`docs/control_structures.md`](docs/control_structures.md) ;
- **tests**: please see [`docs/tests.md`](docs/tests.md) ;
- **global functions**: please browse through [`docs/global_functions.md`](docs/global_functions.md).
- **global variables**: please open [`docs/global_variables.md`](docs/global_variables.md).
- **methods**: please take a peek at [`docs/methods.md`](docs/methods.md).

## Limitations

- **format**: `format` does **not** take `python`'s string format syntax as a parameter, instead it takes Go's. Essentially `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`
- **escape** / **force_escape**: Unlike Jinja's behavior, the `escape`-filter is applied immediately. Therefore there is no need for a `force_escape` filter
- Only subsets of native `python` types (`bool`, `int`, `float`, `str`, `dict` and `list`) methods have been re-implemented in Go and can slightly differ from the original ones

## Development

### Guidelines

Please read through the [contribution guidelines](./CONTRIBUTING.md) before diving into any work.

### Requirements

- Install go `>= 1.26` by following the [official documentation](https://go.dev/doc/install) ;
- Install `ginkgo` by [any means you see fit](https://onsi.github.io/ginkgo/).

### Tests

The unit tests can be run using:

```sh
ginkgo run -p ./...
```

## Tribute

A massive thank you to the original author [@noirbizarre](https://github.com/noirbizarre) for doing the initial work in https://github.com/noirbizarre/gonja which this project was forked from and [@NikolaLohinski](https://github.com/NikolaLohinski) who has been maintaining the project since.
