{{ "foobar"|float }}
{{ nil|float }}
{{ "5.5"|float }}
{{ 5|float }}
{{ "5.6"|int|float }}
{{ -100|float }}
{% if 5.5 == 5.500000 %}5.5 is 5.500000{% endif %}
{% if 5.5 != 5.500001 %}5.5 is not 5.500001{% endif %}
