Begin
{% macro greetings(to, from=simple.name, name2="guest") %}
Greetings to {{ to }} from {{ from }}. Howdy, {% if name2 == "guest" %}anonymous guest{% else %}{{ name2 }}{% endif %}!
{% endmacro %}
{{ greetings('') }}
{{ greetings(10) }}
{{ greetings("john") }}
{{ greetings("john", "michelle") }}
{{ greetings("john", "michelle", "johann") }}

{% macro test2(loop, value) %}map[{{ loop.index0 }}] = {{ value }}{% endmacro %}
{% for item in simple.misc_list %}
{{ test2(loop, item) }}{% endfor %}

issue #39 (deactivate auto-escape of macros)
{% macro html_test(name) %}
<p>Hello {{ name }}.</p>
{% endmacro %}
{{ html_test("Max") }}

Importing macros
{% from "macro.helper" import imported_macro, imported_macro as renamed_macro, imported_macro as html_test %}
{{ imported_macro("User1") }}
{{ renamed_macro("User2") }}
{{ html_test("Max") }}

Chaining macros{% from "macro2.helper" import greeter_macro %}
{{ greeter_macro() }}
End
