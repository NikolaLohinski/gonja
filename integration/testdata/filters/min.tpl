{{ [1, 2, 3]|min }}
{{ ["a", "B"]|min }}
{{ ["a", "B"]|min(case_sensitive=true) }}
{{ []|min }}
{{ [{'value': 1}, {'value': 2}]|min(attribute="value") }}
