{% for comment in complex.comments %}[{{ loop.index }} {{ loop.index0 }} {{ loop.first }} {{ loop.last }} {{ loop.revindex }} {{ loop.revindex0 }}] {{ comment.Author.Name }}

{# nested loop #}
{% set parent = loop %}
{% for char in comment.Text %}{{ parent.index0 }}.{{ loop.index0 }}:{{ char|safe }} {% endfor %}

{% endfor %}

reversed
'{% for item in simple.multiple_item_list|reverse %}{{ item }} {% endfor %}'

sorted string map
'{% for key in simple.strmap|sort %}{{ key }} {% endfor %}'

sorted int map
'{% for key in simple.intmap|sort %}{{ key }} {% endfor %}'

sorted int list
'{% for key in simple.unsorted_int_list|sort %}{{ key }} {% endfor %}'

reversed sorted int list
'{% for key in simple.unsorted_int_list|sort|reverse %}{{ key }} {% endfor %}'

reversed sorted string map
'{% for key in simple.strmap|sort|reverse %}{{ key }} {% endfor %}'

reversed sorted int map
'{% for key in simple.intmap|sort|reverse %}{{ key }} {% endfor %}'

key, value sorted string map
{%- for key, value in simple.strmap|dictsort %}
{{ key }}: {{ value }}
{%- endfor %}

(key, value) 2-list key
{%- for key, value in [['key', 'value'], ['2nd key', '2nd value']] %}
{{ key }}: {{ value }}
{%- endfor %}

(key, value) 2-tuple key
{%- for key, value in (('key', 'value'), ('2nd key', '2nd value')) %}
{{ key }}: {{ value }}
{%- endfor %}

If expression
{%- for person in persons if person.Gender is equalto "male" %}
{{ person.FirstName }} {{ person.LastName }}
{%- endfor %}

If expression and loop
{%- for person in persons if person.Gender is equalto "male" %}
{{ person.FirstName }} {{ person.LastName }} {{ loop.index }} {{ loop.revindex }} {{ loop.first }} {{ loop.last }}
{%- endfor %}

Cycle
{%- for idx in range(4) %}
{{ idx }} {{ loop.Cycle('even', 'odd') }}
{%- endfor %}

Else
{%- for idx in [] %}
Should not go there
{%- else %}
Nothing
{%- endfor %}

Changed
{%- for idx in [1, 2, 2, 3]  %}
{{ idx }}: {{ loop.Changed(idx) }}
{%- endfor %}

Prev/Next items
{%- for idx in range(3)  %}
{{ idx }}: prev: {{ loop.PrevItem }} next: {{ loop.NextItem }}
{%- endfor %}

Prev/Next items 2-tuple
{%- for k, v in [(1, 'first'), (2, 'second'), (3, 'third')]  %}
{{ k }} {{ v }}: prev: {{ loop.PrevItem }} next: {{ loop.NextItem }}
{%- endfor %}

Prev/Next items with if
{%- for idx in range(6) if idx is even  %}
{{ idx }}: prev: {{ loop.PrevItem }} next: {{ loop.NextItem }}
{%- endfor %}

Else with if
{% set keywords = ["yes", "hello"] -%}
{%- for w in keywords if w in ["no"] -%}
{%- if loop.first -%}
match
{%- endif -%}
{%- else -%}
no match
{%- endfor -%}