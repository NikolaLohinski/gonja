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
)

func main() {
	tpl, err := gonja.FromString("Hello {{ name | capitalize }}!")
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(gonja.Context{"name": "bob"})
	if err != nil {
		panic(err)
	}
	fmt.Println(out) // Prints: Hello Bob!
}
```

## Documentation

* For a details on how the template language works, please refer to [the Jinja documentation](https://jinja.palletsprojects.com) ;
* `gonja` API documentation is available on [godoc](https://godoc.org/github.com/nikolalohinski/gonja/v2) ;
* filters: please refer to [`docs/filters.md`](docs/filters.md) ;
* control structures: please take a look at [`docs/control_structures.md`](docs/control_structures.md) ;
* tests: please see [`docs/tests.md`](docs/tests.md) ;
* global functions: please browse through [`docs/global_functions.md`](docs/global_functions.md).

## Limitations 

* **format**: `format` does **not** take Python's string format syntax as a parameter, instead it takes Go's. Essentially `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`.
* **escape** / **force_escape**: Unlike Jinja's behavior, the `escape`-filter is applied immediately. Therefore there is no need for a `force_escape`-filter yet.

## Tribute

A massive thank you to the original author [@noirbizarre](https://github.com/noirbizarre) for doing the initial work in https://github.com/noirbizarre/gonja which this project was forked from.