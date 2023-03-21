## Tests

A test can be used in blocks and/or expressions to trigger conditional behavior, for example:

```
{% if variable is string %}
   This was a string: {{ variable }}
{% elif variable is sequence %}
   This was a list: {{ variable | join(",") }}
{% end if%}
```

Any test that is also implemented in the `python` version of the Jinja engine will be marked with the following tag:

| 🐍 `python` |
|-------------|

For any of those, the [official documentation for `python`'s Jinja implementation](https://jinja.palletsprojects.com/en/3.1.x/templates/#list-of-builtin-tests)  can be used as additional reference.


### The `callable` test
| 🐍 `python` |
|-------------|

Return whether the object is callable (i.e., some kind of function).

### The `defined` and `undefined` tests
| 🐍 `python` |
|-------------|

Tells whether a variable is `defined` or `undefined`.

### The `divisibleby` test
| 🐍 `python` |
|-------------|

Check if a variable is divisible by a number.
```
{% if 2048 is divisibleby 512 %}
    Yes it is modulo 4
{% endif %}
```

### The `eq`/`equalto`/`==` and `ne`/`!=` tests
| 🐍 `python` |
|-------------|

Classic arithmetic equality and inequality comparisons.

### The `ge`/`>=`, `gt`/`>`, `le`/`<=` and `lt`/`<` tests
| 🐍 `python` |
|-------------|

Classic arithmetic comparisons.

### The `even` and `odd` tests
| 🐍 `python` |
|-------------|

Tells whether a given number can be divided by 2 (`even`) or not (`odd`).

### The `in` test
| 🐍 `python` |
|-------------|

Return whether the input contains the argument:
* on strings, tells whether the provided substring is part of the tested one ;
* on lists, tells whether the argument in the tested list ;
* on dictionaries, tells whether the argument is a key of the dictionary.
```
{{ "foo" is in "foobar" }}            // True
{{ 4 is in [1, 2, 3] }}               // False
{{ "key" is in {"key": "value"} }}    // True
{{ "value" is in {"key": "value"} }}  // False
```

### The `iterable` tests
| 🐍 `python` |
|-------------|

Check if it’s possible to iterate over the tested input, i.e the object is either a list, a dictionary or a string.

### The `empty` test
Check if the input is empty. Works on strings, lists and dictionaries.

### The `none` test
| 🐍 `python` |
|-------------|

Return `True` if the input is `nil` or `None`

### The `mapping`,`sequence`, `number` and `string` tests
| 🐍 `python` |
|-------------|

Classic type casting tests.