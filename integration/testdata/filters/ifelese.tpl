true  = {{ "foo" in "foo bar" | ifelse("yes", "no") }}
false = {{ "baz" in "foo bar" | ifelse("yes", {"value": "no"}) }}