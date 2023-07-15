{% macro input(name, value='', type='text', size=20) -%}
    <input type="{{ type }}" name="{{ name }}" value="{{
        value|e }}" size="{{ size }}">
{%- endmacro -%}
<p>{{ input('username') }}</p>
<p>{{ input('username', type='password') }}</p>
<p>{{ input(name='address', size=42) }}</p>
