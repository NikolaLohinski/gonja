{% set row_class = cycler('odd', 'even') -%}
<ul class="browser">
{%- for num in [1, 2, 3] %}
  <li class="{{ row_class.next() }}">{{ num }}</li>
{%- endfor %}
{%- for num in simple.fixed_item_list %}
  <li class="{{ row_class.next() }}">{{ num }}</li>
{%- endfor %}
</ul>
