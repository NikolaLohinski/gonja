<table>
{%- for row in [0, 1, 2, 3, 4, 5, 6]|batch(3) %}
  <tr>
  {%- for column in row %}
    <td>{{ column }}</td>
  {%- endfor %}
  </tr>
{%- endfor %}
</table>

<table>
{%- for row in [0, 1, 2, 3, 4, 5, 6]|batch(3, 'filled') %}
  <tr>
  {%- for column in row %}
    <td>{{ column }}</td>
  {%- endfor %}
  </tr>
{%- endfor %}
</table>
