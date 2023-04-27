{{ "simple if when true" if "foo" is string }}
{{ "simple if when false" if "nope" is number }}
{{ "if and else when true" if 2 is number else "never" }}
{{ "never" if 2 is odd else "if and else when false" }}
{{ "never" if None else "when condition is nil" }}
{{ "never" if '' else "when condition is empty string" }}
{{ "never" if 0 else "when condition is 0" }}