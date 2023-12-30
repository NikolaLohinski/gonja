{% if nothing %}false{% else %}true{% endif %}
{% if simple %}simple != nil{% endif %}
{% if simple.uint %}uint != 0{% endif %}
{% if simple.float %}float != 0.0{% endif %}
{% if not simple %}false{% else %}!simple{% endif %}
{% if not simple.uint %}false{% else %}!simple.uint{% endif %}
{% if not simple.float %}false{% else %}!simple.float{% endif %}
{% if "Text" in complex.post %}text field in complex.post{% endif %}
{% if 5 in simple.intmap %}5 in simple.intmap{% endif %}
{% if not 0.0 %}!0.0{% endif %}
{% if not 0 %}!0{% endif %}
{% if not complex.post %}true{% else %}false{% endif %}
{% if simple.number == 43 %}no{% else %}42{% endif %}
{% if simple.number < 42 %}false{% elif simple.number > 42 %}no{% elif simple.number >= 42 %}yes{% else %}no{% endif %}
{% if simple.number < 42 %}false{% elif simple.number > 42 %}no{% elif simple.number != 42 %}no{% else %}yes{% endif %}
{% if 0 %}!0{% elif nothing %}nothing{% else %}true{% endif %}
{% if 0 %}!0{% elif simple.float %}simple.float{% else %}false{% endif %}
{% if 0 %}!0{% elif not simple.float %}false{% elif "Text" in complex.post%}Elseif with no else{% endif %}
{% if simple.str != '' and 'foo' not in simple.str | lower %}https://github.com/nikolalohinski/gonja/v2/issues/7{% endif %}