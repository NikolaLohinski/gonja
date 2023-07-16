Start '{% include "includes.helper" %}' End
Start '{% include "includes.helper" ignore missing %}' End
{% with number=7, what_am_i="guest" -%}
Start '{% include simple.included_file|lower %}' End
{%- endwith %}
Start '{% include "includes.helper.not_exists" ignore missing %}' End
{% with number=7, what_am_i="guest" -%}
Start '{% include simple.included_file_not_exists ignore missing %}' End
{%- endwith %}
