# Global Variables

Global variables are jinja variables available in the global scope by default.

```
{{ gonja.version }}
```

## The `gonja` object      

A dictionary containing information about the `gonja` library, with the following properties:
* `version` - the version of the library in use, which is `0.0.0+trunk` if using any commit from `master` branch
