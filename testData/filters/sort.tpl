{{ [3, 4, 5, 0, 1, 2]|sort }}
{{ [3, 4, 5, 0, 1, 2]|sort(true) }}
{{ "bacd"|sort }}
{{ "bacd"|sort(true) }}
{{ "bACd"|sort(case_sensitive=true) }}
{{ "bACd"|sort(true, case_sensitive=true) }}
{{ ['a', 'AA', 'aaa', 'AAA', 'AAAA']|sort }}
{{ ['a', 'AA', 'aaa', 'AAA', 'AAAA']|sort(true) }}
{{ ['a', 'AA', 'aaa', 'AAA', 'AAAA']|sort(case_sensitive=true) }}
{{ ['a', 'AA', 'aaa', 'AAA', 'AAAA']|sort(true, true) }}
{{ {"a": "a", "key": "value", "other": 42, "ALL": "CAPS"}|sort }}
{{ {"a": "a", "key": "value", "other": 42, "ALL": "CAPS"}|sort(true) }}
{{ {"a": "a", "key": "value", "other": 42, "ALL": "CAPS"}|sort(case_sensitive=true) }}
{{ {"a": "a", "key": "value", "other": 42, "ALL": "CAPS"}|sort(true, true) }}
{{ simple.casedStrmap|sort }}
{{ simple.casedStrmap|sort(true) }}
{{ simple.casedStrmap|sort(case_sensitive=true) }}
{{ simple.casedStrmap|sort(true, true) }}
{{ [] | sort }}
