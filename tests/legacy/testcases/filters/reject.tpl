{{ [true, false, 0, 1]|reject }}
{{ [0, 1, 2, 3, 4, 5]|reject('odd') }}
{{ [0, 1, 2, 3, 4, 5]|reject('ge', 3) }}
