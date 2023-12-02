{%- set var = '{ "string": "one","list": ["item"], "nested": { "field": 123 }}' | fromjson -%}
{{ var.string }}
{{ var.list }}
{{ var.nested.field }}