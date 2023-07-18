<ul>
{%- for group in persons|groupby('Gender') %}
    <li>{{ group.grouper }}<ul>
    {%- for person in group.list %}
        <li>{{ person.FirstName }} {{ person.LastName }}</li>
    {%- endfor %}
    </ul></li>
{%- endfor %}
</ul>

<ul>
{%- for group in groupable|groupby('grouper') %}
    <li>{{ group.grouper }}<ul>
    {%- for value in group.list %}
        <li>{{ value.value }}</li>
    {%- endfor %}
    </ul></li>
{%- endfor %}
</ul>
