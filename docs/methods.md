# Methods

Methods are Go implementations of class functions available on native `python` types. For example:

```
{{ "hello".upper() }}
```

The following clickable admonition can be used to browse the `python` dedicated documentation for additional details:

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html) |
| ------------------------------------------------------------- |


## The `bool` type      

_Booleans are subtypes of integers in `python`._

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html#additional-methods-on-integer-types) |
| ------------------------------------------------------------------------------------------------- |

### The `bit_length()` method

Returns the number of bits necessary to represent an integer in binary, excluding the sign and leading zeros.

### The `bit_count()` method

Returns the number of ones in the binary representation of the absolute value of the integer. This is also known as the population count.

## The `int` type      

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html#additional-methods-on-integer-types) |
| ------------------------------------------------------------------------------------------------- |

### The `is_integer()` method

Returns `True`. Exists for duck type compatibility with `float.is_integer()`.


## The `float` type      

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html#additional-methods-on-float) |
| ----------------------------------------------------------------------------------------- |

### The `is_integer()` method

Returns `True` if the float instance is finite with integral value, and `False` otherwise.

## The `str` type

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html#string-methods) |
| ---------------------------------------------------------------------------- |

### The `upper()` method

Returns a copy of the string with all the cased characters converted to uppercase.

### The `startswith(prefix[, start[, end]])` method

Returns `True` if string starts with the prefix, otherwise return `False`. The `prefix` parameter can also be a tuple (or a list which is not supported in `python`) of prefixes to look for. 

With optional `start`, test string beginning at that position. With optional `end`, stop comparing string at that position.

### The `encode(encoding='utf-8', errors='strict')` method

Return the string encoded to bytes. Encoding defaults to `'utf-8'` and only `'iso-8859-1'` is also supported. `errors` controls how encoding errors are handled and can be set to `'strict'` or `'ignore'`.

## The `list` type      

| [🐍 `python`](https://docs.python.org/3/tutorial/datastructures.html#data-structures) |
| ------------------------------------------------------------------------------------ |

### The `reverse()` method

Reverses the elements of the list in place.

### The `append(x)` method

Adds an item to the end of the list.

### The `copy()` method

Returns a shallow copy of the list.

## The `dict` type      

| [🐍 `python`](https://docs.python.org/3/library/stdtypes.html#mapping-types-dict) |
| -------------------------------------------------------------------------------- |

### The `keys()` method

Returns a list of the dictionary’s keys.