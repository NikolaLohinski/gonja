{{ simple.nothing | default("n/a") }}
{{ nothing | default(simple.number) }}
{{ simple.number | default("n/a") }}
{{ 5 | d("n/a") }}
{{ false | default("false should not hit default") }}
{{ false | default("false with boolean argument should hit default", True)}}
{{ '' | default("empty string with boolean argument should hit default", True)}}