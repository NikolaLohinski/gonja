{{ {"key": "will disappear", "still": "there"} | unset("key") }}
{{ {} | unset("nope") }}