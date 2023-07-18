{%- if (fail | try) is undefined %}
Now you see me without errors!
{%- endif %}
{%- if (no | try) %}
But here you don't!
{%- endif %}