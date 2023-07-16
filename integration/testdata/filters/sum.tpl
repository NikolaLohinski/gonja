{{ [1, 2, 3, 4, 5, 6]|sum }}
{{ [1, 2, 3, 4, 5, 6]|sum(start=10) }}
{{ [1, 2, 3, 4, 5, 6.1]|sum }}
{{ [{'value': 23}, {'value': 1}, {'value': 18}]|sum('value') }}
{{ [{'real': {'value': 23}}, {'real': {'value': 1}}, {'real': {'value': 18}}]|sum('real.value') }}
{{ [('foo', 23), ('bar', 1), ('baz', 18)]|sum(1) }}
