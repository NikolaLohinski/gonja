sort the dict by key, case insensitive
{% for key, value in simple.casedStrmap|dictsort -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by key, case insensitive, reverse order
{% for key, value in simple.casedStrmap|dictsort(reverse=true) -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by key, case sensitive
{% for key, value in simple.casedStrmap|dictsort(true) -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by key, case sensitive, reverse order
{% for key, value in simple.casedStrmap|dictsort(true, reverse=true) -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by value, case insensitive
{% for key, value in simple.casedStrmap|dictsort(by='value') -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by value, case insensitive, reverse order
{% for key, value in simple.casedStrmap|dictsort(by='value', reverse=true) -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by value, case sensitive
{% for key, value in simple.casedStrmap|dictsort(true, 'value') -%}
{{ key }}: {{ value }}
{% endfor %}

sort the dict by value, case sensitive, reverse order
{% for key, value in simple.casedStrmap|dictsort(true, 'value', reverse=true) -%}
{{ key }}: {{ value }}
{% endfor %}
