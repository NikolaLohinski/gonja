# indent 4
{{ simple.long_text|indent }}
# indent 2
{{ simple.long_text|indent(2) }}
# indent 2, first=true
{{ simple.long_text|indent(2, true) }}
# indent 2, blank=true
{{ simple.long_text|indent(2, false, true) }}
# indent 2, first=true, blank=true
{{ simple.long_text|indent(2, true, true) }}
