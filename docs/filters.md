## Filters

Variables can be modified by filters. Filters are separated from the variable by a pipe symbol (|) and may have optional arguments in parentheses. Multiple filters can be chained. The output of one filter is applied to the next.

```
{{ "a,b,c" | split(",") | tojson }}

{% set flag = "off" | bool %}

{% for x in {"first": 1, "second": 2} | values %}
  {{- x }}
{% endfor %}
```

Any filter that is also implemented in the `python` version of the Jinja engine will be marked with the following tag:

| 🐍 `python` |
|-------------|

For any of those, the [official documentation for `python`'s Jinja implementation](https://jinja.palletsprojects.com/en/3.1.x/templates/#list-of-builtin-tests)  can be used as additional reference.

### The `abs` filter
| 🐍 `python` |
|-------------|

Return the absolute value of the integer or float passed.

### The `add`, `append` and `insert` filters

The `insert` filter is meant to add a key value pair to a dict. It expects a key and value to operate.

The `append` filter adds an item to a list. It expects one value to work.

The `add` filter is just the combination of both with type reflection to decide what to do.

```
{%- set object = {"existing": "value", "overridden": 123} | insert("other", true) | add("overridden", "new") -%}
{%- set array = ["one"] | add("two") | append("three") -%}
{{ object | tojson }}
{{ array | tojson }}
```
Will render into:
```
{"existing":"value","other":true,"overridden":"new"}
["one","two","three"]
```

### The `attr` filter
| 🐍 `python` |
|-------------|

Get an attribute of an object. However, items are not looked up.

### The `basename` filter

The `basename` filter is meant to be used with filesystem paths to retrieve the last element, which is usually the file name or the folder name.
```
{{ "one/two/three" | basename }}
```
Will render into:
```
three
```

### The `batch` filter
| 🐍 `python` |
|-------------|

A filter that batches items. It returns a list of lists with the given number of items and will fill missing items if the second parameter `fille_with` is passed:
```html
<table>
{%- for row in items|batch(3, '&nbsp;') %}
  <tr>
  {%- for column in row %}
    <td>{{ column }}</td>
  {%- endfor %}
  </tr>
{%- endfor %}
</table>
```

### The `bool` filter

The `bool` filter is meant to cast a string, an int, a bool or `nil` to a boolean value. _Truthful_ values are:
* Any string once lowercased such as `"on"`, `"yes"` `"1"` or `"true"`
* `1` as an integer or `1.0` as a float
* A `True` boolean

_False_ values are:
* Any string once lowercased such as `"off"`, `"no"` `"0"` or `"false"` 
* An empty string
* `0` as an integer or `0.0` as a float
* A `False` boolean
* A `nil` or `None` value

Any other type passed will cause the `bool` filter to fail.

### The `capitalize` filter
| 🐍 `python` |
|-------------|

Capitalize a value. The first character will be uppercase, all others lowercase.

### The `center` filter
| 🐍 `python` |
|-------------|

Centers the value in a field of a given width.

### The `concat` filter

The `concat` filter is meant to concatenate lists together and can take any number of lists to append together.

```
{%- set array = ["one"] | concat(["two"],["three"]) -%}
{{ array | tojson }}
```
Will render into:
```
["one","two","three"]
```

### The `default`/`d` filter
| 🐍 `python` |
|-------------|

If the value is undefined it will return the passed default value, otherwise the value of the variable:
```
{{ my_variable | d('my_variable is not defined') }}
```

### The `dictsort` filter
| 🐍 `python` |
|-------------|

Sort a dict and yield (key, value) pairs. Dictionaries may not be in the order you want to display them in, so sort them first.

### The `dir` filter
The `dir` filter is meant to be used with filesytem paths to retrieve the path to the containing folder.

```
{{ "one/two/three" | dir }}
```
Will render into:
```
one/two
```

### The `escape`/`e` filter
| 🐍 `python` |
|-------------|

Replace the characters &, <, >, ', and " in the string with HTML-safe sequences. Use this if you need to display text that might contain such characters in HTML.

### The `fail` filter

The `fail` filter is meant to error out explicitly in a given place of the template.

```
{{ "error message to output" | fail }}
```

### The `file` filter

The `file` filter is meant to load a local file into a variable. It works with both absolute and relative (to the place it's called from) paths. The `file` filter does not process the file as a template but simply loads the contents of it.

```
{% set content = "some/path" | file %}
{{ content }}
```

### The `fileset` filter

The `fileset` filter is a filesystem filter meant to be used with the `include` statement to dynamically include files.
It supports glob patterns (using `*`) and double glob patterns (using `**`) in paths, and operates relatively to the
folder that contains the file being rendered.

```
{% for path in "folder/*" | fileset %}
{% include path %}
{% endfor %}
```

### The `filesizeformat` filter
| 🐍 `python` |
|-------------|

Format the value like a ‘human-readable’ file size (i.e. 13 kB, 4.1 MB, 102 Bytes, etc).

### The `first` filter
| 🐍 `python` |
|-------------|

Return the first item of a sequence.

### The `flatten` filter

The `flatten` filter is meant to reduce a list of lists to a list of the underlying elements.

```
{%- set array = [ ["one"], ["two", "three"] ] | flatten -%}
{{ array | tojson }}
```
Will render into:
```
["one","two","three"]
```

### The `float` filter
| 🐍 `python` |
|-------------|

Convert the value into a floating point number

### The `forceescape` filter
| 🐍 `python` |
|-------------|

Enforce HTML escaping. This will probably double escape variables.

### The `format` filter
| 🐍 `python` |
|-------------|

Apply the given values to a printf-style format string, like string % values.
```
{{ "%s, %s!"|format(greeting, name) }}
Hello, World!
```

### The `fromjson` filter

The `fromjson` filter is meant to parse a JSON string into a useable object.

```
{%- set object = "{ \"nested\": { \"field\": \"value\" } }" | fromjson -%}
{{ object.nested.field }}
```
Will render into:
```
value
```

### The `fromyaml` filter

The `fromyaml` filter is meant to parse a YAML string into a useable object.

```
{%- set object = "nested:\n  field: value\n" | fromyaml -%}
{{ object.nested.field }}
```
Will render into:
```
value
```

### The `get` filter

The `get` filter helps getting an item in a map with a dynamic key:

```
{% set tac = "key" %}
{{ tic | get(tac) }}
```
With the following YAML context:
```
tic:
  key: toe
```
Will render into:
```
toe
```

The filter has the following keyword attributes:

- `strict`: a boolean to fail if the key is missing from the map. Defaults to `False` ;
- `default`: any value to pass as default if the key is not found. This takes precedence over the `strict` attribute if defined. Defaults to nil value ;

### The `groupby` filter
| 🐍 `python` |
|-------------|

Group a sequence of objects by an attribute.

For example, a list of User objects with a city attribute can be rendered in groups. In this example, grouper refers to the city value of the group.

```html
<ul>{% for city, items in users | groupby("city") %}
  <li>{{ city }}
    <ul>{% for user in items %}
      <li>{{ user.name }}
    {% endfor %}</ul>
  </li>
{% endfor %}</ul>
```

### The `ifelse` filter

The `ifelse` filter is meant to perform ternary conditions as follows:

```
true is {{ "foo" in "foo bar" | ifelse("yes", "no") }}
false is {{ "yolo" in "foo bar" | ifelse("yes", "no") }}
```
Which will render into:
```
true is yes
false is no
```

### The `indent` filter
| 🐍 `python` |
|-------------|

Return a copy of the string with each line indented by 4 spaces. The first line and blank lines are not indented by default.

### The `int` filter
| 🐍 `python` |
|-------------|

Convert the value into an integer.

### The `join` filter
| 🐍 `python` |
|-------------|

Return a string which is the concatenation of the strings in the sequence. The separator between elements is an empty string per default,

### The `keys` filter

The `keys` filter is meant to get the keys of a map as a list:

```
{{ letters | keys | sort | join(" > ") }}
```
With the following YAML context:
```
letters:
  a: hey
  b: bee
  c: see
```
Will render into:
```
a > b > c
```

Note that the order of keys is not guaranteed as there is no ordering in Golang maps.

### The `last` filter
| 🐍 `python` |
|-------------|

Return the last item of a sequence.

### The `length` filter
| 🐍 `python` |
|-------------|

Return the number of items in a container.

### The `list` filter
| 🐍 `python` |
|-------------|

Convert the value into a list. If it was a string the returned list will be a list of characters.

### The `lower` filter
| 🐍 `python` |
|-------------|

Convert a value to lowercase.

### The `map` filter

| 🐍 `python` |
|-------------|

Applies a filter on a sequence of objects or looks up an attribute. This is useful when dealing with lists of objects but you are really only interested in a certain value of it.

The basic usage is mapping on an attribute. Imagine you have a list of users but you are only interested in a list of usernames:

```
Users on this page: {{ users | map(attribute='username') | join(', ') }}
```

### The `max` and `min` filters
| 🐍 `python` |
|-------------|

Return the largest/smallest item from the sequence.

### The `pprint` filter
| 🐍 `python` |
|-------------|

Pretty print a variable.

### The `random` filter
| 🐍 `python` |
|-------------|

Return a random item from the sequence.

### The `rejectattr` filter
| 🐍 `python` |
|-------------|
Filters a sequence of objects by applying a test to the specified attribute of each object, and rejecting the objects with the test succeeding.

If no test is specified, the attribute’s value will be evaluated as a boolean.

```
{{ users | rejectattr("is_active") }}
{{ users | rejectattr("email", "none") }}
```

### The `reject` filter
| 🐍 `python` |
|-------------|

Filters a sequence of objects by applying a test to each object, and rejecting the objects with the test succeeding.

If no test is specified, each object will be evaluated as a boolean.

Example usage:

```
{{ numbers|reject("odd") }}
```

### The `replace` filter
| 🐍 `python` |
|-------------|

Return a copy of the value with all occurrences of a substring replaced with a new one.

### The `reverse` filter
| 🐍 `python` |
|-------------|

Reverse the object or return an iterator that iterates over it the other way round.

### The `round` filter
| 🐍 `python` |
|-------------|

Round the number to a given precision. The first parameter specifies the precision (default is 0), the second the rounding method:

* `common` rounds either up or down
* `ceil` always rounds up
* `floor` always rounds down

If you don’t specify a method `common` is used.

### The `safe` filter
| 🐍 `python` |
|-------------|

Mark the value as safe which means that in an environment with automatic escaping enabled this variable will not be escaped.

### The `selectattr` filter
| 🐍 `python` |
|-------------|

Filters a sequence of objects by applying a test to the specified attribute of each object, and only selecting the objects with the test succeeding.

If no test is specified, the attribute’s value will be evaluated as a boolean.

```
{{ users | selectattr("is_active") }}
{{ users | selectattr("email", "none") }}
```

### The `select` filter

| 🐍 `python` |
|-------------|
Filters a sequence of objects by applying a test to each object, and only selecting the objects with the test succeeding.

If no test is specified, each object will be evaluated as a boolean.

```
{{ numbers | select("odd") }}
{{ numbers | select("odd") }}
{{ numbers | select("divisibleby", 3) }}
{{ numbers | select("lessthan", 42) }}
{{ strings | select("equalto", "mystring") }}
```

### The `slice` filter

| 🐍 `python` |
|-------------|

Slice an iterator and return a list of lists containing those items. Useful if you want to create a div containing three ul tags that represent columns:

```
<div class="columnwrapper">
  {%- for column in items|slice(3) %}
    <ul class="column-{{ loop.index }}">
    {%- for item in column %}
      <li>{{ item }}</li>
    {%- endfor %}
    </ul>
  {%- endfor %}
</div>
```

If you pass it a second argument it’s used to fill missing values on the last iteration.

### The `sort` filter
| 🐍 `python` |
|-------------|

Sort an iterable input.

### The `split` filter

The `split` filter is meant to split a string into a list of strings using a given delimiter.

```
{%- set array = "one/two/three" | split("/") -%}
{{ array | tojson }}
```
Will render into:
```
["one","two","three"]
```

### The `string` filter
| 🐍 `python` |
|-------------|

Convert an object to a string if it isn’t already.

### The `striptags` filter
| 🐍 `python` |
|-------------|

Strip SGML/XML tags and replace adjacent whitespace by one space.

### The `sum` filter
| 🐍 `python` |
|-------------|

Returns the sum of a sequence of numbers plus the value of parameter `start` (which defaults to 0). When the sequence is empty it returns `start`.

It is also possible to sum up only certain attributes:

```
Total: {{ items | sum(attribute='price') }}
```

### The `title` filter
| 🐍 `python` |
|-------------|

Return a titlecased version of the value. I.e. words will start with uppercase letters, all remaining characters are lowercase.

### The `tojson` filter
| 🐍 `python` |
|-------------|

Serialize an object to a string of JSON. It takes an `indent` parameter to do pretty printing.

### The `toyaml` filter


The `toyaml` filter is meant to render a given object as YAML. It takes an optional argument called `indent` to
specify the indentation to apply to the result, which defaults to `2` spaces.

```
{%- set object = "{ \"nested\": { \"field\": \"value\" } }" | fromjson -%}
{{ object | toyaml }}
```
Will render into:
```
nested:
  field: value
```

### The `trim` filter
| 🐍 `python` |
|-------------|

Strip leading and trailing characters, by default whitespace.

### The `truncate` filter

| 🐍 `python` |
|-------------|

Return a truncated copy of the string. The length is specified with the first parameter which defaults to 255.

### The `try` filter

The `try` filter is meant to gracefully evaluate an expression. It returns an `undefined` value if the passed expression is undefined or throws an error. Otherwise, it returns the value passed in the context of the pipeline.

```
{%- if (empty.missing | try) is undefined -%}
	Now you see {{ value | try }}!
{%- endif -%}
```
With the following YAML context and `strict_undefined` set to `true`:
```
empty: {}
value: me
```
Will render into:
```
Now you see me!
```

This is useful when `strict_undefined = true` is set but you need to handle a missing key without throwing errors in a given template ;

### The `unique` filter

| 🐍 `python` |
|-------------|

Returns a list of unique items from the given iterable.

```
{{ ['foo', 'bar', 'foobar', 'FooBar'] | unique }}
```
Will render:
```
['foo', 'bar', 'foobar']
```

Parameters:
* case_sensitive (default: false): Treat upper and lower case strings as distinct.
* attribute (default: None): Filter objects with unique values for this attribute.

### The `unset` filter

The `unset` filter is meant to remove a key/value pair from a dict.

```
{%- set object = {"existing": "value", "disappear": 123} | unset("disappear") -%}
{{ object | tojson }}
```
Will render into:
```
{"existing":"value"}
```

### The `upper` filter
| 🐍 `python` |
|-------------|

Convert a value to uppercase.

### The `urlencode` filter
| 🐍 `python` |
|-------------|

Quote data for use in a URL path or query using UTF-8.

### The `urlize` filter
| 🐍 `python` |
|-------------|

Convert URLs in text into clickable links.

### The `values` filter

The `values` filter is meant to get the values of a map as a list:

```
{{ numbers | values | sort | join(" > ") }}
```
With the following YAML context:
```
numbers:
  first: 1
  second: 2
  third: 3
```
Will render into:
```
1 > 2 > 3
```

### The `wordcount` filter
| 🐍 `python` |
|-------------|

Count the words in that string.

### The `wordwrap` filter
| 🐍 `python` |
|-------------|

Wrap a string to the given width. Existing newlines are treated as paragraphs to be wrapped separately.

### The `xmlattr` filter
| 🐍 `python` |
|-------------|

Create an SGML/XML attribute string based on the items in a dict.