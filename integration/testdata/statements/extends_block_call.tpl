{% extends "../inheritance/base.tpl" %}

{% block title %}Title{% endblock %}

{% block body %}Body: {{ self.title() }}{% endblock %}
