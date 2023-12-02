{{ [1, 2, 3]|max }}
{{ ["a", "B"]|max }}
{{ ["a", "B"]|max(case_sensitive=true) }}
{{ []|max }}
{{ [{'value': 1}, {'value': 2}]|max(attribute="value") }}
