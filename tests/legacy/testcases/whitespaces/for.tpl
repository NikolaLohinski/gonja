##########
{% for x in [1,2,3] %}
{{ x }}
{% endfor %}
##########
{% for x in [1,2,3] -%}
{{ x }}
{% endfor %}
##########
{% for x in [1,2,3] %}
{{ x }}
{%- endfor %}
##########
{% for x in [1,2,3] -%}
{{ x }}
{%- endfor %}
##########
{% for x in [1,2,3] %}
{{- x }}
{% endfor %}
##########
{% for x in [1,2,3] %}
{{- x -}}
{% endfor %}
##########
{% for x in [1,2,3] -%}
{{ x }}
{%- if x is even -%}
(is modulo)
{%- endif %}
{% endfor %}
##########