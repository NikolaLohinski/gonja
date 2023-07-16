{{ "foo bar baz qux"|truncate(100) }}
{{ "foo bar baz qux"|truncate(9) }}
{{ "foo bar baz qux"|truncate(9, leeway=5) }}
{{ "foo bar baz qux"|truncate(9, True) }}
{{ "foo bar baz qux"|truncate(9, True, '…') }}
{{ "foo bar baz qux"|truncate(15) }}
{{ "foo bar baz qux"|truncate(11, leeway=5) }}
{{ "foo bar baz qux"|truncate(11, False) }}
