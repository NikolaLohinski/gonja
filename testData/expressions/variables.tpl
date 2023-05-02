{{ 1 }}
{{ -5 }}
{{ "hallo" }}
{{ true }}
{{ false }}
{{ None }}
{{ nil }}
{{ simple.uint }}
{%- set str = "UINT" %}
{{ simple[str | lower] }}
{{ simple["uint"] }}
{{ simple.nil }}
{{ simple.str }}
{{ simple.bool_false }}
{{ simple.bool_true }}
{{ simple.uint }}
{{ simple.uint|int }}
{{ simple.uint|float }}
{{ simple.multiple_item_list.10 }}
{{ simple.multiple_item_list.4 }}
{{ simple["missing"] }}