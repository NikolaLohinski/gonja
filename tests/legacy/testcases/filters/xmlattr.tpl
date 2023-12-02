|{{ {'class': 'my-class', 'missing': none, 'id': 'list-42'}|xmlattr|safe }}
|{{ {'class': 'my-class', 'missing': none, 'id': 'list-42'}|xmlattr(false)|safe }}
|{{ {'key': 'value"'}|xmlattr|safe }}
