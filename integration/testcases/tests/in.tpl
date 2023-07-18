o/foo: {{ "o" is in "foo" }}
foo/foo: {{ "foo" is in "foo" }}
b/foo: {{ "b" is in "foo" }}
1/(1,2): {{ 1 is in (1, 2) }}
3/(1,2): {{ 3 is in (1, 2) }}
1/[1,2]: {{ 1 is in [1, 2] }}
3/[1,2]: {{ 3 is in [1, 2] }}
foo/{"foo": 1}: {{ "foo" is in {"foo": 1}}}
baz/{"bar": 1}: {{ "baz" is in {"bar": 1}}}
