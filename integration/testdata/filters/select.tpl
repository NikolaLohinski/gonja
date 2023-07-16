{{ [true, false, 0, 1] | select }}
{{ [0, 1, 2, 3, 4, 5] | select('odd') }}
{{ [0, 1, 2, 3, 4, 5] | select('ge', 3) }}
{{ ["one", "two", "three"] | select("!=", "two") | join(" is smaller than ") }}
