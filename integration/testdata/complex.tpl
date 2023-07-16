{# A more complex template using gonja #}
<!DOCTYPE html>
<html>

<head>
	<title>My blog page</title>
</head>

<body>
	<h1>Blogpost</h1>
	<div id="content">
		{{ complex.post.Text|safe }}
	</div>

	<h1>Comments</h1>

	{% for comment in complex.comments %}
		<h2>{{ loop.index }}. Comment ({{ loop.revindex}} comment{% if loop.revindex > 1 %}s{% endif %} left)</h2>
		<p>From: {{ comment.Author.Name }} ({% if comment.Author.Validated %}validated{% else %}not validated{% endif %})</p>

		{% if complex.is_admin(comment.Author) %}
			<p>This user is an admin (verify: {{ comment.Author.IsAdmin() }})!</p>
		{% else %}
			<p>This user is not admin!</p>
		{% endif %}

		<p>Written {{ comment.Date }}</p>
		<p>{{ comment.Text|striptags }}</p>
	{% endfor %}
</body>

</html>
