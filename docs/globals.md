## Globals

Globals (or _global functions_) are helpers available in the global scope by default.

```
{% for index in range(10) %}
counting {{ index + 1 }}
{% endfor %}
```

Any global function that is also implemented in the `python` version of the Jinja engine will be marked with the following tag:

| 🐍 `python` |
|-------------|

For any of those, the [official documentation for `python`'s Jinja implementation](https://jinja.palletsprojects.com/en/3.1.x/templates/#list-of-global-functions) can be used as additional reference.


### The `dict` function      
| 🐍 `python` |
|-------------|

A convenient alternative to dict literals. `{'foo': 'bar'}` is the same as `dict(foo='bar')`.

### The `namespace` function 
| 🐍 `python` |
|-------------|

Creates a new container that allows attribute assignment using the `{% set %}` tag:

```
{% set ns = namespace() %}
{% set ns.foo = 'bar' %}
```

The main purpose of this is to allow carrying a value from within a loop body to an outer scope. Initial values can be provided as a dict, as keyword arguments, or both (same behavior as Python’s dict constructor):

```
{% set ns = namespace(found=false) %}
{% for item in items %}
    {% if item.check_something() %}
        {% set ns.found = true %}
    {% endif %}
    * {{ item.title }}
{% endfor %}
Found item having something: {{ ns.found }}
```

### The `range` function     
| 🐍 `python` |
|-------------|

Return a list containing an arithmetic progression of integers. `range(i, j)` returns _[i, i+1, i+2, ..., j-1]_; the `start` (!) defaults to `0`. When a `step` is given, it specifies the increment (or decrement). For example, `range(4)` and `range(0, 4, 1)` return _[0, 1, 2, 3]_.

### The `cycler` function
| 🐍 `python` |
|-------------|

Cycle through values by yielding them one at a time, then restarting once the end is reached.

Similar to `loop.cycle`, but can be used outside loops or across multiple loops. For example, render a list of folders and files in a list, alternating giving them “odd” and “even” classes.

```html
{% set row_class = cycler("odd", "even") %}
<ul class="browser">
{% for folder in folders %}
  <li class="folder {{ row_class.next() }}">{{ folder }}
{% endfor %}
{% for file in files %}
  <li class="file {{ row_class.next() }}">{{ file }}
{% endfor %}
</ul>
```

### The `joiner` function    
| 🐍 `python` |
|-------------|

A tiny helper that can be used to “join” multiple sections. A `joiner` is passed a string and will return that string every time it’s called, except the first time (in which case it returns an empty string). You can use this to join things:

```html
{% set pipe = joiner("|") %}
{% if categories %} {{ pipe() }}
    Categories: {{ categories|join(", ") }}
{% endif %}
{% if author %} {{ pipe() }}
    Author: {{ author() }}
{% endif %}
{% if can_edit %} {{ pipe() }}
    <a href="?action=edit">Edit</a>
{% endif %}
```

### The `lipsum` function    
| 🐍 `python` |
|-------------|

Generates some lorem ipsum for the template. By default, five paragraphs of HTML are generated with each paragraph between 20 and 100 words. If html is False, regular text is returned. This is useful to generate simple contents for layout testing.
