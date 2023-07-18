{{ "foo" | match("bar") }}
{{ "123" | match("^[0-9]+$") }}
{{ "nope" | match("^[0-9]+$") }}