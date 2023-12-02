# precision 0
{{ 42.55|round }}
{{ 42.55|round(0, 'ceil') }}
{{ 42.55|round(0, 'floor') }}
{{ 42.45|round }}
{{ 42.45|round(0, 'ceil') }}
{{ 42.45|round(0, 'floor') }}
# precision 1
{{ 42.55|round(1) }}
{{ 42.55|round(1, 'ceil') }}
{{ 42.55|round(1, 'floor') }}
{{ 42.54|round(1) }}
{{ 42.54|round(1, 'ceil') }}
{{ 42.54|round(1, 'floor') }}
