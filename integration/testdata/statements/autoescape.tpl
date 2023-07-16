{{ "<script>alert('xss');</script>" }}
{% autoescape true %}
{{ "<script>alert('xss');</script>" }}
{% endautoescape %}
{% autoescape false %}
{{ "<script>alert('xss');</script>" }}
{% endautoescape %}
{% autoescape false %}
{{ "<script>alert('xss');</script>"|escape }}
{% endautoescape %}
