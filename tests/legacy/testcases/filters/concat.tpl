{%- set one = [] | concat(["one"]) -%}
{%- set multiple = one | concat(["two"], ["three"]) -%}
{{ one }}
{{ multiple }}