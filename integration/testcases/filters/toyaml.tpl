{{ {} | toyaml }}
{{ {"simple": {"nested": "field"}} | toyaml }}
{{ ["array"] | toyaml }}
{{ "string" | toyaml }}
{{ 42 | toyaml }}
{{ {"indented": {4: "spaces"}} | toyaml(indent=4) }}
{{ {"indented": {4: "spaces"}} | toyaml(indent=4) }}