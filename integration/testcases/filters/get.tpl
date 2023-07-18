{%- set dictionary = {"field": "content"} -%}
{%- set key = "field" -%}
{{- dictionary | get(key) -}}
{{- dictionary | get("not strict, should just print empty") -}}
{{- dictionary | get("with default", default=" default") -}}