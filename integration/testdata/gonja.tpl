{%- if "v" ~ gonja.version != CI_COMMIT_TAG -%}
v{{- gonja.version }} != {{ CI_COMMIT_TAG }}
{%- else -%}
"v" ~ gonja.version == CI_COMMIT_TAG
{%- endif -%}
