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

Most methods are implemented; Except for these:
- [maketrans](https://docs.python.org/3/library/stdtypes.html#str.maketrans) - complex
- [translate](https://docs.python.org/3/library/stdtypes.html#str.translate) - complex
- [isidentifier](https://docs.python.org/3/library/stdtypes.html#str.isidentifier) - not too relevant outside python.
- [index](https://docs.python.org/3/library/stdtypes.html#str.index) - have find without errors.
- [rindex](https://docs.python.org/3/library/stdtypes.html#str.rindex) - Have rfind without errors


We strive for 100% python compatibility but due to limited resources we are bound by the differences in golang and python. 
For example in go `unicode.IsNumber("三")` returns false for chinese numbers "三" but in python `!char.isnumeric()` returns true.

### The [capitalize](https://docs.python.org/3/library/stdtypes.html#str.capitalize)() method

Return a copy of the string with its first character capitalized and the rest lowercased.

### The [capwords](https://docs.python.org/3/library/string.html#string.capwords)()

Split the argument into words using str.split(), capitalize each word using str.capitalize(), and join the capitalized words using str.join(). If the optional second argument sep is absent or None, runs of whitespace characters are replaced by a single space and leading and trailing whitespace are removed, otherwise sep is used to split and join the words.

### The [casefold](https://docs.python.org/3/library/stdtypes.html#str.casefold)() method

Return a casefolded copy of the string. Casefolded strings may be used for caseless matching.

Casefolding is similar to lowercasing but more aggressive because it is intended to remove all case distinctions in a string. For example, the German lowercase letter 'ß' is equivalent to "ss". Since it is already lowercase, lower() would do nothing to 'ß'; casefold() converts it to "ss".

### The [center](https://docs.python.org/3/library/stdtypes.html#str.center)(width[, fillchar]) method

Return centered in a string of length width. Padding is done using the specified fillchar (default is an ASCII space). The original string is returned if width is less than or equal to len(s).

### The [count](https://docs.python.org/3/library/stdtypes.html#str.count)(sub[, start[, end]]) method

Return the number of non-overlapping occurrences of substring sub in the range [start, end]. Optional arguments start and end are interpreted as in slice notation.

If sub is empty, returns the number of empty strings between characters which is the length of the string plus one.

### The [encode](https://docs.python.org/3/library/stdtypes.html#str.encode)(encoding='utf-8', errors='strict') method

-Return the string encoded to bytes. Encoding defaults to `'utf-8'` and only `'iso-8859-1'` is also supported. `errors` controls how encoding errors are handled and can be set to `'strict'` or `'ignore'`.

### The [endswith](https://docs.python.org/3/library/stdtypes.html#str.endswith)(suffix[, start[, end]]) method

Return True if the string ends with the specified suffix, otherwise return False. suffix can also be a tuple of suffixes to look for. With optional start, test beginning at that position. With optional end, stop comparing at that position.

### The [expandtabs](https://docs.python.org/3/library/stdtypes.html#str.expandtabs)(tabsize=8) method

Return a copy of the string where all tab characters are replaced by one or more spaces, depending on the current column and the given tab size. Tab positions occur every tabsize characters (default is 8, giving tab positions at columns 0, 8, 16 and so on). To expand the string, the current column is set to zero and the string is examined character by character. If the character is a tab (\t), one or more space characters are inserted in the result until the current column is equal to the next tab position. (The tab character itself is not copied.) If the character is a newline (\n) or return (\r), it is copied and the current column is reset to zero. Any other character is copied unchanged and the current column is incremented by one regardless of how the character is represented when printed.

### The [find](https://docs.python.org/3/library/stdtypes.html#str.find)(sub[, start[, end]]) method

Return the lowest index in the string where substring sub is found within the slice s[start:end]. Optional arguments start and end are interpreted as in slice notation. Return -1 if sub is not found.

### The [format](https://docs.python.org/3/library/stdtypes.html#str.format)(*args, **kwargs) method

Perform a string formatting operation. The string on which this method is called can contain literal text or replacement fields delimited by braces {}. Each replacement field contains either the numeric index of a positional argument, or the name of a keyword argument. Returns a copy of the string where each replacement field is replaced with the string value of the corresponding argument.

See Format [String Syntax](https://docs.python.org/3/library/string.html#formatstrings) for a description of the various formatting options that can be specified in format strings.

There are differences in python versions. We try to capture this with ["dialects" and default to 3.11](https://github.com/NikolaLohinski/gonja/blob/master/builtins/methods/pystring/dialect.go). Override the DefaultDialect to get the desired behavior. 


### The [format_map](https://docs.python.org/3/library/stdtypes.html#str.format_map)(mapping) method

Aliased to format since the difference in golang is not relevant.

### The [isalnum](https://docs.python.org/3/library/stdtypes.html#str.isalnum)() method

Return True if all characters in the string are alphanumeric and there is at least one character, False otherwise. A character c is alphanumeric if one of the following returns True: c.isalpha(), c.isdecimal(), c.isdigit(), or c.isnumeric().

### The [isalpha](https://docs.python.org/3/library/stdtypes.html#str.isalpha)() method

Return True if all characters in the string are alphabetic and there is at least one character, False otherwise. Alphabetic characters are those characters defined in the Unicode character database as “Letter”, i.e., those with general category property being one of “Lm”, “Lt”, “Lu”, “Ll”, or “Lo”. Note that this is different from the Alphabetic property defined in the section 4.10 ‘Letters, Alphabetic, and Ideographic’ of the Unicode Standard.

### The [isascii](https://docs.python.org/3/library/stdtypes.html#str.isascii)() method

Return True if the string is empty or all characters in the string are ASCII, False otherwise. ASCII characters have code points in the range U+0000-U+007F.

### The [isdecimal](https://docs.python.org/3/library/stdtypes.html#str.isdecimal)() method

Return True if all characters in the string are decimal characters and there is at least one character, False otherwise. Decimal characters are those that can be used to form numbers in base 10, e.g. U+0660, ARABIC-INDIC DIGIT ZERO. Formally a decimal character is a character in the Unicode General Category “Nd”.

### The [isdigit](https://docs.python.org/3/library/stdtypes.html#str.isdigit)() method

Return True if all characters in the string are digits and there is at least one character, False otherwise. Digits include decimal characters and digits that need special handling, such as the compatibility superscript digits. This covers digits which cannot be used to form numbers in base 10, like the Kharosthi numbers. Formally, a digit is a character that has the property value Numeric_Type=Digit or Numeric_Type=Decimal.

### The [islower](https://docs.python.org/3/library/stdtypes.html#str.islower)() method

Return True if all cased characters [4] in the string are lowercase and there is at least one cased character, False otherwise.

### The [isnumeric](https://docs.python.org/3/library/stdtypes.html#str.isnumeric)() method

Return True if all characters in the string are numeric characters, and there is at least one character, False otherwise. Numeric characters include digit characters, and all characters that have the Unicode numeric value property, e.g. U+2155, VULGAR FRACTION ONE FIFTH. Formally, numeric characters are those with the property value Numeric_Type=Digit, Numeric_Type=Decimal or Numeric_Type=Numeric.

### The [isprintable](https://docs.python.org/3/library/stdtypes.html#str.isprintable)() method

Return True if all characters in the string are printable or the string is empty, False otherwise. Nonprintable characters are those characters defined in the Unicode character database as “Other” or “Separator”, excepting the ASCII space (0x20) which is considered printable. (Note that printable characters in this context are those which should not be escaped when repr() is invoked on a string. It has no bearing on the handling of strings written to sys.stdout or sys.stderr.)

### The [isspace](https://docs.python.org/3/library/stdtypes.html#str.isspace)() method

Return True if there are only whitespace characters in the string and there is at least one character, False otherwise.

A character is whitespace if in the Unicode character database (see unicodedata), either its general category is Zs (“Separator, space”), or its bidirectional class is one of WS, B, or S.

### The [istitle](https://docs.python.org/3/library/stdtypes.html#str.istitle)() method

Return True if the string is a titlecased string and there is at least one character, for example uppercase characters may only follow uncased characters and lowercase characters only cased ones. Return False otherwise.

### The [isupper](https://docs.python.org/3/library/stdtypes.html#str.isupper)() method

Return True if all cased characters [4] in the string are uppercase and there is at least one cased character, False otherwise.

### The [join](https://docs.python.org/3/library/stdtypes.html#str.join)(iterable) method

Return a string which is the concatenation of the strings in iterable. A TypeError will be raised if there are any non-string values in iterable, including bytes objects. The separator between elements is the string providing this method.

### The [ljust](https://docs.python.org/3/library/stdtypes.html#str.ljust)(width[, fillchar]) method

Return the string left justified in a string of length width. Padding is done using the specified fillchar (default is an ASCII space). The original string is returned if width is less than or equal to len(s).

### The [lower](https://docs.python.org/3/library/stdtypes.html#str.lower)() method

Return a copy of the string with all the cased characters [4] converted to lowercase.

The lowercasing algorithm used is described in section 3.13 ‘Default Case Folding’ of the Unicode Standard.

### The [lstrip](https://docs.python.org/3/library/stdtypes.html#str.lstrip)([chars]) method

Return a copy of the string with leading characters removed. The chars argument is a string specifying the set of characters to be removed. If omitted or None, the chars argument defaults to removing whitespace. The chars argument is not a prefix; rather, all combinations of its values are stripped:

### The [partition](https://docs.python.org/3/library/stdtypes.html#str.partition)(sep) method

Split the string at the first occurrence of sep, and return a 3-tuple containing the part before the separator, the separator itself, and the part after the separator. If the separator is not found, return a 3-tuple containing the string itself, followed by two empty strings.

### The [removeprefix](https://docs.python.org/3/library/stdtypes.html#str.removeprefix)(prefix, /) method

If the string starts with the prefix string, return string[len(prefix):]. Otherwise, return a copy of the original string:

### The [removesuffix](https://docs.python.org/3/library/stdtypes.html#str.removesuffix)(suffix, /) method

If the string ends with the suffix string and that suffix is not empty, return string[:-len(suffix)]. Otherwise, return a copy of the original string:

### The [replace](https://docs.python.org/3/library/stdtypes.html#str.replace)(old, new[, count]) method

Return a copy of the string with all occurrences of substring old replaced by new. If the optional argument count is given, only the first count occurrences are replaced.

### The [rfind](https://docs.python.org/3/library/stdtypes.html#str.rfind)(sub[, start[, end]]) method

Return the highest index in the string where substring sub is found, such that sub is contained within s[start:end]. Optional arguments start and end are interpreted as in slice notation. Return -1 on failure.

### The [rjust](https://docs.python.org/3/library/stdtypes.html#str.rjust)(width[, fillchar]) method

Return the string right justified in a string of length width. Padding is done using the specified fillchar (default is an ASCII space). The original string is returned if width is less than or equal to len(s).

### The [rpartition](https://docs.python.org/3/library/stdtypes.html#str.rpartition)(sep) method

Split the string at the last occurrence of sep, and return a 3-tuple containing the part before the separator, the separator itself, and the part after the separator. If the separator is not found, return a 3-tuple containing two empty strings, followed by the string itself.

### The [rsplit](https://docs.python.org/3/library/stdtypes.html#str.rsplit)(sep=None, maxsplit=-1) method

Return a list of the words in the string, using sep as the delimiter string. If maxsplit is given, at most maxsplit splits are done, the rightmost ones. If sep is not specified or None, any whitespace string is a separator. Except for splitting from the right, rsplit() behaves like split() which is described in detail below.

### The [rstrip](https://docs.python.org/3/library/stdtypes.html#str.rstrip)([chars]) method

Return a copy of the string with trailing characters removed. The chars argument is a string specifying the set of characters to be removed. If omitted or None, the chars argument defaults to removing whitespace. The chars argument is not a suffix; rather, all combinations of its values are stripped

### The [split](https://docs.python.org/3/library/stdtypes.html#str.split)(sep=None, maxsplit=-1) method

Return a list of the words in the string, using sep as the delimiter string. If maxsplit is given, at most maxsplit splits are done (thus, the list will have at most maxsplit+1 elements). If maxsplit is not specified or -1, then there is no limit on the number of splits (all possible splits are made).

If sep is given, consecutive delimiters are not grouped together and are deemed to delimit empty strings (for example, '1,,2'.split(',') returns ['1', '', '2']). The sep argument may consist of multiple characters as a single delimiter (to split with multiple delimiters, use re.split()). Splitting an empty string with a specified separator returns [''].

### The [splitlines](https://docs.python.org/3/library/stdtypes.html#str.splitlines)(keepends=False) method

Return a list of the lines in the string, breaking at line boundaries. Line breaks are not included in the resulting list unless keepends is given and true.

This method splits on the following line boundaries. In particular, the boundaries are a superset of universal newlines.

### The [startswith](https://docs.python.org/3/library/stdtypes.html#str.startswith)(prefix[, start[, end]]) method

Return True if string starts with the prefix, otherwise return False. prefix can also be a tuple of prefixes to look for. With optional start, test string beginning at that position. With optional end, stop comparing string at that position.

### The [strip](https://docs.python.org/3/library/stdtypes.html#str.strip)([chars]) method

Return a copy of the string with the leading and trailing characters removed. The chars argument is a string specifying the set of characters to be removed. If omitted or None, the chars argument defaults to removing whitespace. The chars argument is not a prefix or suffix; rather, all combinations of its values are stripped

### The [swapcase](https://docs.python.org/3/library/stdtypes.html#str.swapcase)() method

Return a copy of the string with uppercase characters converted to lowercase and vice versa. Note that it is not necessarily true that s.swapcase().swapcase() == s.

### The [title](https://docs.python.org/3/library/stdtypes.html#str.title)() method

Return a titlecased version of the string where words start with an uppercase character and the remaining characters are lowercase.

### The [upper](https://docs.python.org/3/library/stdtypes.html#str.upper)() method

Return a copy of the string with all the cased characters [4] converted to uppercase. Note that s.upper().isupper() might be False if s contains uncased characters or if the Unicode category of the resulting character(s) is not “Lu” (Letter, uppercase), but e.g. “Lt” (Letter, titlecase).

The uppercasing algorithm used is described in section 3.13 ‘Default Case Folding’ of the Unicode Standard.

### The [zfill](https://docs.python.org/3/library/stdtypes.html#str.zfill)(width) method

Return a copy of the string left filled with ASCII '0' digits to make a string of length width. A leading sign prefix ('+'/'-') is handled by inserting the padding after the sign character rather than before. The original string is returned if width is less than or equal to len(s).


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
