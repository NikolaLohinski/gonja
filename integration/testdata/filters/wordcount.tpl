{{ ""|wordcount }}
{% filter wordcount %}{{ lorem(25, "w") }}{% endfilter %}
