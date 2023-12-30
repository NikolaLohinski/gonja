{{ "http://www.john-doe.de"|urlize|safe }}
{{ "http://www.john-doe.de"|urlize(rel='nofollow')|safe }}
{{ "http://www.john-doe.de"|urlize(rel='nofollow', target='_blank')|safe }}
{{ "http://www.john-doe.de"|urlize(rel='noopener')|safe }}
{{ "www.john-doe.de"|urlize|safe }}
{{ "john-doe.de"|urlize|safe }}
--
{% filter urlize|safe %}
Please mail me at demo@example.com or visit mit on:
- lorem ipsum github.com/nikolalohinski/gonja/v2 lorem ipsum
- lorem ipsum http://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de/test="test" lorem ipsum
{% endfilter %}
--
{% filter urlize(target='_blank', rel="nofollow")|safe %}
Please mail me at demo@example.com or visit mit on:
- lorem ipsum github.com/nikolalohinski/gonja/v2 lorem ipsum
- lorem ipsum http://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de/test="test" lorem ipsum
{% endfilter %}
--
{% filter urlize(15)|safe %}
Please mail me at demo@example.com or visit mit on:
- lorem ipsum github.com/nikolalohinski/gonja/v2 lorem ipsum
- lorem ipsum http://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum https://www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de lorem ipsum
- lorem ipsum www.john-doe.de/test="test" lorem ipsum
{% endfilter %}
