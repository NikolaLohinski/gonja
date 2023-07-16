{{ {} | tojson }}
{{ ['one'] | tojson }}
{{ True | tojson }}
{{ 'test' | tojson }}
{{ simple.strmap | tojson }}
{{ simple.strmap | tojson(4) }}
{{ simple.strmap | tojson(indent=2) }}
