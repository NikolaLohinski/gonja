{% set ns = namespace(found=false) -%}
{%- for item in [1, 2, 3] -%}
    {%- if item == 3 -%}{% set ns.found = true -%}{%- endif -%}
    * {{ item }}
{% endfor -%}
Found item having something: {{ ns.found }}
