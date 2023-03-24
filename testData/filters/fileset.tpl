{% for path in "files/**/*txt" | fileset %}
folder: {{ path | dir | basename }}
file: {{ path | basename }}
{%- endfor %}