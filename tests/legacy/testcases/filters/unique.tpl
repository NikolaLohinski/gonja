{{ ['foo', 'bar', 'foobar', 'FooBar', 'foobar', 'FooBar']|unique }}
{{ ['foo', 'bar', 'foobar', 'FooBar', 'foobar', 'FooBar']|unique(true) }}
{{ [1, 2, 3, 3, 2, 1]|unique }}
{{ [1, 2, 3, 3, 2, 1]|unique(true) }}
{{ [{'value': 1}, {'value': 2}, {'value': 3}, {'value': 3}, {'value': 2}, {'value': 1}]|unique(attribute='value') }}
