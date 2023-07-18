{% set pipe = joiner("|") -%}
{% for i in [0, 1, 2] %}{{ pipe() }}{{ i }}{% endfor %}
