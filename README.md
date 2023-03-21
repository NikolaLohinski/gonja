# Gonja

`gonja` is pure `go` implementation of the [Jinja](https://jinja.palletsprojects.com/en/3.1.x/) template engine. It aims to be _mostly_ compatible with the original `python` implementation but also provides additional features to compensate the lack of `python` scripting capabilities.

## Usage

### As a library

Install/update using `go get`:
```
go get github.com/nikolalohinski/gonja
```

### As a `terraform` provider

This `gonja` library has been packaged as a `terraform` provider. For more information, please refer to the [dedicated documentation](https://registry.terraform.io/providers/NikolaLohinski/jinja/latest/docs).

## Example

```golang
// TODO
```

## Documentation

* For a details on how the templating language works, please refer to [the Jinja documentation](https://jinja.palletsprojects.com) ;
* `gonja` API documentation is available on [godoc](https://godoc.org/github.com/nikolalohinski/gonja) ;
* `filters`: please refer to [`docs/filters.md`](docs/filters.md) ;
* `statements`: please refer to [`docs/statements.md`](docs/statments.md) ;
* `tests`: please refer to [`docs/tests.md`](docs/tests.md) ;
* `globals`: please refer to [`docs/globals.md`](docs/globals.md).

## Known caveats 

### Filters

 * **format**: `format` does **not** take Python's string format syntax as a parameter, instead it takes Go's. Essentially `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`.
 * **escape** / **force_escape**: Unlike Jinja's behaviour, the `escape`-filter is applied immediately. Therefore there is no need for a `force_escape`-filter yet.

### Panic attacks

### Code hanging

## Tribute

A massive thank you to the original author [@noirbizarre](https://github.com/noirbizarre) for doing the initial work in https://github.com/noirbizarre/gonja which this project was forked from.