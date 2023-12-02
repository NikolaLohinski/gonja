{% set parsed = "nested:\n  value: 123\n  foo: bar" | fromyaml %}
{{ parsed.nested.value }}
{{ parsed.nested.foo }}